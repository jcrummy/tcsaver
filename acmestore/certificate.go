package acmestore

import "encoding/base64"

// Certificate is the certificate information
type Certificate struct {
	Domains Domain `json:"domain"`
	Cert    string `json:"certificate"`
	Key     string `json:"key"`
}

// Domain information for a certifiate
type Domain struct {
	Main string   `json:"main"`
	SANs []string `json:"sans"`
}

// CertPEM retrieves the PEM encoded certificate
func (c Certificate) CertPEM() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(c.Cert))
}

// KeyPEM retrieves the PEM encoded key
func (c Certificate) KeyPEM() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(c.Key))
}

// MainDomain returns the main domain name for the certificate
func (c Certificate) MainDomain() string {
	return c.Domains.Main
}

// IncludesDomain returns true if the certificate is for the provided domain
func (c Certificate) IncludesDomain(domain string) bool {
	if c.Domains.Main == domain {
		return true
	}
	for i := 0; i < len(c.Domains.SANs); i++ {
		if c.Domains.SANs[i] == domain {
			return true
		}
	}
	return false
}
