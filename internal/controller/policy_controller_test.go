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
	"k8s.io/apimachinery/pkg/api/errors"
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
	testNS                  = "default"
	testSeverityHigh        = "HIGH"
	testSeverityCritical    = "CRITICAL"
	testOfflinePolicyName   = "offline-policy"
	testConvergedPolicyName = "Converged Policy"
)

// --- mock ClientProviderInterface for tests ---

// mockClientProvider implements ClientProviderInterface without needing a K8s client.
type mockClientProvider struct {
	url    string
	mu     sync.Mutex
	api    *dtapi.APIClient
	getErr error // when non-nil, Get returns this error (simulates auth failure)
}

func (p *mockClientProvider) Get(ctx context.Context) (context.Context, *dtapi.APIClient, error) {
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

// --- httptest-based mock server for DependencyTrack API ---

type mockDTServer struct {
	mu                  sync.Mutex
	policy              *dtapi.Policy
	policies            map[string]*dtapi.Policy // uuid -> policy, used for GET /v1/policy listing
	failCreatePolicy    bool                     // when true, /v1/policy PUT returns 500
	failCreateCondition bool                     // when true, /v1/policy/{uuid}/condition PUT returns 500
	failListPolicies    bool                     // when true, /v1/policy GET returns 500
	failGetPolicy       bool                     // when true, /v1/policy/{uuid} GET returns 500
	failGetPolicyUUID   string                   // if set, return 500 on this specific UUID GET
	failDeletePolicy    bool                     // when true, DELETE /v1/policy/{uuid} returns 500
}

const (
	apiPathPrefix           = "/api"
	policyPath              = "/v1/policy"
	policyPathPrefix        = policyPath + "/"
	policyConditionPath     = policyPath + "/condition"
	policyConditionPrefix   = policyConditionPath + "/"
	validateCredentialsPath = "/v1/user/validate-credentials"
)

func (s *mockDTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path, _ := strings.CutPrefix(r.URL.Path, apiPathPrefix)

	switch {
	case path == policyPath:
		s.handlePolicyCollection(w, r)
	case path == policyConditionPath || strings.HasPrefix(path, policyConditionPrefix):
		s.handlePolicyCondition(w, r, path)
	case strings.HasPrefix(path, policyPathPrefix):
		s.handlePolicy(w, r, path)
	case path == validateCredentialsPath:
		writeJSON(w, "mock-token")
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *mockDTServer) handlePolicyCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listPolicies(w)
	case http.MethodPut:
		s.createPolicy(w, r)
	case http.MethodPost:
		s.updatePolicy(w, r)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *mockDTServer) listPolicies(w http.ResponseWriter) {
	if s.failListPolicies {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	all := make([]dtapi.Policy, 0, len(s.policies)+1)
	for _, policy := range s.policies {
		all = append(all, *policy)
	}
	if s.policy != nil {
		all = append(all, *s.policy)
	}
	writeJSON(w, all)
}

func (s *mockDTServer) createPolicy(w http.ResponseWriter, r *http.Request) {
	if s.failCreatePolicy {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var policy dtapi.Policy
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	policy.Uuid = "created-uuid-123"
	if policy.Operator == "" {
		policy.Operator = "OR"
	}
	s.policy = &policy
	writeJSON(w, s.policy)
}

func (s *mockDTServer) updatePolicy(w http.ResponseWriter, r *http.Request) {
	var policy dtapi.Policy
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if s.policy != nil {
		s.policy.Name = policy.Name
		if policy.Operator != "" {
			s.policy.Operator = policy.Operator
		}
		if policy.ViolationState != "" {
			s.policy.ViolationState = policy.ViolationState
		}
	}
	writeJSON(w, s.policy)
}

func (s *mockDTServer) handlePolicy(w http.ResponseWriter, r *http.Request, path string) {
	if strings.HasSuffix(path, "/condition") && r.Method == http.MethodPut {
		s.createCondition(w, r, path)
		return
	}

	uuid, _ := strings.CutPrefix(path, policyPathPrefix)
	switch r.Method {
	case http.MethodGet:
		s.getPolicy(w, uuid)
	case http.MethodDelete:
		s.deletePolicy(w, uuid)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *mockDTServer) getPolicy(w http.ResponseWriter, uuid string) {
	if s.failGetPolicy || (s.failGetPolicyUUID != "" && s.failGetPolicyUUID == uuid) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if s.policy == nil || s.policy.GetUuid() != uuid {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, s.policy)
}

func (s *mockDTServer) deletePolicy(w http.ResponseWriter, uuid string) {
	if s.failDeletePolicy {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if s.policy == nil || s.policy.GetUuid() != uuid {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s.policy = nil
	w.WriteHeader(http.StatusNoContent)
}

func (s *mockDTServer) createCondition(w http.ResponseWriter, r *http.Request, path string) {
	if s.failCreateCondition {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	conditionPath, _ := strings.CutPrefix(path, policyPathPrefix)
	uuid := strings.TrimSuffix(conditionPath, "/condition")
	var condition dtapi.PolicyCondition
	if err := json.NewDecoder(r.Body).Decode(&condition); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if condition.Uuid == "" {
		condition.Uuid = "cond-uuid-" + uuid
	}
	writeJSON(w, condition)
}

func (s *mockDTServer) handlePolicyCondition(w http.ResponseWriter, r *http.Request, path string) {
	if strings.HasPrefix(path, policyConditionPrefix) && r.Method == http.MethodDelete {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if path != policyConditionPath || r.Method != http.MethodPost {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var condition dtapi.PolicyCondition
	if err := json.NewDecoder(r.Body).Decode(&condition); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	writeJSON(w, condition)
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

var _ = Describe("Policy Controller", func() {
	var (
		mockServer   *httptest.Server
		mockDT       *mockDTServer
		ctx          context.Context
		fakeRecorder *record.FakeRecorder
		controller   *PolicyReconciler
	)

	BeforeEach(func() {
		mockDT = &mockDTServer{policy: nil}
		mockServer = httptest.NewServer(mockDT)
		DeferCleanup(mockServer.Close)

		fakeRecorder = record.NewFakeRecorder(10)
		controller = &PolicyReconciler{
			Client:   k8sClient,
			Scheme:   k8sClient.Scheme(),
			Recorder: fakeRecorder,
			DTProvider: &mockClientProvider{
				url: mockServer.URL,
			},
		}
		ctx = context.Background()
	})

	Context("When reconciling a resource", func() {
		const resourceName = "test-policy"

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: testNS,
		}
		policy := &dependencytrackv1alpha1.Policy{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Policy")
			err := k8sClient.Get(ctx, typeNamespacedName, policy)
			if err != nil && errors.IsNotFound(err) {
				resource := &dependencytrackv1alpha1.Policy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: testNS,
					},
					Spec: dependencytrackv1alpha1.PolicySpec{
						Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
						Name:           "Test Policy",
						ViolationState: dependencytrackv1alpha1.ViolationStateFail,
						Conditions: []dependencytrackv1alpha1.PolicyCondition{
							{
								Subject:  dependencytrackv1alpha1.PolicyConditionSubjectSeverity,
								Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs,
								Value:    testSeverityHigh,
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())
			By("Cleanup the specific resource instance Policy")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			_, err := controller.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify the finalizer was added.
			policy = &dependencytrackv1alpha1.Policy{}
			err = k8sClient.Get(ctx, typeNamespacedName, policy)
			Expect(err).NotTo(HaveOccurred())
			Expect(policy.Finalizers).To(ContainElement(policyFinalizer))
		})
	})

	Context("When reconciling with a new DT policy creation", func() {
		const resourceName = "test-policy-new"

		ctx := context.Background()

		BeforeEach(func() {
			resource := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       resourceName,
					Namespace:  testNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "New Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateFail,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{
							Subject:  dependencytrackv1alpha1.PolicyConditionSubjectLicense,
							Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs,
							Value:    "GPL-3.0",
						},
						{
							Subject:  dependencytrackv1alpha1.PolicyConditionSubjectSeverity,
							Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs,
							Value:    testSeverityCritical,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, resource)).To(Succeed())
		})

		AfterEach(func() {
			resource := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should create a new DT policy and conditions", func() {
			// Reset server to have no existing policy (nil triggers create).
			mockDT.mu.Lock()
			mockDT.policy = nil
			mockDT.mu.Unlock()

			// First reconcile adds finalizer (envtest strips finalizers from Create).
			_, err := controller.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      resourceName,
					Namespace: testNS,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Second reconcile does the actual DT API work.
			_, err = controller.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      resourceName,
					Namespace: testNS,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify status was updated.
			policy := &dependencytrackv1alpha1.Policy{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, policy)
			Expect(err).NotTo(HaveOccurred())
			Expect(policy.Status.UUID).To(Equal("created-uuid-123"))
			Expect(policy.Status.Name).To(Equal("New Policy"))

			cond := meta.FindStatusCondition(policy.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("PolicySynced"))
		})
	})

	Context("When reconciling a policy that no longer exists in DT", func() {
		const resourceName = "test-policy-orphan"

		ctx := context.Background()

		BeforeEach(func() {
			resource := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       resourceName,
					Namespace:  testNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Orphan Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{
							Subject:  dependencytrackv1alpha1.PolicyConditionSubjectCPE,
							Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs,
							Value:    "app:myapp",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, resource)).To(Succeed())
		})

		AfterEach(func() {
			resource := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should recreate the policy when DT returns 404", func() {
			// Pre-populate status UUID to simulate orphaned state.
			policy := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, policy)
			Expect(err).NotTo(HaveOccurred())
			policy.Status.UUID = "orphaned-uuid-999"
			Expect(k8sClient.Status().Update(ctx, policy)).To(Succeed())

			// Server has no matching policy — triggers 404 path → recreation.

			_, err = controller.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      resourceName,
					Namespace: testNS,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify the policy was recreated with new UUID.
			policy = &dependencytrackv1alpha1.Policy{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, policy)
			Expect(err).NotTo(HaveOccurred())
			Expect(policy.Status.UUID).To(Equal("created-uuid-123"))
			Expect(policy.Status.Name).To(Equal("Orphan Policy"))
		})
	})

	Context("policyConditionToDT conversion", func() {
		It("should pass native subject, IS operator, and value through unchanged", func() {
			specCond := dependencytrackv1alpha1.PolicyCondition{
				Subject:  dependencytrackv1alpha1.PolicyConditionSubjectSeverity,
				Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs,
				Value:    "CRITICAL",
			}
			dtCond := policyConditionToDT(specCond)
			Expect(dtCond.GetSubject()).To(Equal("SEVERITY"))
			Expect(dtCond.GetOperator()).To(Equal("IS"))
			Expect(dtCond.GetValue()).To(Equal("CRITICAL"))
			Expect(dtCond.GetViolationType()).To(BeEmpty())
		})

		It("should pass IS_NOT through unchanged", func() {
			specCond := dependencytrackv1alpha1.PolicyCondition{
				Subject:  dependencytrackv1alpha1.PolicyConditionSubjectLicense,
				Operator: dependencytrackv1alpha1.PolicyConditionOperatorIsNot,
				Value:    "GPL-3.0-only",
			}
			dtCond := policyConditionToDT(specCond)
			Expect(dtCond.GetSubject()).To(Equal("LICENSE"))
			Expect(dtCond.GetOperator()).To(Equal("IS_NOT"))
			Expect(dtCond.GetValue()).To(Equal("GPL-3.0-only"))
			Expect(dtCond.GetViolationType()).To(BeEmpty())
		})
	})

	Context("boolPtr helper", func() {
		It("should return a pointer to true", func() {
			b := boolPtr(true)
			Expect(*b).To(BeTrue())
		})

		It("should return a pointer to false", func() {
			b := boolPtr(false)
			Expect(*b).To(BeFalse())
		})
	})

	Context("when a same-name DT policy already exists (unowned)", func() {
		const resourceName = "test-policy-conflict"

		ctx := context.Background()

		BeforeEach(func() {
			resource := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       resourceName,
					Namespace:  testNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Shared Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{
							Subject:  dependencytrackv1alpha1.PolicyConditionSubjectSeverity,
							Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs,
							Value:    testSeverityCritical,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, resource)).To(Succeed())
		})

		AfterEach(func() {
			resource := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should refuse to create when a same-name DT policy exists and report OwnershipConflict", func() {
			// Pre-populate the DT mock server with a policy that shares the same name
			// but has a different UUID (simulating an unowned external policy).
			mockDT.mu.Lock()
			mockDT.policy = &dtapi.Policy{
				Name:           "Shared Policy",
				Uuid:           "external-uuid-abc",
				Global:         boolPtr(true),
				Operator:       string(dependencytrackv1alpha1.PolicyOperatorAny),
				ViolationState: string(dependencytrackv1alpha1.ViolationStateWarn),
			}
			mockDT.mu.Unlock()

			_, err := controller.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      resourceName,
					Namespace: testNS,
				},
			})
			// Reconcile returns nil error — the conflict is captured in status.
			Expect(err).NotTo(HaveOccurred())

			// Verify the Ready condition is False with OwnershipConflict reason.
			policy := &dependencytrackv1alpha1.Policy{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, policy)
			Expect(err).NotTo(HaveOccurred())

			cond := meta.FindStatusCondition(policy.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("OwnershipConflict"))
			Expect(cond.Message).To(ContainSubstring("Shared Policy"))
			Expect(cond.Message).To(ContainSubstring("already owned"))
			// Status UUID must NOT be set — ownership was refused.
			Expect(policy.Status.UUID).To(BeEmpty())
		})
	})

	Context("when a K8s-owned UUID is already set", func() {
		const resourceName = "test-policy-uuid-owned"

		ctx := context.Background()

		BeforeEach(func() {
			resource := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       resourceName,
					Namespace:  testNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Owned Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateInfo,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{
							Subject:  dependencytrackv1alpha1.PolicyConditionSubjectLicense,
							Operator: dependencytrackv1alpha1.PolicyConditionOperatorIsNot,
							Value:    "AGPL-3.0",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, resource)).To(Succeed())
		})

		AfterEach(func() {
			resource := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should use the recorded UUID and not list policies for a UUID-owned resource", func() {
			// Pre-populate status UUID so getOrCreateDTPPolicy goes the UUID-get path
			policy := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, policy)
			Expect(err).NotTo(HaveOccurred())
			policy.Status.UUID = "created-uuid-123"
			Expect(k8sClient.Status().Update(ctx, policy)).To(Succeed())

			_, err = controller.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      resourceName,
					Namespace: testNS,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			policy = &dependencytrackv1alpha1.Policy{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: testNS}, policy)
			Expect(err).NotTo(HaveOccurred())
			Expect(policy.Status.UUID).To(Equal("created-uuid-123"))
			Expect(policy.Status.Name).To(Equal("Owned Policy"))

			cond := meta.FindStatusCondition(policy.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("PolicySynced"))
		})
	})

	Context("when a UUID-owned policy has drifted fields", func() {
		ctx := context.Background()

		It("should converge drifted name, operator, and violation state on the same remote UUID", func() {
			const driftResourceName = "test-policy-drift-converge"
			const driftNS = testNS

			// Create the K8s Policy CR with drifted spec relative to the pre-existing DT policy.
			driftResource := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       driftResourceName,
					Namespace:  driftNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAll,
					Name:           testConvergedPolicyName,
					ViolationState: dependencytrackv1alpha1.ViolationStateFail,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: "0.1"},
					},
				},
			}
			Expect(k8sClient.Create(ctx, driftResource)).To(Succeed())
			DeferCleanup(func() {
				res := &dependencytrackv1alpha1.Policy{}
				if e := k8sClient.Get(ctx, types.NamespacedName{Name: driftResourceName, Namespace: driftNS}, res); e == nil {
					_ = k8sClient.Delete(ctx, res)
				}
			})

			// Pre-populate the mock DT server with a policy that has the same UUID
			// but drifted fields (different name, operator, and ViolationState) — simulating
			// out-of-band manual edits in the DependencyTrack UI.
			mockDT.mu.Lock()
			mockDT.policy = &dtapi.Policy{
				Name:           "Old Stale Name",
				Uuid:           "converge-uuid-456",
				Global:         boolPtr(true),
				Operator:       string(dependencytrackv1alpha1.PolicyOperatorAny),
				ViolationState: string(dependencytrackv1alpha1.ViolationStateFail),
			}
			mockDT.mu.Unlock()

			// Pre-set K8s status UUID to match the remote policy.
			k8sPolicy := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: driftResourceName, Namespace: driftNS}, k8sPolicy)
			Expect(err).NotTo(HaveOccurred())
			k8sPolicy.Status.UUID = "converge-uuid-456"
			Expect(k8sClient.Status().Update(ctx, k8sPolicy)).To(Succeed())

			// Reconcile — the controller should detect name + violationState drift
			// and issue an update to converge.
			_, err = controller.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      driftResourceName,
					Namespace: driftNS,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify the K8s status is converged.
			k8sPolicy = &dependencytrackv1alpha1.Policy{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: driftResourceName, Namespace: driftNS}, k8sPolicy)
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sPolicy.Status.UUID).To(Equal("converge-uuid-456"))
			Expect(k8sPolicy.Status.Name).To(Equal(testConvergedPolicyName))

			cond := meta.FindStatusCondition(k8sPolicy.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("PolicySynced"))

			// Verify the mock server received the converged fields.
			mockDT.mu.Lock()
			defer mockDT.mu.Unlock()
			Expect(mockDT.policy.GetName()).To(Equal(testConvergedPolicyName))
			Expect(mockDT.policy.GetOperator()).To(Equal(string(dependencytrackv1alpha1.PolicyOperatorAll)))
			Expect(mockDT.policy.GetViolationState()).To(Equal(string(dependencytrackv1alpha1.ViolationStateFail)))
			Expect(mockDT.policy.GetUuid()).To(Equal("converge-uuid-456"))
		})

		It("should not send an update when fields already match", func() {
			const noopResourceName = "test-policy-drift-noop"
			const noopNS = testNS

			// Create the K8s Policy CR with matching spec — no drift.
			noopResource := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       noopResourceName,
					Namespace:  noopNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           testConvergedPolicyName,
					ViolationState: dependencytrackv1alpha1.ViolationStateFail,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: "0.1"},
					},
				},
			}
			Expect(k8sClient.Create(ctx, noopResource)).To(Succeed())
			DeferCleanup(func() {
				res := &dependencytrackv1alpha1.Policy{}
				if e := k8sClient.Get(ctx, types.NamespacedName{Name: noopResourceName, Namespace: noopNS}, res); e == nil {
					_ = k8sClient.Delete(ctx, res)
				}
			})

			// Pre-populate the mock DT server with a policy that already matches
			// the K8s spec — no drift.
			mockDT.mu.Lock()
			mockDT.policy = &dtapi.Policy{
				Name:           testConvergedPolicyName,
				Uuid:           "converge-uuid-noop",
				Global:         boolPtr(true),
				Operator:       string(dependencytrackv1alpha1.PolicyOperatorAny),
				ViolationState: string(dependencytrackv1alpha1.ViolationStateFail),
			}
			mockDT.mu.Unlock()

			// Pre-set K8s status UUID to match.
			k8sPolicy := &dependencytrackv1alpha1.Policy{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: noopResourceName, Namespace: noopNS}, k8sPolicy)
			Expect(err).NotTo(HaveOccurred())
			k8sPolicy.Status.UUID = "converge-uuid-noop"
			Expect(k8sClient.Status().Update(ctx, k8sPolicy)).To(Succeed())

			// Reconcile — no update should be needed since fields already match.
			_, err = controller.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      noopResourceName,
					Namespace: noopNS,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify the K8s status is Ready.
			k8sPolicy = &dependencytrackv1alpha1.Policy{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: noopResourceName, Namespace: noopNS}, k8sPolicy)
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sPolicy.Status.UUID).To(Equal("converge-uuid-noop"))

			cond := meta.FindStatusCondition(k8sPolicy.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("PolicySynced"))

			// Verify the remote fields remain unchanged.
			mockDT.mu.Lock()
			defer mockDT.mu.Unlock()
			Expect(mockDT.policy.GetName()).To(Equal(testConvergedPolicyName))
			Expect(mockDT.policy.GetViolationState()).To(Equal(string(dependencytrackv1alpha1.ViolationStateFail)))
		})
	})

	Context("findPolicyByName", func() {
		It("should return true when name matches", func() {
			list := []dtapi.Policy{
				{Name: "Existing", Uuid: "uuid-1"},
				{Name: "Different", Uuid: "uuid-2"},
			}
			found, uuid := findPolicyByName(list, "Existing")
			Expect(found).To(BeTrue())
			Expect(uuid).To(Equal("uuid-1"))
		})

		It("should return false when name does not match", func() {
			list := []dtapi.Policy{
				{Name: "Existing", Uuid: "uuid-1"},
				{Name: "Different", Uuid: "uuid-2"},
			}
			found, uuid := findPolicyByName(list, "Missing")
			Expect(found).To(BeFalse())
			Expect(uuid).To(BeEmpty())
		})

		It("should handle empty list", func() {
			found, uuid := findPolicyByName([]dtapi.Policy{}, "Any")
			Expect(found).To(BeFalse())
			Expect(uuid).To(BeEmpty())
		})
	})

	Context("negative tests — remote API failure status reasons", func() {
		const negNS = testNS

		BeforeEach(func() {
			mockDT.failCreatePolicy = false
			mockDT.failCreateCondition = false
			mockDT.failListPolicies = false
			mockDT.failGetPolicy = false
			mockDT.failGetPolicyUUID = ""
			mockDT.failDeletePolicy = false
		})

		It("should report PolicyCreateFailed when CREATE returns HTTP 500", func() {
			const name = "test-neg-create-fail"
			mockDT.failCreatePolicy = true

			res := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  negNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Create Fail Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: testSeverityHigh},
					},
				},
			}
			Expect(k8sClient.Create(ctx, res)).To(Succeed())
			DeferCleanup(func() {
				rx := &dependencytrackv1alpha1.Policy{}
				if e := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, rx); e == nil {
					_ = k8sClient.Delete(ctx, rx)
				}
			})

			_, err := controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).To(HaveOccurred())

			p := &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PolicyCreateFailed"))
			Expect(cond.Message).To(ContainSubstring("Failed to create policy in DependencyTrack"))
		})

		It("should report PolicyListFailed when LIST returns HTTP 500", func() {
			const name = "test-neg-list-fail"
			mockDT.failListPolicies = true

			res := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  negNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "List Fail Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: testSeverityHigh},
					},
				},
			}
			Expect(k8sClient.Create(ctx, res)).To(Succeed())
			DeferCleanup(func() {
				rx := &dependencytrackv1alpha1.Policy{}
				if e := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, rx); e == nil {
					_ = k8sClient.Delete(ctx, rx)
				}
			})

			_, err := controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).To(HaveOccurred())

			p := &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PolicyListFailed"))
			Expect(cond.Message).To(ContainSubstring("Failed to list policies from DependencyTrack"))
		})

		It("should report PolicyGetFailed when GET owned UUID returns HTTP 500", func() {
			const name = "test-neg-get-fail"
			mockDT.failGetPolicyUUID = "pre-existing-uuid-abc"

			res := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  negNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Get Fail Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: testSeverityHigh},
					},
				},
			}
			Expect(k8sClient.Create(ctx, res)).To(Succeed())
			DeferCleanup(func() {
				rx := &dependencytrackv1alpha1.Policy{}
				if e := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, rx); e == nil {
					_ = k8sClient.Delete(ctx, rx)
				}
			})

			p := &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			p.Status.UUID = "pre-existing-uuid-abc"
			Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())

			_, err := controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).To(HaveOccurred())

			p = &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PolicyGetFailed"))
			Expect(cond.Message).To(ContainSubstring("Failed to get policy from DependencyTrack"))
		})

		It("should report ConditionCreateFailed when CREATE CONDITION returns HTTP 500", func() {
			const name = "test-neg-cond-create-fail"
			mockDT.failCreateCondition = true

			// Pre-create the policy so controller goes the UUID-get path and then condition-create path.
			mockDT.policy = &dtapi.Policy{
				Name:           "Cond Fail Policy",
				Uuid:           "cond-fail-uuid-123",
				Global:         boolPtr(true),
				Operator:       string(dependencytrackv1alpha1.PolicyOperatorAny),
				ViolationState: string(dependencytrackv1alpha1.ViolationStateWarn),
			}

			res := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  negNS,
					Finalizers: []string{policyFinalizer},
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Cond Fail Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: testSeverityHigh},
					},
				},
			}
			Expect(k8sClient.Create(ctx, res)).To(Succeed())
			DeferCleanup(func() {
				rx := &dependencytrackv1alpha1.Policy{}
				if e := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, rx); e == nil {
					_ = k8sClient.Delete(ctx, rx)
				}
			})

			p := &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			p.Status.UUID = "cond-fail-uuid-123"
			Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())

			_, err := controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).To(HaveOccurred())

			p = &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("ConditionCreateFailed"))
			Expect(cond.Message).To(ContainSubstring("Failed to create policy condition in DependencyTrack"))
		})

		// ---- T02: Deletion finalizer retention tests ----

		It("should remove the finalizer after successful deletion", func() {
			const name = "test-del-success"

			mockDT.policy = &dtapi.Policy{
				Name:           "Delete Success Policy",
				Uuid:           "del-success-uuid-123",
				Global:         boolPtr(true),
				Operator:       string(dependencytrackv1alpha1.PolicyOperatorAny),
				ViolationState: string(dependencytrackv1alpha1.ViolationStateWarn),
			}

			res := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: negNS,
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Delete Success Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: testSeverityHigh},
					},
				},
			}
			Expect(k8sClient.Create(ctx, res)).To(Succeed())

			// First reconcile adds the finalizer (envtest strips them from Create responses).
			_, err := controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).NotTo(HaveOccurred())

			p := &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			p.Status.UUID = "del-success-uuid-123"
			Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())

			// Second reconcile: delete the resource to trigger the deletion path.
			Expect(k8sClient.Delete(ctx, p)).To(Succeed())

			_, err = controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).NotTo(HaveOccurred())

			// Verify the controller recorded the deletion event.
			Eventually(fakeRecorder.Events).Should(Receive(ContainSubstring("PolicyDeleted")))
		})

		It("should remove the finalizer when DT returns 404 during delete (already absent)", func() {
			const name = "test-del-404"

			// Do NOT set mockDT.policy — the mock server returns 404 for any UUID GET/DELETE.
			mockDT.policy = nil

			res := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: negNS,
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Delete 404 Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: testSeverityHigh},
					},
				},
			}
			Expect(k8sClient.Create(ctx, res)).To(Succeed())

			// First reconcile adds the finalizer.
			_, err := controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).NotTo(HaveOccurred())

			p := &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			// Set a UUID so the controller tries to delete, but the mock has no matching policy → 404.
			p.Status.UUID = "nonexistent-uuid-999"
			Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())

			// Second reconcile: delete the resource to trigger the deletion path.
			Expect(k8sClient.Delete(ctx, p)).To(Succeed())

			// Reconcile should NOT return an error — 404 means already absent.
			_, err = controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).NotTo(HaveOccurred())

			// Verify the finalizer was removed — envtest may GC the object
			// once the finalizer is gone, so handle both outcomes.
			p = &dependencytrackv1alpha1.Policy{}
			getErr := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)
			if getErr == nil {
				Expect(controllerutil.ContainsFinalizer(p, policyFinalizer)).To(BeFalse())
			} else {
				Expect(errors.IsNotFound(getErr)).To(BeTrue())
			}
		})

		It("should retain the finalizer and set PolicyDeleteFailed on delete error", func() {
			const name = "test-del-fail-retains"

			mockDT.failDeletePolicy = true
			mockDT.policy = &dtapi.Policy{
				Name:           "Delete Fail Policy",
				Uuid:           "del-fail-uuid-123",
				Global:         boolPtr(true),
				Operator:       string(dependencytrackv1alpha1.PolicyOperatorAny),
				ViolationState: string(dependencytrackv1alpha1.ViolationStateWarn),
			}

			res := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: negNS,
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Delete Fail Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: testSeverityHigh},
					},
				},
			}
			Expect(k8sClient.Create(ctx, res)).To(Succeed())

			// First reconcile adds the finalizer.
			_, err := controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).NotTo(HaveOccurred())

			p := &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			p.Status.UUID = "del-fail-uuid-123"
			Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())

			// Second reconcile: delete the resource to trigger the deletion path.
			Expect(k8sClient.Delete(ctx, p)).To(Succeed())

			// Reconcile should return an error (retryable).
			_, err = controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).To(HaveOccurred())

			// Verify finalizer is STILL present (retryable).
			p = &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(p, policyFinalizer)).To(BeTrue())

			// Verify status condition reports the failure.
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PolicyDeleteFailed"))
			Expect(cond.Message).To(ContainSubstring("Failed to delete policy from DependencyTrack"))
		})

		It("should set CredentialsError when auth fails during deletion", func() {
			const name = "test-del-auth-fail"

			mockDT.policy = &dtapi.Policy{
				Name:           "Auth Fail Policy",
				Uuid:           "del-auth-uuid-123",
				Global:         boolPtr(true),
				Operator:       string(dependencytrackv1alpha1.PolicyOperatorAny),
				ViolationState: string(dependencytrackv1alpha1.ViolationStateWarn),
			}

			res := &dependencytrackv1alpha1.Policy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: negNS,
				},
				Spec: dependencytrackv1alpha1.PolicySpec{
					Operator:       dependencytrackv1alpha1.PolicyOperatorAny,
					Name:           "Auth Fail Policy",
					ViolationState: dependencytrackv1alpha1.ViolationStateWarn,
					Conditions: []dependencytrackv1alpha1.PolicyCondition{
						{Subject: dependencytrackv1alpha1.PolicyConditionSubjectSeverity, Operator: dependencytrackv1alpha1.PolicyConditionOperatorIs, Value: testSeverityHigh},
					},
				},
			}
			Expect(k8sClient.Create(ctx, res)).To(Succeed())

			// First reconcile adds the finalizer.
			_, err := controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).NotTo(HaveOccurred())

			p := &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			p.Status.UUID = "del-auth-uuid-123"
			Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())

			// Second reconcile: delete the resource to trigger the deletion path.
			Expect(k8sClient.Delete(ctx, p)).To(Succeed())

			// Use a bad DTProvider that will fail authentication.
			badProvider := &mockClientProvider{
				url:    "http://127.0.0.1:1",
				getErr: fmt.Errorf("auth failed: unreachable"), // force Get() to return error
			}
			controller.DTProvider = badProvider

			_, err = controller.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: negNS}})
			Expect(err).To(HaveOccurred())

			// Verify finalizer is retained.
			p = &dependencytrackv1alpha1.Policy{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: negNS}, p)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(p, policyFinalizer)).To(BeTrue())

			// Verify status condition reports auth failure.
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("CredentialsError"))

			// Restore the DT provider.
			controller.DTProvider = &mockClientProvider{url: mockServer.URL}
		})
	})
})
