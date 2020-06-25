// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package msalbase

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type AuthorityType int

const (
	AuthorityTypeNone AuthorityType = iota
	AuthorityTypeAad
	AuthorityTypeAdfs
	AuthorityTypeB2C
)

type AuthorityInfo struct {
	Host                  string
	CanonicalAuthorityURI string
	AuthorityType         AuthorityType
	userRealmURIPrefix    string
	ValidateAuthority     bool
	Tenant                string
}

func canonicalizeAuthorityURI(input string) string {
	val := input
	// todo: ensure ends with /
	return strings.ToLower(val)
}

func validateAuthorityURI(input string) error {
	return nil
}

func getFirstPathSegment(u *url.URL) (string, error) {
	pathParts := strings.Split(u.EscapedPath(), "/")
	if len(pathParts) >= 2 {
		return pathParts[1], nil
	}

	return "", errors.New("Authority does not have two segments")
}

func createAuthorityInfo(authorityType AuthorityType, authorityURI string, validateAuthority bool) (*AuthorityInfo, error) {

	u, err := url.Parse(authorityURI)
	if err != nil {
		return nil, err
	}

	host := u.Hostname()
	userRealmURIPrefix := fmt.Sprintf("https://%v/common/userrealm/", host)

	tenant, err := getFirstPathSegment(u)
	if err != nil {
		return nil, err
	}

	canonicalAuthorityURI := fmt.Sprintf("https://%v/%v/", host, tenant)

	return &AuthorityInfo{host, canonicalAuthorityURI, authorityType, userRealmURIPrefix, validateAuthority, tenant}, nil
}

func CreateAuthorityInfoFromAuthorityUri(authorityURI string, validateAuthority bool) (*AuthorityInfo, error) {
	canonicalURI := canonicalizeAuthorityURI(authorityURI)
	err := validateAuthorityURI(canonicalURI)
	if err != nil {
		return nil, err
	}

	// todo: check for other authority types...
	authorityType := AuthorityTypeAad

	return createAuthorityInfo(authorityType, canonicalURI, validateAuthority)
}
