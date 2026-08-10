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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// cloneTeamSpec performs a deep-enough copy of TeamSpec to assert that
// ValidateTeam leaves the caller's values (especially the Groups slice and
// the OIDC pointer) bit-for-bit unaltered.
func cloneTeamSpec(in TeamSpec) TeamSpec {
	out := in
	if in.Permissions != nil {
		out.Permissions = append([]string{}, in.Permissions...)
	}
	if in.OIDC != nil {
		cloned := &TeamOIDCConfig{}
		if in.OIDC.Groups != nil {
			cloned.Groups = append([]string{}, in.OIDC.Groups...)
		}
		out.OIDC = cloned
	}
	return out
}

func TestValidateTeam(t *testing.T) {
	type tc struct {
		name    string
		spec    TeamSpec
		wantErr bool
		substr  string // optional fragment asserted to appear in the error text
		substrs []string
	}

	cases := []tc{
		{
			name:    "nil_oidc_disables_validation_no_op",
			spec:    TeamSpec{OIDC: nil},
			wantErr: false,
		},
		{
			name:    "explicit_nil_groups_with_enabled_config_allowed",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: nil}},
			wantErr: false,
		},
		{
			name:    "empty_groups_defers_mapping_clearance_to_controller",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{}}},
			wantErr: false,
		},
		{
			name:    "single_wellformed_group_accepted",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"Admins"}}},
			wantErr: false,
		},
		{
			name:    "mixed_case_variants_are_distinct_under_exact_trimmed_equality",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"Admins", "admins", "ADMINS"}}},
			wantErr: false,
		},
		{
			name:    "unicode_cased_pair_preserved_distinct",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"Département-Vente", "département-vente"}}},
			wantErr: false,
		},
		{
			name:    "surrounding_whitespace_retained_verbatim_when_unique",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{" Ops-Engineers "}}},
			wantErr: false,
		},
		{
			name:    "whitespace_only_entry_rejected",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{" "}}},
			wantErr: true,
			substr:  "is empty or consists only of whitespace",
		},
		{
			name:    "tabbed_whitespace_only_entry_rejected",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"\t\v\f\r\n\u00A0"}}},
			wantErr: true,
			substr:  "is empty or consists only of whitespace",
		},
		{
			name:    "literal_empty_string_entry_rejected",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{""}}},
			wantErr: true,
			substr:  "is empty or consists only of whitespace",
		},
		{
			name:    "valid_followed_by_blank_flags_the_blank",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"Ops", " "}}},
			wantErr: true,
			substr:  "is empty or consists only of whitespace",
		},
		{
			name:    "identical_spelling_duplicate_detected",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"Admin", "Admin"}}},
			wantErr: true,
			substr:  "duplicates",
		},
		{
			name:    "leading_trailing_space_variant_collapses_to_trimmed_twin",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"Admin", " Admin ", "\tAdmin\n"}}},
			wantErr: true,
			substr:  "duplicates",
		},
		{
			name:    "reverse_order_still_detects_duplicate_of_trimmed_twin",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{" Admin ", "Admin"}}},
			wantErr: true,
			substr:  "duplicates",
		},
		{
			name:    "two_blanks_each_independently_reported",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"", " "}}},
			wantErr: true,
			substrs: []string{"is empty or consists only of whitespace", "[1]"},
		},
		{
			name:    "aggregated_errors_surface_every_offense_in_one_denial",
			spec:    TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{"", "", "Dup", " Dup "}}},
			wantErr: true,
			substrs: []string{"is empty or consists only of whitespace", "duplicates"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			before := cloneTeamSpec(c.spec)

			err := ValidateTeam(c.spec)

			if c.wantErr {
				require.Error(t, err, "expected validation failure for %q", c.name)
				if c.substr != "" {
					assert.Contains(t, err.Error(), c.substr)
				}
				for _, s := range c.substrs {
					assert.Contains(t, err.Error(), s, "missing fragment %q in error %q", s, err.Error())
				}
			} else {
				assert.NoError(t, err, "did not expect validation failure for %q", c.name)
			}

			// ValidateTeam is a pure predicate: it must not mutate the
			// caller-owned spec, groups, or OIDC pointer.
			assert.Equal(t, before, c.spec, "ValidateTeam mutated the input spec")
			if c.spec.OIDC != nil && c.spec.OIDC.Groups != nil {
				assert.Equal(t, before.OIDC.Groups, c.spec.OIDC.Groups,
					"ValidateTeam altered the OIDC.Groups slice contents/order")
			}
		})
	}
}

// TestValidateTeamDoesNotNormalize documents the promise relied upon by the
// controller (T04): the admission predicate never normalizes, sorts, or
// deduplicates in place. Operators retain full control over upstream casing
// and whitespace so IdP claim matching is unaffected.
func TestValidateTeamDoesNotNormalize(t *testing.T) {
	adversarial := []string{
		"Admins",
		" admins ", // leading/trailing space + lower
		"",         // empty (would be rejected if fed to the webhook)
	}
	spec := TeamSpec{OIDC: &TeamOIDCConfig{Groups: adversarial}}

	before := append([]string(nil), adversarial...)
	_ = ValidateTeam(spec) // ignore result; this case is expected to fail

	assert.Equal(t, before, spec.OIDC.Groups,
		"groups must round-trip unmodified (no trim/sort/dedupe)")
}
