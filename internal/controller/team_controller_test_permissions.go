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
	"testing"
)

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
			input:    []string{"PORTFOLIO_VIEW"},
			sep:      ",",
			expected: "PORTFOLIO_VIEW",
		},
		{
			name:     "multiple elements",
			input:    []string{"VIEW_PORTFOLIO", "PORTFOLIO_VIEW"},
			sep:      ",",
			expected: "VIEW_PORTFOLIO,PORTFOLIO_VIEW",
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
