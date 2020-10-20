package acmestore

import (
	"encoding/json"
	"io"
)

// Store contains json file structure used by Traefik
type Store struct {
	Items map[string]Provider
}

// Provider represents the data saved in the acme.json file
type Provider struct {
	// Account - we ignore this section
	Certificates []Certificate `json:"Certificates"`
}

// NewStore returns an ACME store based on a provided io.Reader
func NewStore(data io.Reader) (*Store, error) {
	ret := Store{}
	ret.Items = make(map[string]Provider, 0)
	decoder := json.NewDecoder(data)
	err := decoder.Decode(&ret.Items)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// Find looks for certificate for a specific domain
func (s *Store) Find(domain string) (*Certificate, error) {
	for _, provider := range s.Items {
		for _, cert := range provider.Certificates {
			if cert.IncludesDomain(domain) {
				return &cert, nil
			}
		}
	}
	return nil, ErrDomainNotFound
}
