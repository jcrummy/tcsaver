package acmestore

import (
	"encoding/json"
	"io"
)

// ACME is the json file structure used by Traefik
type ACME map[string]Provider

// Provider represents the data saved in the acme.json file
type Provider struct {
	// Account - we ignore this section
	Certificates []Certificate `json:"Certificates"`
}

// NewStore returns an ACME store based on a provided io.Reader
func NewStore(data io.Reader) (*ACME, error) {
	ret := ACME{}
	decoder := json.NewDecoder(data)
	err := decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// Find looks for certificate for a specific domain
func (a *ACME) Find(domain string) (*Certificate, error) {
	for _, provider := range a {
		for _, cert := range provider.Certificates {
			if cert.IncludesDomain(domain) {
				return &cert, nil
			}
		}
	}
	return nil, ErrDomainNotFound
}
