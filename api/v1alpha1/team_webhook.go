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

package v1alpha1

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-dependencytrack-mko-dev-v1alpha1-team,mutating=false,failurePolicy=Fail,sideEffects=None,resources=teams,verbs=create;update,groups=dependencytrack.mko.dev,versions=v1alpha1,name=team.validating.dependencytrack.mko.dev,admissionReviewVersions=v1

// ValidateTeam is a pure, side-effect-free admission predicate for TeamSpec.
//
// It guards Team.spec.oidc.groups so malformed entries can never reach the
// controller or the DependencyTrack OIDC API. The rules mirror the
// documented semantics of TeamOIDCConfig (see team_types.go):
//
//   - A nil *TeamOIDCConfig disables OIDC management entirely and is a
//     no-op (nothing to validate).
//   - A non-nil config carrying a nil or empty Groups slice is permitted:
//     it is an intentional "map this team to zero groups" state whose
//     eventual convergence (clearing remote mappings) is owned by the
//     controller (S02). Validation must not conflate this with malformation.
//   - Each surviving entry is evaluated on its space-trimmed form. Blank
//     (empty or whitespace-only) entries are rejected.
//   - Duplicate detection uses trimmed CASE-SENSORY (exact) equality, so
//     "Admins" and "ADMINS" are distinct groups; only identically-spelled
//     trimmed values collide. Original casing is left untouched.
//   - Leading/trailing-whitespace variants collapse to the same trimmed
//     key, so "Admin" followed by " Admin " is flagged as a duplicate of its
//     trimmed twin.
//
// Multiple offenses are accumulated and returned as a single joined error so
// a deny response communicates the full set of problems in one shot.
func ValidateTeam(spec TeamSpec) error {
	if spec.OIDC == nil {
		// OIDC management disabled: nothing to validate.
		return nil
	}

	seen := make(map[string]string, len(spec.OIDC.Groups))
	var errs []error

	for i, raw := range spec.OIDC.Groups {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			errs = append(errs, fmt.Errorf(
				"spec.oidc.groups[%d]: %q is empty or consists only of whitespace", i, raw))
			continue
		}
		if prev, ok := seen[trimmed]; ok {
			errs = append(errs, fmt.Errorf(
				"spec.oidc.groups[%d]: %q duplicates %q (matched after trimming)", i, raw, prev))
			continue
		}
		seen[trimmed] = raw
	}

	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

// teamValidator adapts the pure ValidateTeam predicate to the
// controller-runtime admission.CustomValidator contract, so a single
// implementation backs both create and update admission reviews.
type teamValidator struct{}

var _ admission.CustomValidator = &teamValidator{}

func (v *teamValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, validateTeamObject(obj)
}

func (v *teamValidator) ValidateUpdate(_ context.Context, _ runtime.Object, newObj runtime.Object) (admission.Warnings, error) {
	// Updates are denied against the incoming (proposed) spec; the prior
	// state is irrelevant because every admission decision is total over the
	// requested spec.
	return nil, validateTeamObject(newObj)
}

func (v *teamValidator) ValidateDelete(context.Context, runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func validateTeamObject(obj runtime.Object) error {
	t, ok := obj.(*Team)
	if !ok {
		return fmt.Errorf("unexpected admission object type %T, expected *%s.Team", obj, GroupVersion.Group)
	}
	return ValidateTeam(t.Spec)
}

// SetupWebhookWithManager registers the Team validating admission webhook
// with mgr. The webhook is fail-closed (failurePolicy=Fail) and declares
// sideEffects=None, so a misbehaving or unreachable webhook server denies
// the request rather than admitting possibly-corrupt group mappings.
//
// Registration is invoked from cmd/main.go (during manager bootstrap),
// which wires the Team validator into the shared webhook server so the
// admission path is served alongside the rest of the operator.
func (r *Team) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(&teamValidator{}).
		Complete()
}
