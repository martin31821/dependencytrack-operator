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

// driftLogSuccessfulPublish returns true when the spec and remote
// logSuccessfulPublish values disagree.
func driftLogSuccessfulPublish(remote *bool, spec *bool) bool {
	if remote == nil && spec == nil {
		return false
	}
	if remote == nil || spec == nil {
		return true
	}
	return *remote != *spec
}

// driftNotifyChildren returns true when the spec and remote notifyChildren
// values disagree.
func driftNotifyChildren(remote *bool, spec *bool) bool {
	if remote == nil && spec == nil {
		return false
	}
	if remote == nil || spec == nil {
		return true
	}
	return *remote != *spec
}

// driftScheduleCron returns true when the spec and remote scheduleCron values
// disagree.  The remote is a *string (from dtapi.NotificationRule); the spec
// is a plain string (from NotificationRuleSpec), so we dereference the remote
// pointer (defaulting to empty string) before comparing.
func driftScheduleCron(remote *string, spec string) bool {
	remoteVal := ""
	if remote != nil {
		remoteVal = *remote
	}
	if remoteVal == "" && spec == "" {
		return false
	}
	if remoteVal == "" || spec == "" {
		return true
	}
	return remoteVal != spec
}

// driftScheduleSkipUnchanged returns true when the spec and remote
// scheduleSkipUnchanged values disagree.
func driftScheduleSkipUnchanged(remote *bool, spec *bool) bool {
	if remote == nil && spec == nil {
		return false
	}
	if remote == nil || spec == nil {
		return true
	}
	return *remote != *spec
}
