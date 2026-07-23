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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// ======================================================================
// Unit tests — indexFuncNotificationRuleBySecretRef
// ======================================================================

var _ = Describe("indexFuncNotificationRuleBySecretRef", func() {
	var indexFn func(client.Object) []string

	BeforeEach(func() {
		indexFn = indexFuncNotificationRuleBySecretRef
	})

	It("returns nil for non-NotificationRule objects", func() {
		result := indexFn(&corev1.Secret{})
		Expect(result).To(BeNil())
	})

	It("returns nil when SecretRef is nil", func() {
		rule := &dependencytrackv1alpha1.NotificationRule{
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Name:        "test",
				Scope:       dependencytrackv1alpha1.NotificationRuleScopePortfolio,
				TriggerType: dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
				Level:       dependencytrackv1alpha1.NotificationRuleLevelWarn,
				PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{
					Name: "pub",
				},
			},
		}
		result := indexFn(rule)
		Expect(result).To(BeNil())
	})

	It("returns the Secret name when SecretRef is set", func() {
		rule := &dependencytrackv1alpha1.NotificationRule{
			Spec: dependencytrackv1alpha1.NotificationRuleSpec{
				Name:        "test",
				Scope:       dependencytrackv1alpha1.NotificationRuleScopePortfolio,
				TriggerType: dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
				Level:       dependencytrackv1alpha1.NotificationRuleLevelWarn,
				PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{
					Name: "pub",
				},
				PublisherConfigSecretRef: &dependencytrackv1alpha1.PublisherConfigSecretRef{
					Name: "my-secret",
					Key:  "config",
				},
			},
		}
		result := indexFn(rule)
		Expect(result).To(Equal([]string{"my-secret"}))
	})
})

// ======================================================================
// Unit tests — enqueueRequestsForSecret mapper
// ======================================================================

var _ = Describe("enqueueRequestsForSecret mapper", func() {
	var fakeClient client.Client

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		Expect(dependencytrackv1alpha1.AddToScheme(scheme)).To(Succeed())

		fakeClient = fakeclient.NewClientBuilder().
			WithScheme(scheme).
			WithIndex(
				&dependencytrackv1alpha1.NotificationRule{},
				indexFieldNotificationRuleSecretRef,
				indexFuncNotificationRuleBySecretRef,
			).
			Build()
	})

	Context("When a NotificationRule references the Secret name", func() {
		BeforeEach(func() {
			rule := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-rule",
					Namespace: "default",
				},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name:        "my-rule",
					Scope:       dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType: dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:       dependencytrackv1alpha1.NotificationRuleLevelWarn,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{
						Name: "my-pub",
					},
					PublisherConfigSecretRef: &dependencytrackv1alpha1.PublisherConfigSecretRef{
						Name: "my-secret",
						Key:  "config",
					},
				},
			}
			Expect(fakeClient.Create(context.Background(), rule)).To(Succeed())
		})

		It("returns the rule's enqueue request", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-secret",
					Namespace: "default",
				},
			}
			reqs := enqueueRequestsForSecret(context.Background(), fakeClient, secret)
			Expect(reqs).To(HaveLen(1))
			Expect(reqs[0].NamespacedName).To(Equal(types.NamespacedName{
				Name:      "my-rule",
				Namespace: "default",
			}))
		})
	})

	Context("When a NotificationRule references a different Secret", func() {
		BeforeEach(func() {
			rule := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "other-rule",
					Namespace: "default",
				},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name:        "other-rule",
					Scope:       dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType: dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:       dependencytrackv1alpha1.NotificationRuleLevelWarn,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{
						Name: "other-pub",
					},
					PublisherConfigSecretRef: &dependencytrackv1alpha1.PublisherConfigSecretRef{
						Name: "other-secret",
						Key:  "config",
					},
				},
			}
			Expect(fakeClient.Create(context.Background(), rule)).To(Succeed())
		})

		It("returns no requests for a non-matching Secret", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "different-secret",
					Namespace: "default",
				},
			}
			reqs := enqueueRequestsForSecret(context.Background(), fakeClient, secret)
			Expect(reqs).To(BeEmpty())
		})
	})

	Context("When a NotificationRule has no SecretRef", func() {
		BeforeEach(func() {
			rule := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "no-ref-rule",
					Namespace: "default",
				},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name:        "no-ref-rule",
					Scope:       dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType: dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:       dependencytrackv1alpha1.NotificationRuleLevelWarn,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{
						Name: "no-ref-pub",
					},
				},
			}
			Expect(fakeClient.Create(context.Background(), rule)).To(Succeed())
		})

		It("returns no requests regardless of Secret name", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "any-secret",
					Namespace: "default",
				},
			}
			reqs := enqueueRequestsForSecret(context.Background(), fakeClient, secret)
			Expect(reqs).To(BeEmpty())
		})
	})

	Context("When multiple rules reference the same Secret", func() {
		BeforeEach(func() {
			// Rule 1 references "shared-secret"
			rule1 := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rule-one",
					Namespace: "default",
				},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name:        "rule-one",
					Scope:       dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType: dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:       dependencytrackv1alpha1.NotificationRuleLevelWarn,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{
						Name: "pub",
					},
					PublisherConfigSecretRef: &dependencytrackv1alpha1.PublisherConfigSecretRef{
						Name: "shared-secret",
						Key:  "config",
					},
				},
			}
			// Rule 2 references "shared-secret" in the same namespace
			rule2 := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rule-two",
					Namespace: "default",
				},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name:        "rule-two",
					Scope:       dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType: dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:       dependencytrackv1alpha1.NotificationRuleLevelWarn,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{
						Name: "pub",
					},
					PublisherConfigSecretRef: &dependencytrackv1alpha1.PublisherConfigSecretRef{
						Name: "shared-secret",
						Key:  "config",
					},
				},
			}
			Expect(fakeClient.Create(context.Background(), rule1)).To(Succeed())
			Expect(fakeClient.Create(context.Background(), rule2)).To(Succeed())
		})

		It("returns only the rules that reference the Secret — bounded fan-out", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "shared-secret",
					Namespace: "default",
				},
			}
			reqs := enqueueRequestsForSecret(context.Background(), fakeClient, secret)
			Expect(reqs).To(HaveLen(2))
			// Both rules should be enqueued — fan-out is bounded to referencing rules only
			names := make(map[string]struct{})
			for _, req := range reqs {
				names[req.Name] = struct{}{}
			}
			Expect(names).To(HaveKey("rule-one"))
			Expect(names).To(HaveKey("rule-two"))
		})
	})

	Context("When a rule and Secret are in different namespaces but share a name", func() {
		BeforeEach(func() {
			// Rule in "other-ns" referencing "shared-secret"
			rule := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "other-rule",
					Namespace: "other-ns",
				},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name:        "other-rule",
					Scope:       dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType: dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:       dependencytrackv1alpha1.NotificationRuleLevelWarn,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{
						Name: "other-pub",
					},
					PublisherConfigSecretRef: &dependencytrackv1alpha1.PublisherConfigSecretRef{
						Name: "shared-secret",
						Key:  "config",
					},
				},
			}
			Expect(fakeClient.Create(context.Background(), rule)).To(Succeed())
		})

		It("does not enqueue the rule when a Secret with the same name changes in a different namespace", func() {
			// Secret in "default" namespace with same name as the one referenced by the rule in "other-ns"
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "shared-secret",
					Namespace: "default",
				},
			}
			reqs := enqueueRequestsForSecret(context.Background(), fakeClient, secret)
			Expect(reqs).To(BeEmpty())
		})
	})
})
