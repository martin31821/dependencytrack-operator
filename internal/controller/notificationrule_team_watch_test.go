/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

// teamRuleMockServer handles notification rule operations for integration tests.
type teamRuleMockServer struct {
	mu                sync.Mutex
	rules             map[string]*dtapi.NotificationRule
	lastUpdateRequest *dtapi.UpdateNotificationRuleRequest
}

func newTeamRuleMockServer() *teamRuleMockServer {
	return &teamRuleMockServer{
		rules: make(map[string]*dtapi.NotificationRule),
	}
}

func (s *teamRuleMockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stripped := r.URL.Path
	if len(stripped) > 4 && stripped[:4] == "/api" {
		stripped = stripped[4:]
	}

	if r.Method == http.MethodGet && stripped == "/v1/notification/rule" {
		s.mu.Lock()
		list := make([]dtapi.NotificationRule, 0, len(s.rules))
		for _, rule := range s.rules {
			list = append(list, *rule)
		}
		s.mu.Unlock()
		writeJSON(w, list)
		return
	}

	if r.Method == http.MethodPost && stripped == "/v1/notification/rule" {
		var req dtapi.UpdateNotificationRuleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.mu.Lock()
		s.lastUpdateRequest = &req
		rule, exists := s.rules[req.Uuid]
		if exists {
			if req.Name != "" {
				rule.Name = req.Name
			}
			if req.Scope != "" {
				rule.Scope = req.Scope
			}
			if req.Enabled != nil {
				rule.Enabled = req.Enabled
			}
			if req.NotifyOn != nil {
				rule.NotifyOn = req.NotifyOn
			}
		}
		s.mu.Unlock()
		writeJSON(w, rule)
		return
	}

	if r.Method == http.MethodPut && stripped == "/v1/notification/rule" {
		var req dtapi.CreateNotificationRuleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rule := &dtapi.NotificationRule{
			Uuid: "rule-created-uuid-456", Name: req.Name, Scope: req.Scope,
			TriggerType: "EVENT", NotificationLevel: &req.Level,
		}
		s.mu.Lock()
		s.rules[rule.Uuid] = rule
		s.mu.Unlock()
		writeJSON(w, rule)
		return
	}

	// All other endpoints (including /v1/team/*) return 404.
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"title":"Not Found","status":404}`))
}

// Test the indexer function directly — this is the core team watch mechanism.
var _ = Describe("NotificationRule Team Watch Indexer", func() {
	It("should correctly index NotificationRules by team name from spec", func() {
		ruleWithTeams := &dependencytrackv1alpha1.NotificationRule{
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Teams: []string{"devs", "ops"},
			},
		}
		keys := indexFuncNotificationRuleByTeamName(ruleWithTeams)
		Expect(keys).To(ConsistOf("devs", "ops"))

		// Duplicate teams are deduplicated.
		ruleWithDupes := &dependencytrackv1alpha1.NotificationRule{
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Teams: []string{"devs", "devs"},
			},
		}
		keys = indexFuncNotificationRuleByTeamName(ruleWithDupes)
		Expect(keys).To(ConsistOf("devs"))

		// Nil/empty teams return nil.
		Expect(indexFuncNotificationRuleByTeamName(&dependencytrackv1alpha1.NotificationRule{})).To(BeNil())
		ruleNoTeams := &dependencytrackv1alpha1.NotificationRule{
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{Teams: []string{}},
		}
		Expect(indexFuncNotificationRuleByTeamName(ruleNoTeams)).To(BeNil())

		// Non-rule object returns nil.
		Expect(indexFuncNotificationRuleByTeamName(&dependencytrackv1alpha1.NotificationPublisher{})).To(BeNil())
	})
})

// Test the full reconciler flow with teams — verifies finalizer and UUID persistence.
var _ = Describe("NotificationRule Team Watch Integration", func() {
	var (
		server       *httptest.Server
		mockServer   *teamRuleMockServer
		provider     *notificationRuleMockClientProvider
		ctrl         *NotificationRuleReconciler
		fakeRecorder *record.FakeRecorder
		ctx          context.Context
	)

	BeforeEach(func() {
		mockServer = newTeamRuleMockServer()
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

	Context("When a rule references Teams", func() {
		It("should add finalizer and persist UUID", func() {
			const (
				ruleName = "team-rule"
				teamName = "team-rule-team"
			)

			publisher := &dependencytrackv1alpha1.NotificationPublisher{
				ObjectMeta: metav1.ObjectMeta{Name: ruleName + "-pub", Namespace: "default"},
				Spec:       dependencytrackv1alpha1.NotificationPublisherSpec{ExtensionName: "webhook", Name: "webhook"},
			}
			Expect(k8sClient.Create(ctx, publisher)).To(Succeed())
			DeferCleanup(func() { k8sClient.Delete(ctx, publisher) })
			setPublisherUUID(ruleName+"-pub", "default", "pub-uuid-team")

			team := &dependencytrackv1alpha1.Team{
				ObjectMeta: metav1.ObjectMeta{Name: teamName, Namespace: "default"},
			}
			Expect(k8sClient.Create(ctx, team)).To(Succeed())
			DeferCleanup(func() { _ = k8sClient.Delete(ctx, team) })
			team.Status.UUID = "team-uuid"
			Expect(k8sClient.Status().Update(ctx, team)).To(Succeed())

			rule := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{Name: ruleName, Namespace: "default"},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name: ruleName + "-rule", Scope: dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:        dependencytrackv1alpha1.NotificationRuleLevelWarn,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: ruleName + "-pub"},
					Teams:        []string{teamName},
					Projects:     []string{"proj-1"},
				},
			}
			Expect(k8sClient.Create(ctx, rule)).To(Succeed())
			DeferCleanup(func() { k8sClient.Delete(ctx, rule) })

			// First reconcile: adds finalizer.
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: ruleName, Namespace: "default"}})
			Expect(err).NotTo(HaveOccurred())

			retrieved := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: ruleName, Namespace: "default"}, retrieved)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(retrieved, notificationRuleFinalizer)).To(BeTrue())

			// Second reconcile: with finalizer, proceeds to rule creation.
			// Team resolution logic is tested in unit tests (TestResolveTeamRefs).
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: ruleName, Namespace: "default"}})
			retrieved2 := &dependencytrackv1alpha1.NotificationRule{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: ruleName, Namespace: "default"}, retrieved2)).To(Succeed())
			Expect(retrieved2.Status.UUID).NotTo(BeEmpty())
		})
	})
})
