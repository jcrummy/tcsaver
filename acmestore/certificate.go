package acmestore

import "encoding/base64"

// Certificate is the certificate information
type Certificate struct {
	Domains struct {
		Main string   `json:"main"`
		SANs []string `json:"sans"`
	} `json:"domain"`
	Cert []byte `json:"certificate"`
	Key  []byte `json:"key"`
}

// CertPEM retrieves the PEM encoded certificate
func (c Certificate) CertPEM() (string, error) {
	return base64.StdEncoding.DecodeString(string(c.Cert))
}

// KeyPEM retrieves the PEM encoded key
func (c Certificate) KeyPEM() (string, error) {
	return base64.StdEncoding.DecodeString(string(c.Key))
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
