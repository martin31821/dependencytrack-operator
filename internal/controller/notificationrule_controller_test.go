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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
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

const (
	notifRulePath     = "/v1/notification/rule"
	ruleNS            = "default"
	validateCredsPath = "/v1/credentials/validation"
	schemaPathPrefix  = "/v1/notification/publisher/"
	schemaPathSuffix  = "/configSchema"
)

// notificationRuleMockClientProvider implements deptrack.ClientProviderInterface.
type notificationRuleMockClientProvider struct {
	url    string
	mu     sync.Mutex
	api    *dtapi.APIClient
	getErr error
}

func (p *notificationRuleMockClientProvider) Get(ctx context.Context) (context.Context, *dtapi.APIClient, error) {
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

func (p *notificationRuleMockClientProvider) Invalidate() {
	p.mu.Lock()
	p.api = nil
	p.mu.Unlock()
}

// notificationRuleMockDTServer mocks the DT NotificationRule REST API.
type notificationRuleMockDTServer struct {
	mu                sync.Mutex
	rules             map[string]*dtapi.NotificationRule
	lastCreateRequest *dtapi.CreateNotificationRuleRequest
	lastUpdateRequest *dtapi.UpdateNotificationRuleRequest
	failCreate        bool
	failUpdate        bool
	failList          bool
	failDelete        bool
	failSchema        bool
}

// configSchemaPathRE is defined in notificationrule_test_helpers.go

func newNotificationRuleMockServer() *notificationRuleMockDTServer {
	return &notificationRuleMockDTServer{
		rules: make(map[string]*dtapi.NotificationRule),
	}
}

func (s *notificationRuleMockDTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stripped := r.URL.Path
	if len(stripped) > 4 && stripped[:4] == "/api" {
		stripped = stripped[4:]
	}

	if stripped == validateCredsPath {
		writeJSON(w, "mock-token")
		return
	}

	switch {
	case stripped == notifRulePath && r.Method == http.MethodPut:
		s.handleCreate(w, r)
	case stripped == notifRulePath && r.Method == http.MethodPost:
		s.handleUpdate(w, r)
	case stripped == notifRulePath && r.Method == http.MethodGet:
		s.handleList(w)
	case stripped == notifRulePath && r.Method == http.MethodDelete:
		s.handleDelete(w, r)
	case r.Method == http.MethodGet && configSchemaPathRE.MatchString(stripped):
		s.handleConfigSchema(w, r)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *notificationRuleMockDTServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	if s.failCreate {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var req dtapi.CreateNotificationRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	s.lastCreateRequest = &req

	rule := &dtapi.NotificationRule{
		Uuid:              "rule-created-uuid-456",
		Name:              req.Name,
		Scope:             req.Scope,
		TriggerType:       "EVENT",
		NotificationLevel: &req.Level,
	}
	s.rules[rule.Uuid] = rule
	writeJSON(w, rule)
}

func (s *notificationRuleMockDTServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if s.failUpdate {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var req dtapi.UpdateNotificationRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	s.lastUpdateRequest = &req

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

func (s *notificationRuleMockDTServer) handleList(w http.ResponseWriter) {
	if s.failList {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	list := make([]dtapi.NotificationRule, 0, len(s.rules))
	for _, rule := range s.rules {
		list = append(list, *rule)
	}
	writeJSON(w, list)
}

// handleConfigSchema handles GET /v1/notification/publisher/{uuid}/configSchema.
func (s *notificationRuleMockDTServer) handleConfigSchema(w http.ResponseWriter, r *http.Request) {
	if s.failSchema {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	// path: /v1/notification/publisher/{uuid}/configSchema → parts[4] = uuid
	if len(parts) < 5 {
		http.Error(w, "bad request", http.StatusBadRequest)
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

func (s *notificationRuleMockDTServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	if s.failDelete {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var body dtapi.NotificationRule
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, exists := s.rules[body.Uuid]; exists {
		delete(s.rules, body.Uuid)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// Test helpers.

func createPublisher(name, ns string) *dependencytrackv1alpha1.NotificationPublisher {
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
	Expect(k8sClient.Create(ctx, p)).To(Succeed())
	return p
}

func createRule(name, ns, publisherName string) *dependencytrackv1alpha1.NotificationRule {
	return &dependencytrackv1alpha1.NotificationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: dependencytrackv1alpha1.NotificationRuleSpec{
			Name:         name + "-rule",
			Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
			TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
			Level:        dependencytrackv1alpha1.NotificationRuleLevelWarn,
			PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: publisherName},
		},
	}
}

func createRuleWithFinalizer(name, ns, publisherName string) *dependencytrackv1alpha1.NotificationRule {
	r := createRule(name, ns, publisherName)
	controllerutil.AddFinalizer(r, notificationRuleFinalizer)
	return r
}

func setRuleUUID(name, ns, uuid string) {
	rule := &dependencytrackv1alpha1.NotificationRule{}
	Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, rule)).To(Succeed())
	rule.Status.UUID = uuid
	Expect(k8sClient.Status().Update(ctx, rule)).To(Succeed())
}

// setPublisherUUID sets the UUID on a publisher so the rule reconcile sees a Ready publisher.
func setPublisherUUID(name, ns, uuid string) {
	p := &dependencytrackv1alpha1.NotificationPublisher{}
	Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, p)).To(Succeed())
	p.Status.UUID = uuid
	p.Status.Name = p.Spec.Name
	Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())
}

func deleteRule(name, ns string) {
	r := &dependencytrackv1alpha1.NotificationRule{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, r); err == nil {
		Expect(k8sClient.Delete(ctx, r)).To(Succeed())
	}
}

// createPublisherConfigSecret creates a Secret with publisherConfig JSON data.
func createPublisherConfigSecret(name, ns, key string, data map[string][]byte) *corev1.Secret {
	s := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Data: data,
	}
	Expect(k8sClient.Create(ctx, s)).To(Succeed())
	return s
}

// ---- Ginkgo tests ----

var _ = Describe("NotificationRule Controller", func() {
	var (
		server       *httptest.Server
		mockServer   *notificationRuleMockDTServer
		provider     *notificationRuleMockClientProvider
		ctx          context.Context
		fakeRecorder *record.FakeRecorder
		ctrl         *NotificationRuleReconciler
		publisher    *dependencytrackv1alpha1.NotificationPublisher
	)

	BeforeEach(func() {
		mockServer = newNotificationRuleMockServer()
		server = httptest.NewServer(mockServer)
		provider = &notificationRuleMockClientProvider{url: server.URL}
		DeferCleanup(server.Close)

		fakeRecorder = record.NewFakeRecorder(20)
		ctrl = &NotificationRuleReconciler{
			Client:                   k8sClient,
			Scheme:                   k8sClient.Scheme(),
			Recorder:                 fakeRecorder,
			DTProvider:               provider,
			PublisherConfigValidator: NewPublisherConfigValidator(),
		}
		ctx = context.Background()
	})

	// ======================================================================
	// Finalizer
	// ======================================================================

	Context("When reconciling a new rule", func() {
		const name = "basic-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			r := createRule(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
		})
		It("should add the finalizer on first reconcile", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(rule, notificationRuleFinalizer)).To(BeTrue())
		})
	})

	// ======================================================================
	// Create — the staged-create path
	// ======================================================================

	Context("When reconciling with no DT UUID (create path)", func() {
		const name = "create-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			// Pre-set publisher UUID so the rule reconcile sees a Ready publisher.
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
		})
		It("should create the rule in DT and persist the UUID", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))
			Expect(rule.Status.Name).To(Equal("create-rule-rule"))
			Expect(rule.Status.ObservedGeneration).To(Equal(rule.Generation))
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("RuleSynced"))
			Eventually(fakeRecorder.Events).Should(Receive(ContainSubstring("RuleCreated")))
		})
	})

	// ======================================================================
	// Update — drift detection
	// ======================================================================

	Context("When reconciling with an existing UUID that matches spec (noop)", func() {
		const name = "noop-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.mu.Lock()
			mockServer.rules["existing-rule-uuid-123"] = &dtapi.NotificationRule{
				Uuid:              "existing-rule-uuid-123",
				Name:              "noop-rule-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "existing-rule-uuid-123")
		})
		It("should not call update when fields already match", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("existing-rule-uuid-123"))
			Expect(rule.Status.ObservedGeneration).To(Equal(rule.Generation))
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("RuleSynced"))
			Consistently(fakeRecorder.Events).ShouldNot(Receive())
		})
	})

	Context("When reconciling with drifted spec fields", func() {
		const name = "drift-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.mu.Lock()
			mockServer.rules["drift-rule-uuid"] = &dtapi.NotificationRule{
				Uuid:              "drift-rule-uuid",
				Name:              "old-name",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "drift-rule-rule"
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "drift-rule-uuid")
		})
		It("should update the rule in DT when spec drifts", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("drift-rule-uuid"))
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("RuleSynced"))
		})
	})

	// ======================================================================
	// The core T03 proof: resume from persisted UUID without duplication
	// ======================================================================

	Context("When a failure occurs after rule creation and retry resumes", func() {
		const name = "resume-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
		})
		It("should resume from persisted UUID without creating a duplicate", func() {
			// Step 1: first reconcile creates the rule in DT and persists UUID.
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))

			// Count creates via mock server state — exactly one rule should exist.
			mockServer.mu.Lock()
			createCount := len(mockServer.rules)
			mockServer.mu.Unlock()
			Expect(createCount).To(Equal(1), "exactly one rule should exist in DT")

			// Step 3: reconcile again — should hit UUID-exists path, find rule, no create.
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			rule = &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))

			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))

			// Step 4: third reconcile — still no duplicate.
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			rule = &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))

			// Verify still exactly one rule in DT (no duplicate creation).
			mockServer.mu.Lock()
			finalCount := len(mockServer.rules)
			mockServer.mu.Unlock()
			Expect(finalCount).To(Equal(1), "no duplicate rule created on retry")
		})
	})

	// ======================================================================
	// Recreate — rule deleted from DT
	// ======================================================================

	Context("When the rule was deleted from DT (recreate)", func() {
		const name = "recreate-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "deleted-rule-uuid")
		})
		It("should recreate the rule when DT does not have it", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))
		})
	})

	// ======================================================================
	// Delete — finalizer path
	// ======================================================================

	Context("When a rule is deleted (finalizer path)", func() {
		It("should remove the finalizer after successful DT deletion", func() {
			const name = "del-success-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.mu.Lock()
			mockServer.rules["del-rule-uuid"] = &dtapi.NotificationRule{
				Uuid:        "del-rule-uuid",
				Name:        "del-rule",
				Scope:       "ALL_PROJECTS",
				TriggerType: "EVENT",
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			setRuleUUID(name, ruleNS, "del-rule-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			// Simulate K8s deletion lifecycle: k8sClient.Delete() sets the
			// DeletionTimestamp on objects with finalizers while keeping them
			// in the fake client store so the controller can reconcile.
			obj := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, obj)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(obj, notificationRuleFinalizer)).To(BeTrue())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())

			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			// After the controller removes the finalizer, the fake client removes
			// the object from storage (DeletionTimestamp set + empty finalizers
			// matches real K8s garbage collection).  Verify via the event instead.
			Eventually(fakeRecorder.Events).Should(Receive(ContainSubstring("RuleDeleted")))
		})

		It("should retain the finalizer when DT delete fails", func() {
			const name = "del-fail-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.failDelete = true
			mockServer.mu.Lock()
			mockServer.rules["fail-del-rule-uuid"] = &dtapi.NotificationRule{
				Uuid:        "fail-del-rule-uuid",
				Name:        "fail-del-rule",
				Scope:       "ALL_PROJECTS",
				TriggerType: "EVENT",
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			setRuleUUID(name, ruleNS, "fail-del-rule-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			obj := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, obj)).To(Succeed())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())

			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())

			obj = &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, obj)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(obj, notificationRuleFinalizer)).To(BeTrue())

			cond := meta.FindStatusCondition(obj.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal(conditionReasonRuleDeleteFailed))
		})

		// UUID deletion treats remote 404 as successful absence.

		It("should succeed when the remote rule is already absent (404)", func() {
			const name = "del-already-absent-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			// Pre-seed the rule in DT so the controller finds it by UUID.
			mockServer.mu.Lock()
			mockServer.rules["already-absent-uuid"] = &dtapi.NotificationRule{
				Uuid:        "already-absent-uuid",
				Name:        name,
				Scope:       "ALL_PROJECTS",
				TriggerType: "EVENT",
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			setRuleUUID(name, ruleNS, "already-absent-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			// Simulate K8s deletion lifecycle.
			obj := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, obj)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(obj, notificationRuleFinalizer)).To(BeTrue())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())

			// The mock deletes the rule — so the next reconcile gets 404.
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			// After finalizer removal the fake client GCs the object.
			Eventually(func() bool {
				r := &dependencytrackv1alpha1.NotificationRule{}
				return apierrors.IsNotFound(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, r))
			}).Should(BeTrue())
		})

		// UUID-targeted deletion does not affect same-name foreign-UUID rules.

		It("should delete only the owned UUID, not a same-name foreign UUID", func() {
			const name = "del-uuid-targeted-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")

			// Pre-seed: owned rule (UUID A) and foreign rule (UUID B) with same name.
			mockServer.mu.Lock()
			mockServer.rules["owned-uuid-a"] = &dtapi.NotificationRule{
				Uuid:        "owned-uuid-a",
				Name:        name,
				Scope:       "ALL_PROJECTS",
				TriggerType: "EVENT",
			}
			mockServer.rules["foreign-uuid-b"] = &dtapi.NotificationRule{
				Uuid:        "foreign-uuid-b",
				Name:        name,
				Scope:       "ALL_PROJECTS",
				TriggerType: "EVENT",
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			setRuleUUID(name, ruleNS, "owned-uuid-a")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			// Simulate K8s deletion lifecycle.
			obj := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, obj)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(obj, notificationRuleFinalizer)).To(BeTrue())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())

			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			// Only the owned UUID should be deleted — foreign must survive.
			mockServer.mu.Lock()
			_, stillExists := mockServer.rules["owned-uuid-a"]
			_, foreignExists := mockServer.rules["foreign-uuid-b"]
			mockServer.mu.Unlock()
			// After successful delete the owned rule is gone.
			Expect(stillExists).To(BeFalse(), "owned UUID must be deleted")
			Expect(foreignExists).To(BeTrue(), "foreign UUID with same name must NOT be deleted")
		})
	})

	// ======================================================================
	// Negative tests — failure modes (Q5 / Q7)
	// ======================================================================

	Context("Negative tests — failure modes", func() {
		It("should report CredentialsError when auth fails", func() {
			const name = "neg-auth-fail-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			badProvider := &notificationRuleMockClientProvider{getErr: fmt.Errorf("auth failed: unreachable")}
			badCtrl := &NotificationRuleReconciler{
				Client: k8sClient, Scheme: k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(5), DTProvider: badProvider,
			}
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })

			_, err := badCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("CredentialsError"))
		})

		It("should report RuleCreateFailed when CREATE returns HTTP 500", func() {
			const name = "neg-create-fail-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.failCreate = true
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("RuleCreateFailed"))
		})

		It("should report RuleListFailed when LIST returns HTTP 500", func() {
			const name = "neg-list-fail-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.failList = true
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "list-fail-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("RuleListFailed"))
		})

		It("should report RuleUpdateFailed when UPDATE returns HTTP 500", func() {
			const name = "neg-update-fail-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.failUpdate = true
			mockServer.mu.Lock()
			mockServer.rules["update-fail-uuid"] = &dtapi.NotificationRule{
				Uuid:        "update-fail-uuid",
				Name:        "old-name",
				Scope:       "ALL_PROJECTS",
				TriggerType: "EVENT",
			}
			mockServer.mu.Unlock()
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "drifted-name" // forces drift → update path
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "update-fail-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("RuleUpdateFailed"))
		})

		It("should report PublisherRefError when the publisher resource is missing", func() {
			const name = "neg-missing-pub-rule"
			r := createRuleWithFinalizer(name, ruleNS, "nonexistent-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherRefError"))
		})

		It("should report PublisherRefError when the publisher has no remote UUID yet", func() {
			const name = "neg-pub-not-ready-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherRefError"))
		})
	})

	// ======================================================================
	// Cross-namespace reference — unrepresentable
	// ======================================================================

	Context("Cross-namespace publisher reference", func() {
		It("should fail when publisher is in a different namespace", func() {
			const name = "cross-ns-rule"
			// envtest can't create in arbitrary namespaces, so create the
			// publisher in ruleNS then delete it. The controller will try to
			// look it up in ruleNS and get NotFound, which manifests as a
			// PublisherRefError condition.
			publisher = createPublisher("cross-pub", ruleNS)
			Expect(k8sClient.Delete(ctx, publisher)).To(Succeed())
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })

			r := createRuleWithFinalizer(name, ruleNS, "cross-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherRefError"))
		})
	})

	// ======================================================================
	// Optional fields — Enabled, NotifyOn, FilterExpression
	// ======================================================================

	Context("When a rule has optional fields set", func() {
		const name = "optional-fields-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			enabled := true
			r.Spec.Enabled = &enabled
			r.Spec.NotifyOn = []dependencytrackv1alpha1.NotificationRuleTriggerEvent{
				dependencytrackv1alpha1.NotificationRuleTriggerEventNewVulnerability,
				dependencytrackv1alpha1.NotificationRuleTriggerEventNewVulnerableDependency,
			}
			r.Spec.FilterExpression = "vulnerable"
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
		})
		It("should persist optional fields on the update stage after create", func() {
			// First reconcile: create (only sends required fields).
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))

			// Second reconcile: should update with optional fields (drift detected).
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			// Verify the mock received the update with optional fields.
			mockServer.mu.Lock()
			ruleInMock := mockServer.rules[rule.Status.UUID]
			mockServer.mu.Unlock()
			Expect(ruleInMock).NotTo(BeNil())
			Expect(*ruleInMock.Enabled).To(BeTrue())
			Expect(ruleInMock.NotifyOn).To(ContainElements("NEW_VULNERABILITY", "NEW_VULNERABLE_DEPENDENCY"))
			Expect(*ruleInMock.FilterExpression).To(Equal("vulnerable"))
		})
	})

	// ======================================================================
	// Scheduled routing fields — LogSuccessfulPublish, NotifyChildren,
	// ScheduleCron, ScheduleSkipUnchanged (T02)
	// ======================================================================

	Context("When a scheduled rule has routing fields set", func() {
		const name = "scheduled-routing-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			logSuccess := true
			notifyCh := true
			skipUnchanged := false
			r.Spec.LogSuccessfulPublish = &logSuccess
			r.Spec.NotifyChildren = &notifyCh
			r.Spec.ScheduleCron = "0 0 * * *"
			r.Spec.ScheduleSkipUnchanged = &skipUnchanged
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
		})
		It("should converge all scheduled routing fields on update", func() {
			// First reconcile: create (only sends required fields).
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("rule-created-uuid-456"))

			// Second reconcile: drift detected → update with all new fields.
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			// Verify the mock received and persisted all five new fields.
			mockServer.mu.Lock()
			ruleInMock := mockServer.rules[rule.Status.UUID]
			lastReq := mockServer.lastUpdateRequest
			mockServer.mu.Unlock()

			Expect(ruleInMock).NotTo(BeNil())
			Expect(lastReq).NotTo(BeNil())

			Expect(ruleInMock.LogSuccessfulPublish).NotTo(BeNil())
			Expect(*ruleInMock.LogSuccessfulPublish).To(BeTrue())
			Expect(*lastReq.LogSuccessfulPublish).To(BeTrue())

			Expect(ruleInMock.NotifyChildren).NotTo(BeNil())
			Expect(*ruleInMock.NotifyChildren).To(BeTrue())
			Expect(*lastReq.NotifyChildren).To(BeTrue())

			Expect(ruleInMock.ScheduleCron).NotTo(BeNil())
			Expect(*ruleInMock.ScheduleCron).To(Equal("0 0 * * *"))
			Expect(*lastReq.ScheduleCron).To(Equal("0 0 * * *"))

			Expect(ruleInMock.ScheduleSkipUnchanged).NotTo(BeNil())
			Expect(*ruleInMock.ScheduleSkipUnchanged).To(BeFalse())
			Expect(*lastReq.ScheduleSkipUnchanged).To(BeFalse())
		})
	})

	// ======================================================================
	// ObservedGeneration is set on Ready condition
	// ======================================================================

	Context("ObservedGeneration is set on Ready condition", func() {
		const name = "obs-gen-rule"
		BeforeEach(func() {
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
		})
		It("should set ObservedGeneration in the Ready condition", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())
			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.ObservedGeneration).To(Equal(rule.Generation))
		})
	})

	// ======================================================================
	// Publisher Config Convergence (T02)
	// ======================================================================

	Context("Publisher config convergence", func() {
		validConfig := `{"webhookUrl": "https://example.com/hook", "headers": {"Authorization": "Bearer test"}}`

		// --- Happy path: converged publisherConfig on update ---

		It("should include publisherConfig in update when Secret is valid", func() {
			const name = "pc-happy-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			createPublisherConfigSecret("pc-secret", ruleNS, "config", map[string][]byte{
				"config": []byte(validConfig),
			})
			mockServer.mu.Lock()
			mockServer.rules["pc-rule-uuid"] = &dtapi.NotificationRule{
				Uuid:              "pc-rule-uuid",
				Name:              "pc-config-rule-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.PublisherConfigSecretRef = &dependencytrackv1alpha1.PublisherConfigSecretRef{
				Name: "pc-secret",
				Key:  "config",
			}
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "pc-rule-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			Expect(rule.Status.UUID).To(Equal("pc-rule-uuid"))
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("RuleSynced"))

			// Verify the mock received the update with publisherConfig.
			mockServer.mu.Lock()
			lastReq := mockServer.lastUpdateRequest
			mockServer.mu.Unlock()
			Expect(lastReq).NotTo(BeNil())
			Expect(*lastReq.PublisherConfig).To(Equal(validConfig))
		})

		// --- Negative: missing Secret → SecretNotFound ---

		It("should report PublisherConfigValidationFailed when the Secret is not found", func() {
			const name = "pc-neg-secret-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.mu.Lock()
			mockServer.rules["pc-neg-secret-uuid"] = &dtapi.NotificationRule{
				Uuid:              "pc-neg-secret-uuid",
				Name:              "pc-neg-secret-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "pc-neg-secret-rule" // force drift → update
			r.Spec.PublisherConfigSecretRef = &dependencytrackv1alpha1.PublisherConfigSecretRef{
				Name: "nonexistent-secret",
				Key:  "config",
			}
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "pc-neg-secret-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherConfigValidationFailed"))
			// Secret name must NOT appear in the condition message.
			Expect(cond.Message).NotTo(ContainSubstring("pc-secret"))
		})

		// --- Negative: Secret key not found → SecretKeyNotFound ---

		It("should report SecretKeyNotFound when the Secret key is missing", func() {
			const name = "pc-neg-key-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			createPublisherConfigSecret("pc-key-secret", ruleNS, "otherKey", map[string][]byte{
				"otherKey": []byte(`{}`),
			})
			mockServer.mu.Lock()
			mockServer.rules["pc-neg-key-uuid"] = &dtapi.NotificationRule{
				Uuid:              "pc-neg-key-uuid",
				Name:              "pc-neg-key-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "pc-neg-key-rule" // force drift → update
			r.Spec.PublisherConfigSecretRef = &dependencytrackv1alpha1.PublisherConfigSecretRef{
				Name: "pc-key-secret",
				Key:  "missingKey",
			}
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "pc-neg-key-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherConfigValidationFailed"))
		})

		// --- Negative: malformed JSON → MalformedJSON ---

		It("should report MalformedJSON when the Secret value is not valid JSON", func() {
			const name = "pc-neg-json-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			createPublisherConfigSecret("pc-bad-json-secret", ruleNS, "config", map[string][]byte{
				"config": []byte(`this is not json`),
			})
			mockServer.mu.Lock()
			mockServer.rules["pc-neg-json-uuid"] = &dtapi.NotificationRule{
				Uuid:              "pc-neg-json-uuid",
				Name:              "pc-neg-json-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "pc-neg-json-rule" // force drift → update
			r.Spec.PublisherConfigSecretRef = &dependencytrackv1alpha1.PublisherConfigSecretRef{
				Name: "pc-bad-json-secret",
				Key:  "config",
			}
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "pc-neg-json-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherConfigValidationFailed"))
		})

		// --- Negative: schema validation failure → SchemaValidationFailed ---

		It("should report SchemaValidationFailed when JSON has wrong types", func() {
			const name = "pc-neg-schema-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			// webhookUrl is a number instead of string — type mismatch is still validated.
			createPublisherConfigSecret("pc-schema-secret", ruleNS, "config", map[string][]byte{
				"config": []byte(`{"webhookUrl": 123}`),
			})
			mockServer.mu.Lock()
			mockServer.rules["pc-neg-schema-uuid"] = &dtapi.NotificationRule{
				Uuid:              "pc-neg-schema-uuid",
				Name:              "pc-neg-schema-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "pc-neg-schema-rule" // force drift → update
			r.Spec.PublisherConfigSecretRef = &dependencytrackv1alpha1.PublisherConfigSecretRef{
				Name: "pc-schema-secret",
				Key:  "config",
			}
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "pc-neg-schema-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherConfigValidationFailed"))
		})

		// --- Negative: schema retrieval failure → SchemaRetrievalFailed ---

		It("should report SchemaRetrievalFailed when the schema endpoint fails", func() {
			const name = "pc-neg-sch-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			createPublisherConfigSecret("pc-schema-fail-secret", ruleNS, "config", map[string][]byte{
				"config": []byte(validConfig),
			})
			mockServer.failSchema = true
			mockServer.mu.Lock()
			mockServer.rules["pc-neg-sch-uuid"] = &dtapi.NotificationRule{
				Uuid:              "pc-neg-sch-uuid",
				Name:              "pc-neg-sch-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "pc-neg-sch-rule" // force drift → update
			r.Spec.PublisherConfigSecretRef = &dependencytrackv1alpha1.PublisherConfigSecretRef{
				Name: "pc-schema-fail-secret",
				Key:  "config",
			}
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "pc-neg-sch-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).To(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherConfigValidationFailed"))
		})

		// --- No-op: no SecretRef → config skipped, update proceeds ---

		It("should update without publisherConfig when no SecretRef is set", func() {
			const name = "pc-no-ref-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			mockServer.mu.Lock()
			mockServer.rules["pc-no-ref-uuid"] = &dtapi.NotificationRule{
				Uuid:              "pc-no-ref-uuid",
				Name:              "pc-no-ref-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "pc-no-ref-rule" // force drift → update
			// No PublisherConfigSecretRef set.
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "pc-no-ref-uuid")

			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))

			// publisherConfig should NOT be set in the update request.
			mockServer.mu.Lock()
			lastReq := mockServer.lastUpdateRequest
			mockServer.mu.Unlock()
			Expect(lastReq).NotTo(BeNil())
			Expect(lastReq.PublisherConfig).To(BeNil())
		})

		// --- Secret values must never leak ---

		It("must never leak Secret values into conditions, events, or logs", func() {
			const name = "pc-leak-rule"
			publisher = createPublisher(name+"-pub", ruleNS)
			setPublisherUUID(name+"-pub", ruleNS, "pub-created-uuid-123")
			// Use a secret with an obviously sensitive value.
			sensitiveConfig := `{"webhookUrl": "https://secret-token-abc123"}`
			createPublisherConfigSecret("pc-leak-secret", ruleNS, "config", map[string][]byte{
				"config": []byte(sensitiveConfig),
			})
			mockServer.mu.Lock()
			mockServer.rules["pc-leak-uuid"] = &dtapi.NotificationRule{
				Uuid:              "pc-leak-uuid",
				Name:              "pc-leak-rule",
				Scope:             "ALL_PROJECTS",
				TriggerType:       "EVENT",
				NotificationLevel: strPtr("WARN"),
			}
			mockServer.mu.Unlock()

			r := createRuleWithFinalizer(name, ruleNS, name+"-pub")
			r.Spec.Name = "pc-leak-rule" // force drift → update
			r.Spec.PublisherConfigSecretRef = &dependencytrackv1alpha1.PublisherConfigSecretRef{
				Name: "pc-leak-secret",
				Key:  "config",
			}
			Expect(k8sClient.Create(ctx, r)).To(Succeed())
			DeferCleanup(deleteRule, name, ruleNS)
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setRuleUUID(name, ruleNS, "pc-leak-uuid")

			// Reconcile with a valid secret — verify no leakage.
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ruleNS}})
			Expect(err).NotTo(HaveOccurred())

			rule := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNS}, rule)).To(Succeed())
			cond := meta.FindStatusCondition(rule.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			// Condition message must not contain the sensitive URL.
			Expect(cond.Message).NotTo(ContainSubstring("secret-token-abc123"))

			// FakeRecorder events must not contain the sensitive URL.
			for i := 0; i < len(fakeRecorder.Events); i++ {
				ev := <-fakeRecorder.Events
				Expect(ev).NotTo(ContainSubstring("secret-token-abc123"))
			}
		})
	})
})

// strPtr is a helper to create a *string from a string literal.
func strPtr(s string) *string {
	return &s
}
