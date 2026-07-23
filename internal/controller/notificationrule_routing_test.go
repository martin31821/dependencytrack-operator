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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("driftLogSuccessfulPublish", func() {
	It("returns false when both are nil", func() {
		Expect(driftLogSuccessfulPublish(nil, nil)).To(BeFalse())
	})
	It("returns false when both point to true", func() {
		t := true
		Expect(driftLogSuccessfulPublish(&t, &t)).To(BeFalse())
	})
	It("returns false when both point to false", func() {
		f := false
		Expect(driftLogSuccessfulPublish(&f, &f)).To(BeFalse())
	})
	It("returns true when remote is nil and spec is true", func() {
		t := true
		Expect(driftLogSuccessfulPublish(nil, &t)).To(BeTrue())
	})
	It("returns true when remote is true and spec is nil", func() {
		t := true
		Expect(driftLogSuccessfulPublish(&t, nil)).To(BeTrue())
	})
	It("returns true when remote is false and spec is true", func() {
		f, t := false, true
		Expect(driftLogSuccessfulPublish(&f, &t)).To(BeTrue())
	})
})

var _ = Describe("driftNotifyChildren", func() {
	It("returns false when both are nil", func() {
		Expect(driftNotifyChildren(nil, nil)).To(BeFalse())
	})
	It("returns false when both point to true", func() {
		t := true
		Expect(driftNotifyChildren(&t, &t)).To(BeFalse())
	})
	It("returns true when remote is nil and spec is true", func() {
		t := true
		Expect(driftNotifyChildren(nil, &t)).To(BeTrue())
	})
	It("returns true when remote is true and spec is false", func() {
		t, f := true, false
		Expect(driftNotifyChildren(&t, &f)).To(BeTrue())
	})
})

var _ = Describe("driftScheduleCron", func() {
	It("returns false when both are empty/nil", func() {
		Expect(driftScheduleCron(nil, "")).To(BeFalse())
	})
	It("returns false when both are the same expression", func() {
		exp := "0 0 * * *"
		Expect(driftScheduleCron(&exp, exp)).To(BeFalse())
	})
	It("returns true when remote is nil and spec has expression", func() {
		exp := "0 0 * * *"
		Expect(driftScheduleCron(nil, exp)).To(BeTrue())
	})
	It("returns true when remote has expression and spec is empty", func() {
		exp := "0 0 * * *"
		Expect(driftScheduleCron(&exp, "")).To(BeTrue())
	})
	It("returns true when expressions differ", func() {
		a := "0 0 * * *"
		b := "0 12 * * *"
		Expect(driftScheduleCron(&a, b)).To(BeTrue())
	})
})

var _ = Describe("driftScheduleSkipUnchanged", func() {
	It("returns false when both are nil", func() {
		Expect(driftScheduleSkipUnchanged(nil, nil)).To(BeFalse())
	})
	It("returns false when both point to true", func() {
		t := true
		Expect(driftScheduleSkipUnchanged(&t, &t)).To(BeFalse())
	})
	It("returns true when remote is nil and spec is false", func() {
		f := false
		Expect(driftScheduleSkipUnchanged(nil, &f)).To(BeTrue())
	})
	It("returns true when remote is true and spec is nil", func() {
		t := true
		Expect(driftScheduleSkipUnchanged(&t, nil)).To(BeTrue())
	})
	It("returns true when values differ", func() {
		t, f := true, false
		Expect(driftScheduleSkipUnchanged(&t, &f)).To(BeTrue())
	})
})
