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

package deptrack

import (
	"context"
	"fmt"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

// tokenTTL is how long we reuse a cached bearer token before re-authenticating.
// DependencyTrack tokens are valid for 24 h; we refresh early to avoid expiry mid-operation.
const tokenTTL = time.Hour

// ClientProvider authenticates with DependencyTrack and caches the bearer token.
// It is safe for concurrent use.
type ClientProvider struct {
	K8sClient  client.Client
	Namespace  string
	SecretName string
	DTURL      string

	mu          sync.Mutex
	apiClient   *dtapi.APIClient
	cachedToken string
	tokenExpiry time.Time
}

// Get returns a context carrying a valid DependencyTrack bearer token and the
// shared API client. The token is refreshed from the credentials Secret when
// it is absent or within 60 s of expiry.
func (p *ClientProvider) Get(ctx context.Context) (context.Context, *dtapi.APIClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.apiClient == nil {
		p.apiClient = NewAPIClient(p.DTURL)
	}

	if p.cachedToken != "" && time.Now().Add(60*time.Second).Before(p.tokenExpiry) {
		return context.WithValue(ctx, dtapi.ContextAccessToken, p.cachedToken), p.apiClient, nil
	}

	secret := &corev1.Secret{}
	if err := p.K8sClient.Get(ctx, types.NamespacedName{Namespace: p.Namespace, Name: p.SecretName}, secret); err != nil {
		return nil, nil, fmt.Errorf("reading credentials secret %s/%s: %w", p.Namespace, p.SecretName, err)
	}

	username := string(secret.Data["username"])
	password := string(secret.Data["password"])

	token, _, err := p.apiClient.UserAPI.ValidateCredentials(ctx).
		Username(username).
		Password(password).
		Execute()
	if err != nil {
		return nil, nil, fmt.Errorf("authenticating with DependencyTrack: %w", err)
	}

	p.cachedToken = token
	p.tokenExpiry = time.Now().Add(tokenTTL)

	return context.WithValue(ctx, dtapi.ContextAccessToken, token), p.apiClient, nil
}

// Invalidate forces the next Get call to re-authenticate. Call this when a
// 401 response is received from DependencyTrack.
func (p *ClientProvider) Invalidate() {
	p.mu.Lock()
	p.cachedToken = ""
	p.mu.Unlock()
}
