// Copyright (c) 2021 MacEwan University. All rights reserved.
//
// This source code is licensed under the MIT-style license found in
// the LICENSE file in the root directory of this source tree.

// Package nonpersistent implements an in-memory (non-persistent) data store. It implements all of the Storer
// interfaces, so it can be used for any and all LTI data.
package nonpersistent

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/macewan-cs/lti-example/pkg/datastore"
)

// Store implements an in-memory datastore.
type Store struct {
	Registrations *sync.Map
	Deployments   *sync.Map
	Nonces        *sync.Map
	LaunchData    *sync.Map
	AccessTokens  *sync.Map
}

// DefaultStore provides a single default datastore as a package variable so that other LTI functions can
// fall back on this datastore whenever the user does not explicitly specify a datastore.
var DefaultStore *Store = New()

// New returns an empty, zeroed sync.Map for each Storer interface.
func New() *Store {
	return &Store{
		Registrations: &sync.Map{},
		Deployments:   &sync.Map{},
		Nonces:        &sync.Map{},
		LaunchData:    &sync.Map{},
		AccessTokens:  &sync.Map{},
	}
}

func registrationIndex(issuer, clientID string) string {
	return issuer + "/" + clientID
}

// StoreRegistration stores a Registration in-memory.
func (s *Store) StoreRegistration(reg datastore.Registration) error {
	// Store the registration both with and without the client ID: later, the registration can be retrieved with or
	// without it. See FindRegistrationByIssuerAndClientID for further details.
	s.Registrations.Store(reg.Issuer, reg)
	s.Registrations.Store(registrationIndex(reg.Issuer, reg.ClientID), reg)
	return nil
}

func deploymentIndex(issuer, deploymentID string) string {
	return issuer + "/" + deploymentID
}

// StoreDeployment stores a deployment ID in-memory.
func (s *Store) StoreDeployment(issuer string, d datastore.Deployment) error {
	if issuer == "" {
		return errors.New("received empty issuer argument")
	}
	if err := datastore.ValidateDeploymentID(d.DeploymentID); err != nil {
		return fmt.Errorf("received invalid deployment ID: %w", err)
	}

	s.Deployments.Store(deploymentIndex(issuer, d.DeploymentID), d)
	return nil
}

// FindRegistrationByIssuerAndClientID looks up and returns either a Registration by the issuer or the datastore error
// ErrRegistrationNotFound.
func (s *Store) FindRegistrationByIssuerAndClientID(issuer, clientID string) (datastore.Registration, error) {
	if issuer == "" {
		return datastore.Registration{}, errors.New("received empty issuer argument")
	}

	index := issuer
	if clientID != "" {
		// Use the client ID to disambiguate multiple registrations for an issuer.  The (optional) client ID
		// parameter can disambiguate between multiple registrations from a single issuer.
		//
		// Source: http://www.imsglobal.org/spec/lti/v1p3/#client_id-login-parameter
		index = registrationIndex(issuer, clientID)
	}

	registration, ok := s.Registrations.Load(index)
	if !ok {
		return datastore.Registration{}, datastore.ErrRegistrationNotFound
	}

	return registration.(datastore.Registration), nil
}

// FindDeployment looks up and returns either a Deployment by the issuer and deployment ID or the datastore error
// ErrDeploymentNotFound.
func (s *Store) FindDeployment(issuer, deploymentID string) (datastore.Deployment, error) {
	if issuer == "" {
		return datastore.Deployment{}, errors.New("received empty issuer argument")
	}
	if err := datastore.ValidateDeploymentID(deploymentID); err != nil {
		return datastore.Deployment{}, fmt.Errorf("received invalid deployment ID: %w", err)
	}

	deployment, ok := s.Deployments.Load(deploymentIndex(issuer, deploymentID))
	if !ok {
		return datastore.Deployment{}, datastore.ErrDeploymentNotFound
	}
	return deployment.(datastore.Deployment), nil
}

// StoreNonce stores a Nonce in-memory. Since the nonce and target_link_uri values have similarly scoped verifications
// required, use the the unique nonce value as a key to store the target_link_uri value. This is used to verify the OIDC
// login request target_link_uri is the same as the claim of the same name in the launch id_token.
func (s *Store) StoreNonce(nonce, targetLinkURI string) error {
	if nonce == "" {
		return errors.New("received empty nonce argument")
	}
	if targetLinkURI == "" {
		return errors.New("received empty issuer argument")
	}

	s.Nonces.Store(nonce, targetLinkURI)
	return nil
}

// TestAndClearNonce looks up a nonce, clears the entry if found, and returns whether it was found via the error
// return. If the nonce wasn't found, it returns the datastore error ErrNonceNotFound. If it was found, it returns nil.
func (s *Store) TestAndClearNonce(nonce, targetLinkURI string) error {
	if nonce == "" {
		return errors.New("received empty nonce argument")
	}
	if targetLinkURI == "" {
		return errors.New("received empty target link uri argument")
	}

	checkURI, ok := s.Nonces.Load(nonce)
	if !ok {
		return datastore.ErrNonceNotFound
	}

	s.Nonces.Delete(nonce)

	if checkURI != targetLinkURI {
		return datastore.ErrNonceTargetLinkURIMismatch
	}

	return nil
}

// StoreLaunchData stores the launch data, i.e. the id_token JWT.
func (s *Store) StoreLaunchData(launchID string, launchData json.RawMessage) error {
	if launchID == "" {
		return errors.New("received empty launchID argument")
	}
	if len(launchData) == 0 {
		return errors.New("received empty launchData argument")
	}

	s.LaunchData.Store(launchID, launchData)
	return nil
}

// FindLaunchData retrieves a cached launchData.
func (s *Store) FindLaunchData(launchID string) (json.RawMessage, error) {
	if launchID == "" {
		return nil, errors.New("received empty launchID argument")
	}

	launchData, ok := s.LaunchData.Load(launchID)
	if !ok {
		return nil, datastore.ErrLaunchDataNotFound
	}
	return launchData.(json.RawMessage), nil
}

func accessTokenIndex(tokenURI, clientID string, scopes []string) string {
	return tokenURI + clientID + strings.Join(scopes[:], " ")
}

// StoreAccessToken stores bearer tokens for potential reuse.
func (s *Store) StoreAccessToken(token datastore.AccessToken) error {
	if token.TokenURI == "" {
		return errors.New("received empty tokenURI")
	}
	if token.ClientID == "" {
		return errors.New("received empty clientID")
	}
	if len(token.Scopes) == 0 {
		return errors.New("received empty scopes")
	}
	if token.Token == "" {
		return errors.New("received empty accessToken")
	}
	zeroTime := time.Time{}
	if token.ExpiryTime == zeroTime {
		return errors.New("received empty expiry time")
	}

	sort.Strings(token.Scopes)

	storeValue, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("error encoding access token to store: %w", err)
	}

	s.AccessTokens.Store(accessTokenIndex(token.TokenURI, token.ClientID, token.Scopes), storeValue)
	return nil
}

// FindAccessToken retrieves bearer tokens for potential reuse.
func (s *Store) FindAccessToken(tokenURI, clientID string, scopes []string) (datastore.AccessToken, error) {
	if tokenURI == "" {
		return datastore.AccessToken{}, errors.New("received empty tokenURI")
	}
	if clientID == "" {
		return datastore.AccessToken{}, errors.New("received empty clientID")
	}
	if len(scopes) == 0 {
		return datastore.AccessToken{}, errors.New("received empty scopes")
	}

	index := accessTokenIndex(tokenURI, clientID, scopes)
	storeValue, ok := s.AccessTokens.Load(index)
	if !ok {
		return datastore.AccessToken{}, datastore.ErrAccessTokenNotFound
	}
	storeBytes, ok := storeValue.([]byte)
	if !ok {
		return datastore.AccessToken{}, errors.New("could not assert access token")
	}

	var accessToken datastore.AccessToken
	err := json.Unmarshal(storeBytes, &accessToken)
	if err != nil {
		return datastore.AccessToken{}, fmt.Errorf("could not decode access token: %w", err)
	}
	if accessToken.ExpiryTime.Before(time.Now()) {
		return datastore.AccessToken{}, datastore.ErrAccessTokenExpired
	}

	return accessToken, nil
}
