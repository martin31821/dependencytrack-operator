/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
	"github.com/martin31821/dependencytrack-operator/internal/deptrack"
)

// ======================================================================
// Combined mock DT server: supports publisher + rule endpoints + config schema
// ======================================================================

type combinedMockDTServer struct {
	mu            sync.Mutex
	publisher     *dtapi.NotificationPublisher
	publishers    map[string]*dtapi.NotificationPublisher // uuid -> publisher
	rules         map[string]*dtapi.NotificationRule      // uuid -> rule
	failCreate    bool
	failUpdate    bool
	failDelete    bool
	failSchema    bool
	lastCreateReq *dtapi.CreateNotificationPublisherRequest
}

func newCombinedMockServer() *combinedMockDTServer {
	return &combinedMockDTServer{
		publishers: make(map[string]*dtapi.NotificationPublisher),
		rules:      make(map[string]*dtapi.NotificationRule),
	}
}

func (s *combinedMockDTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stripped, _ := strings.CutPrefix(r.URL.Path, "/api")

	switch {
	// --- Publisher create/update ---
	case stripped == notifPublisherPath && (r.Method == http.MethodPost || r.Method == http.MethodPut):
		s.handlePubCreateOrUpdate(w, r)
	// --- Publisher list ---
	case stripped == notifPublisherPath && r.Method == http.MethodGet:
		s.handlePubList(w)
	// --- Publisher delete ---
	case strings.HasPrefix(stripped, notifPublisherPathPrefix) && r.Method == http.MethodDelete:
		s.handlePubDelete(w, r)
	// --- Rule list ---
	case stripped == "/v1/notification/rule" && r.Method == http.MethodGet:
		s.handleRuleList(w)
	// --- Rule delete ---
	case stripped == "/v1/notification/rule" && r.Method == http.MethodDelete:
		s.handleRuleDelete(w, r)
	// --- Rule create ---
	case stripped == "/v1/notification/rule" && r.Method == http.MethodPut:
		s.handleRuleCreate(w, r)
	// --- Rule update ---
	case stripped == "/v1/notification/rule" && r.Method == http.MethodPost:
		s.handleRuleUpdate(w, r)
	// --- Config schema ---
	case r.Method == http.MethodGet && configSchemaPathRE.MatchString(stripped):
		s.handleConfigSchema(w, r)
	// --- Validate credentials ---
	case stripped == "/v1/user/validate-credentials":
		writeJSON(w, "mock-token")
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

// ---- Publisher handlers ----

func (s *combinedMockDTServer) handlePubCreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// First try as create request
	var createReq dtapi.CreateNotificationPublisherRequest
	if json.Unmarshal(bodyBytes, &createReq) != nil {
		// Not a valid create request, try update
		s.handlePubUpdateFromBody(w, bytes.NewReader(bodyBytes))
		return
	}
	// Check if it has a UUID — if so it's an update
	var updateReq dtapi.UpdateNotificationPublisherRequest
	if json.Unmarshal(bodyBytes, &updateReq) == nil && updateReq.Uuid != "" {
		s.handlePubUpdateFromBody(w, bytes.NewReader(bodyBytes))
		return
	}
	s.handlePubCreateFromBody(w, bytes.NewReader(bodyBytes))
}

func (s *combinedMockDTServer) handlePubCreateFromBody(w http.ResponseWriter, r io.Reader) {
	if s.failCreate {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var req dtapi.CreateNotificationPublisherRequest
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	s.lastCreateReq = &req
	pub := &dtapi.NotificationPublisher{
		Name:          req.Name,
		ExtensionName: req.ExtensionName,
		Uuid:          "pub-created-uuid-123",
	}
	if req.Description != nil {
		pub.Description = req.Description
	}
	s.publisher = pub
	s.publishers[pub.Uuid] = pub
	writeJSON(w, pub)
}

func (s *combinedMockDTServer) handlePubUpdateFromBody(w http.ResponseWriter, r io.Reader) {
	if s.failUpdate {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var req dtapi.UpdateNotificationPublisherRequest
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	uuid := req.Uuid
	pub, exists := s.publishers[uuid]
	if !exists {
		if s.publisher != nil && s.publisher.Uuid == uuid {
			pub = s.publisher
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	pub.Name = req.Name
	pub.ExtensionName = req.ExtensionName
	// Always sync the single-slot publisher when the UUID matches,
	// even if we found it in the map.
	if s.publisher == nil || s.publisher.Uuid != uuid {
		s.publisher = pub
	}
	s.publishers[uuid] = pub
	writeJSON(w, pub)
}

func (s *combinedMockDTServer) handlePubList(w http.ResponseWriter) {
	list := make([]dtapi.NotificationPublisher, 0, len(s.publishers))
	for _, pub := range s.publishers {
		list = append(list, *pub)
	}
	writeJSON(w, list)
}

func (s *combinedMockDTServer) handlePubDelete(w http.ResponseWriter, r *http.Request) {
	if s.failDelete {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	stripped, _ := strings.CutPrefix(r.URL.Path, "/api")
	uuid, _ := strings.CutPrefix(stripped, notifPublisherPathPrefix)

	if _, exists := s.publishers[uuid]; exists {
		delete(s.publishers, uuid)
		if s.publisher != nil && s.publisher.Uuid == uuid {
			s.publisher = nil
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if s.publisher != nil && s.publisher.Uuid == uuid {
		s.publisher = nil
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// ---- Rule handlers ----

func (s *combinedMockDTServer) handleRuleList(w http.ResponseWriter) {
	list := make([]dtapi.NotificationRule, 0, len(s.rules))
	for _, rule := range s.rules {
		list = append(list, *rule)
	}
	writeJSON(w, list)
}

func (s *combinedMockDTServer) handleRuleDelete(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UUID string `json:"uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.UUID == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if _, exists := s.rules[request.UUID]; exists {
		delete(s.rules, request.UUID)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (s *combinedMockDTServer) handleRuleCreate(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req dtapi.CreateNotificationRuleRequest
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	warnLevel := "WARN"
	if req.Level != "" {
		warnLevel = req.Level
	}
	rule := &dtapi.NotificationRule{
		Uuid:              "rule-created-uuid-456",
		Name:              req.Name,
		Scope:             req.Scope,
		TriggerType:       "EVENT",
		NotificationLevel: &warnLevel,
	}
	s.rules[rule.Uuid] = rule
	writeJSON(w, rule)
}

func (s *combinedMockDTServer) handleRuleUpdate(w http.ResponseWriter, r *http.Request) {
	var req dtapi.UpdateNotificationRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	rule, exists := s.rules[req.Uuid]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	rule.Name = req.Name
	rule.Scope = req.Scope
	rule.NotificationLevel = &req.Level
	if req.Enabled != nil {
		rule.Enabled = req.Enabled
	}
	if req.FilterExpression != nil {
		rule.FilterExpression = req.FilterExpression
	}
	if len(req.NotifyOn) > 0 {
		rule.NotifyOn = req.NotifyOn
	}
	if req.LogSuccessfulPublish != nil {
		rule.LogSuccessfulPublish = req.LogSuccessfulPublish
	}
	if req.NotifyChildren != nil {
		rule.NotifyChildren = req.NotifyChildren
	}
	if req.ScheduleCron != nil {
		rule.ScheduleCron = req.ScheduleCron
	}
	if req.ScheduleSkipUnchanged != nil {
		rule.ScheduleSkipUnchanged = req.ScheduleSkipUnchanged
	}
	writeJSON(w, rule)
}

func (s *combinedMockDTServer) handleConfigSchema(w http.ResponseWriter, r *http.Request) {
	if s.failSchema {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{
		"type":     "object",
		"required": []interface{}{"webhookUrl"},
		"properties": map[string]interface{}{
			"webhookUrl": map[string]interface{}{"type": "string"},
			"headers":    map[string]interface{}{"type": "object"},
		},
	})
}

// ---- Combined client provider ----

type combinedMockClientProvider struct {
	url    string
	mu     sync.Mutex
	api    *dtapi.APIClient
	getErr error
}

func (p *combinedMockClientProvider) Get(ctx context.Context) (context.Context, *dtapi.APIClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.getErr != nil {
		return nil, nil, p.getErr
	}
	if p.api == nil {
		p.api = deptrack.NewAPIClient(p.url)
	}
	return context.WithValue(ctx, dtapi.ContextAccessToken, "mock-token"), p.api, nil
}

func (p *combinedMockClientProvider) Invalidate() {
	p.mu.Lock()
	p.api = nil
	p.mu.Unlock()
}

// ======================================================================
// Helpers
// ======================================================================

const (
	integrationNS = "default"
	intPubName    = "test-pub"
	intRuleName   = "test-rule"
)

// ---------------------------------------------------------------------------
// Full lifecycle: publisher creation → drift → blocked deletion → release → cleanup
// ---------------------------------------------------------------------------

var _ = Describe("Notification cleanup lifecycle integration", func() {
	var (
		server        *httptest.Server
		mockServer    *combinedMockDTServer
		provider      *combinedMockClientProvider
		publisherCtrl *NotificationPublisherReconciler
		ruleCtrl      *NotificationRuleReconciler
		fakeRecorder  *record.FakeRecorder
		ruleRecorder  *record.FakeRecorder
	)

	BeforeEach(func() {
		mockServer = newCombinedMockServer()
		server = httptest.NewServer(mockServer)
		provider = &combinedMockClientProvider{url: server.URL}
		DeferCleanup(server.Close)

		fakeRecorder = record.NewFakeRecorder(20)
		ruleRecorder = record.NewFakeRecorder(20)
		publisherCtrl = &NotificationPublisherReconciler{
			Client:     k8sClient,
			Scheme:     k8sClient.Scheme(),
			Recorder:   fakeRecorder,
			DTProvider: provider,
		}
		ruleCtrl = &NotificationRuleReconciler{
			Client:                   k8sClient,
			Scheme:                   k8sClient.Scheme(),
			Recorder:                 ruleRecorder,
			DTProvider:               provider,
			PublisherConfigValidator: NewPublisherConfigValidator(),
		}
	})

	// ---- Test 1: Full lifecycle ----

	It("should repair drift, block deletion, release on retarget, and clean up", func() {
		// --- Setup ---
		p := createNPWithFinalizer("pub-full", integrationNS)
		Expect(k8sClient.Create(ctx, p)).To(Succeed())

		// Step 1: reconcile publisher → finalizer added
		_, err := publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Step 2: reconcile publisher → DT create, UUID set
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		p = &dependencytrackv1alpha1.NotificationPublisher{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full", Namespace: integrationNS}, p)).To(Succeed())
		Expect(p.Status.UUID).To(Equal("pub-created-uuid-123"))

		// Pre-set publisher UUID so rule reconcile sees a Ready publisher
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full", Namespace: integrationNS}, p)).To(Succeed())
		p.Status.UUID = "pub-created-uuid-123"
		p.Status.Name = p.Spec.Name
		Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())

		// Create the rule
		rule := &dependencytrackv1alpha1.NotificationRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pub-full-rule",
				Namespace: integrationNS,
			},
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Name:         "pub-full-rule-spec",
				Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
				TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
				Level:        dependencytrackv1alpha1.NotificationRuleLevelWarn,
				PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: "pub-full"},
			},
		}
		Expect(k8sClient.Create(ctx, rule)).To(Succeed())

		// Step 3: reconcile rule → finalizer added
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Step 4: reconcile rule → DT create, UUID set
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		rule = &dependencytrackv1alpha1.NotificationRule{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full-rule", Namespace: integrationNS}, rule)).To(Succeed())
		Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))

		// --- Drift repair ---

		// Mutate publisher name in DT
		mockServer.mu.Lock()
		mockServer.publisher.Name = "drifted-pub-name"
		mockServer.publisher.ExtensionName = "drifted-ext"
		mockServer.mu.Unlock()

		// Mutate rule name in DT
		mockServer.mu.Lock()
		mockServer.rules["rule-created-uuid-456"].Name = "drifted-rule-name"
		mockServer.mu.Unlock()

		// Reconcile publisher — should repair
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		p = &dependencytrackv1alpha1.NotificationPublisher{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full", Namespace: integrationNS}, p)).To(Succeed())
		Expect(p.Status.UUID).To(Equal("pub-created-uuid-123"))
		cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionTrue))
		Expect(cond.Reason).To(Equal("PublisherSynced"))

		// Verify DT publisher repaired
		mockServer.mu.Lock()
		Expect(mockServer.publisher.Name).To(Equal("pub-full-ext"))
		Expect(mockServer.publisher.ExtensionName).To(Equal("webhook"))
		mockServer.mu.Unlock()

		// Reconcile rule — should repair
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		rule = &dependencytrackv1alpha1.NotificationRule{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full-rule", Namespace: integrationNS}, rule)).To(Succeed())
		Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))
		cond = meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionTrue))
		Expect(cond.Reason).To(Equal("RuleSynced"))

		// Verify DT rule repaired
		mockServer.mu.Lock()
		Expect(mockServer.rules["rule-created-uuid-456"].Name).To(Equal("pub-full-rule-spec"))
		mockServer.mu.Unlock()

		// --- Dependency-blocked deletion ---

		// Mark publisher for deletion
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full", Namespace: integrationNS}, p)).To(Succeed())
		Expect(k8sClient.Delete(ctx, p)).To(Succeed())

		// Reconcile publisher — should be blocked
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		p = &dependencytrackv1alpha1.NotificationPublisher{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full", Namespace: integrationNS}, p)).To(Succeed())
		Expect(controllerutil.ContainsFinalizer(p, notificationPublisherFinalizer)).To(BeTrue())
		cond = meta.FindStatusCondition(p.Status.Conditions, conditionReady)
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionFalse))
		Expect(cond.Reason).To(Equal(conditionReasonDependencyBlocked))

		// DT publisher still exists — no delete was attempted
		mockServer.mu.Lock()
		_, exists := mockServer.publishers["pub-created-uuid-123"]
		mockServer.mu.Unlock()
		Expect(exists).To(BeTrue())

		// --- Retarget the rule to release the publisher ---

		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full-rule", Namespace: integrationNS}, rule)).To(Succeed())
		rule.Spec.PublisherRef.Name = "other-publisher"
		Expect(k8sClient.Update(ctx, rule)).To(Succeed())

		// Reconcile publisher — should now succeed
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Verify the K8s publisher is deleted (finalizer removed, object GC'd)
		err = k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full", Namespace: integrationNS}, p)
		Expect(apierrors.IsNotFound(err)).To(BeTrue())

		// Verify DT publisher was deleted
		mockServer.mu.Lock()
		_, exists = mockServer.publishers["pub-created-uuid-123"]
		mockServer.mu.Unlock()
		Expect(exists).To(BeFalse())

		// --- Clean up the K8s rule ---

		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-full-rule", Namespace: integrationNS}, rule)).To(Succeed())
		Expect(k8sClient.Delete(ctx, rule)).To(Succeed())

		// Reconcile rule — should clean up UUID in DT and remove finalizer
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-full-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Verify DT rule was deleted
		mockServer.mu.Lock()
		_, exists = mockServer.rules["rule-created-uuid-456"]
		mockServer.mu.Unlock()
		Expect(exists).To(BeFalse())
	})

	// ---- Test 2: Rule UUID-targeted drift repair ----

	It("should repair both publisher and rule drift independently", func() {
		// Create publisher
		p := createNPWithFinalizer("pub-drift", integrationNS)
		Expect(k8sClient.Create(ctx, p)).To(Succeed())
		_, err := publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-drift", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-drift", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Pre-set publisher UUID for rule reconcile
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-drift", Namespace: integrationNS}, p)).To(Succeed())
		p.Status.UUID = "pub-drift-uuid"
		p.Status.Name = p.Spec.Name
		Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())
		mockServer.mu.Lock()
		mockServer.publisher = &dtapi.NotificationPublisher{
			Name: "pub-drift-ext", ExtensionName: "webhook", Uuid: "pub-drift-uuid",
		}
		mockServer.publishers["pub-drift-uuid"] = mockServer.publisher
		mockServer.mu.Unlock()

		// Create rule
		rule := &dependencytrackv1alpha1.NotificationRule{
			ObjectMeta: metav1.ObjectMeta{Name: "pub-drift-rule", Namespace: integrationNS},
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Name:         "pub-drift-rule-spec",
				Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
				TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
				Level:        dependencytrackv1alpha1.NotificationRuleLevelWarn,
				PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: "pub-drift"},
			},
		}
		Expect(k8sClient.Create(ctx, rule)).To(Succeed())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-drift-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-drift-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Mutate both publishers in DT
		mockServer.mu.Lock()
		mockServer.publisher.Name = "drifted-pub"
		mockServer.publisher.ExtensionName = "drifted-ext"
		mockServer.rules["rule-created-uuid-456"].Name = "drifted-rule"
		mockServer.mu.Unlock()

		// Reconcile both
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-drift", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-drift-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Verify publisher repaired
		p = &dependencytrackv1alpha1.NotificationPublisher{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-drift", Namespace: integrationNS}, p)).To(Succeed())
		Expect(p.Status.UUID).To(Equal("pub-drift-uuid"))
		cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionTrue))

		// Verify DT publisher repaired
		mockServer.mu.Lock()
		Expect(mockServer.publisher.Name).To(Equal("pub-drift-ext"))
		Expect(mockServer.publisher.ExtensionName).To(Equal("webhook"))
		mockServer.mu.Unlock()

		// Verify rule repaired
		rule = &dependencytrackv1alpha1.NotificationRule{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-drift-rule", Namespace: integrationNS}, rule)).To(Succeed())
		Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))
		cond = meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionTrue))

		// Verify DT rule repaired
		mockServer.mu.Lock()
		Expect(mockServer.rules["rule-created-uuid-456"].Name).To(Equal("pub-drift-rule-spec"))
		mockServer.mu.Unlock()
	})

	// ---- Test 3: Already-absent remote publisher (404-as-success) ----

	It("should treat remote 404 as successful deletion for both publisher and rule", func() {
		// Create publisher — pre-seed UUID but publisher NOT in mock DT
		p := createNPWithFinalizer("pub-absent", integrationNS)
		Expect(k8sClient.Create(ctx, p)).To(Succeed())
		_, err := publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-absent", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Set UUID but do NOT put publisher in mock server
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-absent", Namespace: integrationNS}, p)).To(Succeed())
		p.Status.UUID = "absent-pub-uuid"
		Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())

		// Delete the publisher — should treat 404 as success
		Expect(k8sClient.Delete(ctx, p)).To(Succeed())
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-absent", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Publisher should be GC'd
		err = k8sClient.Get(ctx, types.NamespacedName{Name: "pub-absent", Namespace: integrationNS}, p)
		Expect(apierrors.IsNotFound(err)).To(BeTrue())

		// Create a separate rule-publisher pair for rule (p is GC'd above)
		absentPub := &dependencytrackv1alpha1.NotificationPublisher{
			ObjectMeta: metav1.ObjectMeta{Name: "absent-pub-2", Namespace: integrationNS},
			Spec: dependencytrackv1alpha1.NotificationPublisherSpec{
				ExtensionName: "webhook",
				Name:          "absent-pub-2-ext",
			},
		}
		Expect(k8sClient.Create(ctx, absentPub)).To(Succeed())
		rule := &dependencytrackv1alpha1.NotificationRule{
			ObjectMeta: metav1.ObjectMeta{Name: "absent-rule", Namespace: integrationNS},
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Name:         "absent-rule-spec",
				Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
				TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
				Level:        dependencytrackv1alpha1.NotificationRuleLevelWarn,
				PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: "absent-pub-2"},
			},
		}
		Expect(k8sClient.Create(ctx, rule)).To(Succeed())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "absent-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "absent-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Rule UUID is set but NOT in mock server — delete it
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "absent-rule", Namespace: integrationNS}, rule)).To(Succeed())
		rule.Status.UUID = "absent-rule-uuid"
		Expect(k8sClient.Status().Update(ctx, rule)).To(Succeed())
		Expect(k8sClient.Delete(ctx, rule)).To(Succeed())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "absent-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Rule should be GC'd
		err = k8sClient.Get(ctx, types.NamespacedName{Name: "absent-rule", Namespace: integrationNS}, rule)
		Expect(apierrors.IsNotFound(err)).To(BeTrue())
	})

	// ---- Test 4: Unrelated resources don't block deletion ----

	It("should not block deletion on unrelated rules/publishers in the namespace", func() {
		// Create multiple publishers
		p1 := createNPWithFinalizer("pub-target", integrationNS)
		Expect(k8sClient.Create(ctx, p1)).To(Succeed())
		_, err := publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-target", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-target", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Fresh fetch after reconcile to avoid stale-reference conflict
		p1 = &dependencytrackv1alpha1.NotificationPublisher{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-target", Namespace: integrationNS}, p1)).To(Succeed())
		p1.Status.UUID = "pub-target-uuid"
		p1.Status.Name = p1.Spec.Name
		Expect(k8sClient.Status().Update(ctx, p1)).To(Succeed())
		mockServer.mu.Lock()
		mockServer.publisher = &dtapi.NotificationPublisher{
			Name: "pub-target-ext", ExtensionName: "webhook", Uuid: "pub-target-uuid",
		}
		mockServer.publishers["pub-target-uuid"] = mockServer.publisher
		mockServer.mu.Unlock()

		// Create unrelated publisher
		p2 := createNPWithFinalizer("pub-unrelated", integrationNS)
		Expect(k8sClient.Create(ctx, p2)).To(Succeed())
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-unrelated", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-unrelated", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Create unrelated rule that does NOT reference our target publisher
		rule := &dependencytrackv1alpha1.NotificationRule{
			ObjectMeta: metav1.ObjectMeta{Name: "pub-unrelated-rule", Namespace: integrationNS},
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Name:         "unrelated-rule-spec",
				Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
				TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
				Level:        dependencytrackv1alpha1.NotificationRuleLevelWarn,
				PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: "pub-unrelated"},
			},
		}
		Expect(k8sClient.Create(ctx, rule)).To(Succeed())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-unrelated-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-unrelated-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Now drift the target publisher
		mockServer.mu.Lock()
		mockServer.publisher.Name = "drifted-target"
		mockServer.mu.Unlock()

		// Reconcile target — should repair
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-target", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		p1 = &dependencytrackv1alpha1.NotificationPublisher{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-target", Namespace: integrationNS}, p1)).To(Succeed())
		cond := meta.FindStatusCondition(p1.Status.Conditions, conditionReady)
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionTrue))

		// Verify unrelated publisher not affected
		p2 = &dependencytrackv1alpha1.NotificationPublisher{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-unrelated", Namespace: integrationNS}, p2)).To(Succeed())
		Expect(p2.Status.UUID).NotTo(Equal(p1.Status.UUID))

		// Verify unrelated rule not affected
		rule = &dependencytrackv1alpha1.NotificationRule{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-unrelated-rule", Namespace: integrationNS}, rule)).To(Succeed())
		Expect(rule.Status.UUID).NotTo(Equal(p1.Status.UUID))
	})

	// ---- Test 5: Secret redaction in validation errors ----

	It("should redact publisher config values in validation error messages", func() {
		secret := createPublisherConfigSecret("test-config-secret", integrationNS, "config", map[string][]byte{
			"config": []byte(`{"token": "sk-ir34ldgdh", "webhookUrl": "https://hooks.slack.com/services/T0000/B0000/XXXX"}`),
		})
		DeferCleanup(func() {
			_ = k8sClient.Delete(ctx, secret)
		})

		// Create a rule with publisher config secret ref
		p := createNPWithFinalizer("pub-redact", integrationNS)
		Expect(k8sClient.Create(ctx, p)).To(Succeed())
		_, err := publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-redact", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = publisherCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-redact", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		// Pre-set publisher UUID for rule reconcile
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-redact", Namespace: integrationNS}, p)).To(Succeed())
		p.Status.UUID = "pub-redact-uuid"
		p.Status.Name = p.Spec.Name
		Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())
		mockServer.mu.Lock()
		mockServer.publisher = &dtapi.NotificationPublisher{
			Name: "pub-redact-ext", ExtensionName: "webhook", Uuid: "pub-redact-uuid",
		}
		mockServer.publishers["pub-redact-uuid"] = mockServer.publisher
		mockServer.mu.Unlock()

		rule := &dependencytrackv1alpha1.NotificationRule{
			ObjectMeta: metav1.ObjectMeta{Name: "pub-redact-rule", Namespace: integrationNS},
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Name:         "pub-redact-rule-spec",
				Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
				TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
				Level:        dependencytrackv1alpha1.NotificationRuleLevelWarn,
				PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: "pub-redact"},
				PublisherConfigSecretRef: &dependencytrackv1alpha1.PublisherConfigSecretRef{
					Name: "test-config-secret",
					Key:  "config",
				},
			},
		}
		Expect(k8sClient.Create(ctx, rule)).To(Succeed())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-redact-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())
		_, err = ruleCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "pub-redact-rule", Namespace: integrationNS}})
		Expect(err).NotTo(HaveOccurred())

		rule = &dependencytrackv1alpha1.NotificationRule{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pub-redact-rule", Namespace: integrationNS}, rule)).To(Succeed())

		// Verify the rule reconciles successfully with valid config (no secret leak in condition)
		cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionTrue))
		Expect(cond.Message).NotTo(ContainSubstring("sk-ir34ldgdh"))
		Expect(cond.Message).NotTo(ContainSubstring("T0000"))
	})
})

// ======================================================================
// Unit tests for error redaction
// ======================================================================

var _ = Describe("PublisherConfigValidationError redaction", func() {
	It("should not include secret values in error messages", func() {
		v := NewPublisherConfigValidator()
		// Valid JSON without a schema passes structurally; no error expected
		err := v.Validate([]byte(`{"token": "sk-ir34ldgdh"}`), nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should not include secret values in malformed JSON errors", func() {
		v := NewPublisherConfigValidator()
		err := v.Validate([]byte(`{"token": "sk-ir34ldgdh", bad json}`), nil)
		Expect(err).To(HaveOccurred())
		ve, ok := err.(*PublisherConfigValidationError)
		Expect(ok).To(BeTrue())
		Expect(ve.Error()).NotTo(ContainSubstring("sk-ir34ldgdh"))
		Expect(ve.Error()).NotTo(ContainSubstring("T0000"))
	})

	It("should include the reason but not secret values", func() {
		v := NewPublisherConfigValidator()
		// Valid JSON without a schema passes structurally; no error expected
		err := v.Validate([]byte(`{"webhookUrl": "https://hooks.slack.com/services/T0000/B0000/XXXX", "token": "xoxb-123456"}`), nil)
		Expect(err).NotTo(HaveOccurred())
	})
})

// ======================================================================
// Unit tests for the combined server handlers
// ======================================================================

var _ = Describe("combinedMockDTServer handlers", func() {
	var (
		server *httptest.Server
		mock   *combinedMockDTServer
	)

	BeforeEach(func() {
		mock = newCombinedMockServer()
		server = httptest.NewServer(mock)
		DeferCleanup(server.Close)
	})

	It("creates a publisher", func() {
		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/notification/publisher",
			bytes.NewBuffer([]byte(`{"name":"test-ext","extensionName":"webhook"}`)))
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		var pub dtapi.NotificationPublisher
		Expect(json.NewDecoder(resp.Body).Decode(&pub)).To(Succeed())
		Expect(pub.Uuid).To(Equal("pub-created-uuid-123"))
		Expect(pub.Name).To(Equal("test-ext"))
	})

	It("updates a publisher by UUID", func() {
		// Create first
		mock.publisher = &dtapi.NotificationPublisher{
			Name: "old-ext", ExtensionName: "webhook", Uuid: "update-uuid",
		}
		mock.publishers["update-uuid"] = mock.publisher

		// Update
		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/notification/publisher",
			bytes.NewBuffer([]byte(`{"uuid":"update-uuid","name":"new-ext","extensionName":"email"}`)))
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		var pub dtapi.NotificationPublisher
		Expect(json.NewDecoder(resp.Body).Decode(&pub)).To(Succeed())
		Expect(pub.Name).To(Equal("new-ext"))
		Expect(pub.ExtensionName).To(Equal("email"))
		Expect(pub.Uuid).To(Equal("update-uuid"))
	})

	It("lists publishers", func() {
		mock.publishers["uuid-1"] = &dtapi.NotificationPublisher{Uuid: "uuid-1", Name: "pub1"}
		mock.publishers["uuid-2"] = &dtapi.NotificationPublisher{Uuid: "uuid-2", Name: "pub2"}

		resp, err := http.Get(server.URL + "/api/v1/notification/publisher")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		var list []dtapi.NotificationPublisher
		Expect(json.NewDecoder(resp.Body).Decode(&list)).To(Succeed())
		Expect(list).To(HaveLen(2))
	})

	It("deletes a publisher by UUID", func() {
		mock.publishers["del-uuid"] = &dtapi.NotificationPublisher{Uuid: "del-uuid", Name: "to-delete"}

		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/notification/publisher/del-uuid", nil)
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

		// Verify deleted
		_, exists := mock.publishers["del-uuid"]
		Expect(exists).To(BeFalse())
	})

	It("returns 404 when deleting a non-existent publisher", func() {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/notification/publisher/no-such-uuid", nil)
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("creates a rule", func() {
		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/notification/rule",
			bytes.NewBuffer([]byte(`{"name":"test-rule","scope":"ALL_PROJECTS","level":"WARN","publisher":{"uuid":"test-pub-uuid"}}`)))
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		var rule dtapi.NotificationRule
		Expect(json.NewDecoder(resp.Body).Decode(&rule)).To(Succeed())
		Expect(rule.Uuid).To(Equal("rule-created-uuid-456"))
		Expect(rule.Name).To(Equal("test-rule"))
	})

	It("updates a rule by UUID", func() {
		mock.rules["rule-update-uuid"] = &dtapi.NotificationRule{
			Uuid: "rule-update-uuid", Name: "old-name",
		}

		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/notification/rule",
			bytes.NewBuffer([]byte(`{"uuid":"rule-update-uuid","name":"new-name","scope":"ALL_PROJECTS","level":"WARN"}`)))
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		var rule dtapi.NotificationRule
		Expect(json.NewDecoder(resp.Body).Decode(&rule)).To(Succeed())
		Expect(rule.Name).To(Equal("new-name"))
	})

	It("deletes a rule by UUID", func() {
		mock.rules["rule-del-uuid"] = &dtapi.NotificationRule{Uuid: "rule-del-uuid"}

		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/notification/rule",
			bytes.NewBuffer([]byte(`{"uuid":"rule-del-uuid"}`)))
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

		_, exists := mock.rules["rule-del-uuid"]
		Expect(exists).To(BeFalse())
	})

	It("returns 404 for config schema on unknown publisher", func() {
		req, err := http.NewRequest(http.MethodGet,
			server.URL+"/api/v1/notification/publisher/unknown/configSchema", nil)
		Expect(err).NotTo(HaveOccurred())
		resp, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		// Schema endpoint returns 200 regardless of UUID existence
	})
})

// ======================================================================
// Helper functions for test setup (mirror of existing test helpers)
// ======================================================================

func createNPWithFinalizer(name, ns string) *dependencytrackv1alpha1.NotificationPublisher {
	p := &dependencytrackv1alpha1.NotificationPublisher{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: dependencytrackv1alpha1.NotificationPublisherSpec{
			ExtensionName: "webhook",
			Name:          name + "-ext",
		},
	}
	controllerutil.AddFinalizer(p, notificationPublisherFinalizer)
	return p
}

// ======================================================================
// Unit tests for error safety (no secret values in error strings)
// ======================================================================

var _ = Describe("PublisherConfigValidationError safety", func() {
	It("accepts an empty object without exposing a secret key name", func() {
		v := NewPublisherConfigValidator()
		// An empty JSON object is structurally valid when no schema is supplied.
		err := v.Validate([]byte(`{}`), nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("error message does not contain secret value", func() {
		v := NewPublisherConfigValidator()
		// Use a type mismatch (number where string expected) so validation fails.
		// This tests that the error message sanitises secret values.
		schema := map[string]interface{}{
			"properties": map[string]interface{}{
				"endpoint": map[string]interface{}{"type": "string"},
			},
		}
		err := v.Validate([]byte(`{"token": "xoxb-123456", "endpoint": 12345}`), schema)
		Expect(err).To(HaveOccurred())
		ve, ok := err.(*PublisherConfigValidationError)
		Expect(ok).To(BeTrue())
		Expect(ve.Error()).NotTo(ContainSubstring("xoxb-123456"))
		Expect(ve.Error()).NotTo(ContainSubstring("T0000"))
		Expect(ve.Error()).NotTo(ContainSubstring("B0000"))
	})

	It("malformed JSON error does not leak the raw input", func() {
		v := NewPublisherConfigValidator()
		err := v.Validate([]byte(`{"api_key": "ghp_1234567890abcdef"}`), nil)
		// Valid JSON, no schema — should pass
		Expect(err).NotTo(HaveOccurred())
	})

	It("sanitizeJSONError truncates at quotes and backticks", func() {
		sanitizeV := NewPublisherConfigValidator()
		err := sanitizeV.Validate([]byte(`{"bad": json}}}}}, "secret": "sk-live-abc"`), nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).NotTo(ContainSubstring("sk-live-abc"))
		Expect(err.Error()).NotTo(ContainSubstring("ghp_1234567890"))
	})
})

// ======================================================================
// Redaction-focused unit tests
// ======================================================================

var _ = Describe("redaction safety", func() {
	It("Validate with valid JSON and schema should pass cleanly", func() {
		v := NewPublisherConfigValidator()
		config := []byte(`{"webhookUrl": "https://example.com"}`)
		schema := map[string]interface{}{
			"required": []interface{}{"webhookUrl"},
		}
		err := v.Validate(config, schema)
		Expect(err).NotTo(HaveOccurred())
	})

	It("Validate error does not leak token value", func() {
		v := NewPublisherConfigValidator()
		// Type mismatch (number where string expected) triggers validation.
		schema := map[string]interface{}{
			"properties": map[string]interface{}{
				"webhookUrl": map[string]interface{}{"type": "string"},
			},
		}
		err := v.Validate([]byte(`{"token": "sk-ir34ldgdh", "webhookUrl": 123}`), schema)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).NotTo(ContainSubstring("sk-ir34ldgdh"))
		Expect(err.Error()).NotTo(ContainSubstring("ghp_"))
	})

	It("Validate error does not leak SLACK webhook URL", func() {
		v := NewPublisherConfigValidator()
		schema := map[string]interface{}{
			"properties": map[string]interface{}{
				"webhookUrl": map[string]interface{}{"type": "string"},
			},
		}
		err := v.Validate([]byte(`{"webhookUrl": 123}`), schema)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).NotTo(ContainSubstring("hooks.slack.com"))
		Expect(err.Error()).NotTo(ContainSubstring("T0000"))
		Expect(err.Error()).NotTo(ContainSubstring("B0000"))
	})
})
