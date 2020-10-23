package config_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/jcrummy/tcsaver/config"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		Name   string
		YAML   string
		Config *config.Config
		Err    error
	}{
		{"All Defined",
			`acmefile: thisacme.json
certdir: /this/cert/dir
keydir: /this/key/dir
domains:
  - firstdomain.com
  - www.seconddomain.com
  - anotherdomain.ca`,
			&config.Config{
				ACMEFile: "thisacme.json",
				CertDir:  "/this/cert/dir",
				KeyDir:   "/this/key/dir",
				Domains:  []string{"firstdomain.com", "www.seconddomain.com", "anotherdomain.ca"},
			},
			nil,
		},
		{"Only domains",
			`domains:
  - firstdomain.com
  - www.seconddomain.com
  - anotherdomain.ca`,
			&config.Config{
				ACMEFile: "/acme.json",
				CertDir:  "/certs",
				KeyDir:   "/private",
				Domains:  []string{"firstdomain.com", "www.seconddomain.com", "anotherdomain.ca"},
			},
			nil,
		},
	}

	for _, test := range tests {
		got, err := config.Load(strings.NewReader(test.YAML))
		if !errors.Is(err, test.Err) {
			t.Errorf("Unexpected error returned loading %s: %s", test.Name, err)
			continue
		}

		if !reflect.DeepEqual(got, test.Config) {
			t.Errorf("Unexpected config returned for %s. Want %s, got %s", test.Name, test.Config, got)
		}
	}
}
