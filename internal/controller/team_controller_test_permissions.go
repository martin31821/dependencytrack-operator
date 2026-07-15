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
	"reflect"
	"testing"

	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

const (
	testPerm = "PORTFOLIO_MANAGEMENT"
	testBOM  = "BOM_UPLOAD"
	testView = "VIEW_PORTFOLIO"
)

func TestPermissionDelta(t *testing.T) {
	tests := []struct {
		name          string
		current       []dtapi.Permission
		desired       []string
		wantAdd       []string
		wantRemove    []string
		wantCanonical []string
	}{
		{
			name:          "adds and removes permissions",
			current:       []dtapi.Permission{{Name: testBOM}, {Name: testView}},
			desired:       []string{testView, testPerm},
			wantAdd:       []string{testPerm},
			wantRemove:    []string{testBOM},
			wantCanonical: []string{testPerm, testView},
		},
		{
			name:          "deduplicates and sorts desired permissions",
			desired:       []string{testView, testBOM, testView},
			wantAdd:       []string{testBOM, testView},
			wantCanonical: []string{testBOM, testView},
		},
		{
			name:       "clears all permissions",
			current:    []dtapi.Permission{{Name: testView}},
			desired:    []string{},
			wantRemove: []string{testView},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdd, gotRemove, gotCanonical := permissionDelta(tt.current, tt.desired)
			if !reflect.DeepEqual(gotAdd, tt.wantAdd) {
				t.Errorf("permissionDelta() additions = %v; want %v", gotAdd, tt.wantAdd)
			}
			if !reflect.DeepEqual(gotRemove, tt.wantRemove) {
				t.Errorf("permissionDelta() removals = %v; want %v", gotRemove, tt.wantRemove)
			}
			if !reflect.DeepEqual(gotCanonical, tt.wantCanonical) {
				t.Errorf("permissionDelta() canonical = %v; want %v", gotCanonical, tt.wantCanonical)
			}
		})
	}
}

func TestJoinString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		sep      string
		expected string
	}{
		{
			name:     "nil slice",
			input:    nil,
			sep:      ",",
			expected: "",
		},
		{
			name:     "empty slice",
			input:    []string{},
			sep:      ",",
			expected: "",
		},
		{
			name:     "single element",
			input:    []string{testPerm},
			sep:      ",",
			expected: testPerm,
		},
		{
			name:     "multiple elements",
			input:    []string{testView, testPerm},
			sep:      ",",
			expected: testView + "," + testPerm,
		},
		{
			name:     "different separator",
			input:    []string{"A", "B"},
			sep:      "|",
			expected: "A|B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinString(tt.input, tt.sep)
			if result != tt.expected {
				t.Errorf("joinString(%v, %q) = %q; want %q", tt.input, tt.sep, result, tt.expected)
			}
		})
	}
}
