// Copyright (c) 2021 MacEwan University. All rights reserved.
//
// This source code is licensed under the MIT-style license found in
// the LICENSE file in the root directory of this source tree.

package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// NRPS implements Names & Roles Provisioning Services functions.
type NRPS struct {
	Endpoint *url.URL
	Limit    int
	NextPage *url.URL
	Target   *Connector
}

// A Membership represents a course membership with a brief class description.
type Membership struct {
	ID      string
	Context LTIContext
	Members []Member
}

// A LTIContext represents a brief course description used in Names & Roles.
type LTIContext struct {
	ID    string
	Label string
	Title string
}

// A Member represents a participant in a LTI-enabled process.
type Member struct {
	Status             string
	Name               string
	Picture            string
	GivenName          string `json:"given_name"`
	FamilyName         string `json:"family_name"`
	MiddleName         string `json:"middle_name"`
	Email              string
	UserID             string `json:"user_id"`
	LisPersonSourceDid string `json:"lis_person_sourcedid"`
	Roles              []string
}

// UpgradeNRPS provides a Connector upgraded for NRPS calls.
func (c *Connector) UpgradeNRPS() (*NRPS, error) {
	// Check for endpoint.
	nrpsRawClaim, ok := c.LaunchToken.Get("https://purl.imsglobal.org/spec/lti-nrps/claim/namesroleservice")
	if !ok {
		return nil, ErrUnsupportedService
	}
	nrpsClaim, ok := nrpsRawClaim.(map[string]interface{})
	if !ok {
		return nil, errors.New("names and roles information improperly formatted")
	}
	nrpsString, ok := nrpsClaim["context_memberships_url"]
	if !ok {
		return nil, errors.New("names and roles endpoint not found")
	}
	nrps, err := url.Parse(nrpsString.(string))
	if err != nil {
		return nil, fmt.Errorf("names and roles endpoint parse error: %w", err)
	}

	return &NRPS{
		Endpoint: nrps,
		Target:   c,
	}, nil
}

// GetMembership gets the launched course (referred to as a Context in LTI) membership from the platform. Using
// GetPagedMemberships as a helper, it checks for next page links, fetching and appending them to the output.
func (n *NRPS) GetMembership() (Membership, error) {
	var (
		limit          int
		hasMore        bool
		membership     Membership
		moreMembership Membership
		err            error
	)

	membership, hasMore, err = n.GetPagedMembership(limit)
	if err != nil {
		return Membership{}, fmt.Errorf("get paged membership error: %w", err)
	}

	for hasMore {
		moreMembership, hasMore, err = n.GetPagedMembership(limit)
		if err != nil {
			return Membership{}, fmt.Errorf("get more membership error: %w", err)
		}
		membership.Members = append(membership.Members, moreMembership.Members...)
	}

	return membership, nil
}

// GetPagedMembership gets paged Memberships for the launched course.
func (n *NRPS) GetPagedMembership(limit int) (Membership, bool, error) {
	if limit < 0 {
		return Membership{}, false, errors.New("invalid paging limit")
	}
	scopes := []string{"https://purl.imsglobal.org/spec/lti-nrps/scope/contextmembership.readonly"}

	query, err := url.ParseQuery(n.Endpoint.RawQuery)
	if err != nil {
		return Membership{}, false, fmt.Errorf("could not parse NRPS query values: %w", err)
	}
	if limit != 0 {
		query.Add("limit", strconv.Itoa(limit))
	}

	// Set the initial limit query parameter.
	pagedURI, err := url.Parse(n.Endpoint.String())
	if err != nil {
		return Membership{}, false, fmt.Errorf("could not parse NRPS endpoint: %w", err)
	}
	pagedURI.RawQuery = query.Encode()
	s := ServiceRequest{
		Scopes: scopes,
		Method: http.MethodGet,
		URI:    pagedURI,
		Accept: "application/vnd.ims.lti-nrps.v2.membershipcontainer+json",
	}

	// If there was a next page set from a previous response, use it.
	if n.NextPage != nil {
		s.URI = n.NextPage
	}
	headers, body, err := n.Target.makeServiceRequest(s)
	if err != nil {
		return Membership{}, false, fmt.Errorf("get paged membership make service request error: %w", err)
	}

	defer body.Close()
	var membership Membership
	err = json.NewDecoder(body).Decode(&membership)
	if err != nil {
		return Membership{}, false, fmt.Errorf("could not decode get paged membership response body: %w", err)
	}

	// Get the next page link from the response headers.
	nextPageLink := headers.Get("link")
	if nextPageLink == "" || !strings.Contains(nextPageLink, `rel="next"`) {
		// If there are no further next page links, set the NRPS NextPage field to nil.
		n.NextPage = nil
		return membership, false, nil
	}

	nextPageString := strings.Trim(nextPageLink, "<>")
	nextPage, err := url.Parse(nextPageString)
	if err != nil {
		return Membership{}, false, fmt.Errorf("could not parse next page URI from response headers: %w", err)
	}
	n.NextPage = nextPage

	return membership, true, nil
}

// GetLaunchingMember returns a Member struct representing the user that performed the launch. Status is not included
// in the launch message.
func (n *NRPS) GetLaunchingMember() (Member, error) {
	var launchingMember Member

	rawLaunchEmail, ok := n.Target.LaunchToken.Get("email")
	if !ok {
		return Member{}, errors.New("launching member email not found")
	}
	launchEmail, ok := rawLaunchEmail.(string)
	if !ok {
		return Member{}, errors.New("could not assert launching member email")
	}
	launchingMember.Email = launchEmail

	rawFamilyName, ok := n.Target.LaunchToken.Get("family_name")
	if !ok {
		return Member{}, errors.New("launching member family name not found")
	}
	familyName, ok := rawFamilyName.(string)
	if !ok {
		return Member{}, errors.New("could not assert launching member family name")
	}
	launchingMember.FamilyName = familyName

	rawGivenName, ok := n.Target.LaunchToken.Get("given_name")
	if !ok {
		return Member{}, errors.New("launching member family name not found")
	}
	givenName, ok := rawGivenName.(string)
	if !ok {
		return Member{}, errors.New("could not assert launching member family name")
	}
	launchingMember.GivenName = givenName

	rawName, ok := n.Target.LaunchToken.Get("name")
	if !ok {
		return Member{}, errors.New("launching member name not found")
	}
	name, ok := rawName.(string)
	if !ok {
		return Member{}, errors.New("could not assert launching member name")
	}
	launchingMember.Name = name

	rawRoles, ok := n.Target.LaunchToken.Get("https://purl.imsglobal.org/spec/lti/claim/roles")
	if !ok {
		return Member{}, errors.New("launching member roles not found")
	}
	rolesInterfaces, ok := rawRoles.([]interface{})
	if !ok {
		return Member{}, errors.New("could not assert launching member roles")
	}
	launchingRoles := convertInterfaceToStringSlice(rolesInterfaces)
	launchingMember.Roles = launchingRoles

	launchingMember.UserID = n.Target.LaunchToken.Subject()

	return launchingMember, nil
}
