package tcsaver

import (
	"log"
	"os"
	"os/signal"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/jcrummy/tcsaver/acmestore"
	"github.com/jcrummy/tcsaver/config"
)

// CLI is a runnable command-line interface
type CLI struct {
	config        config.Config
	watchConfig   bool
	certExtension string
	keyExtension  string
}

// Run is an infinite loop that waits for file changes
func (c *CLI) Run() {
	c.checkEnv()
	err := c.checkConfig()
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}
	err := c.checkACME()
	if err != nil {
		log.Fatalf("Unable to check acme file at least once: %v", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	fileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Unable to create filesystem watcher: %v", err)
	}
	defer fileWatcher.Close()

	fileWatcher.Add(c.config.ACMEFile)
	if c.watchConfig {
		fileWatcher.Add("config.yaml")
	}

	for {
		select {
		case <-shutdown:
			return

		case event := <-fileWatcher.Events:
			switch path.Base(event.Name) {
			case "config.yaml":
				err := c.checkConfig()
				if err != nil {
					log.Printf("Error reloading config file: %v", err)
				}

			case path.Base(c.config.ACMEFile):
				err := c.checkACME()
				if err != nil {
					log.Printf("Error checking acme file: %v", err)
				}
			}
		}
	}
}

func (c *CLI) checkEnv() {
	_, c.watchConfig = os.LookupEnv("TCSAVER_WATCH")

	if ext, ok := os.LookupEnv("TCSAVER_CERT_EXTENSION"); ok {
		c.certExtension = ext
	} else {
		c.certExtension = ".pem"
	}

	if ext, ok := os.LookupEnv("TCSAVER_KEY_EXTENSION"); ok {
		c.keyExtension = ext
	} else {
		c.keyExtension = ".pem"
	}
}

func (c *CLI) checkConfig() error {
	fh, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer fh.Close()
	newConfig, err := config.Load(fh)
	if err != nil {
		return err
	}
	c.config = *newConfig
	return nil
}

func (c *CLI) checkACME() error {
	fh, err := os.Open(c.config.ACMEFile)
	if err != nil {
		return err
	}
	defer fh.Close()
	acme, err := acmestore.NewStore(fh)
	if err != nil {
		return err
	}

	for _, domain := range c.config.Domains {
		c.saveDomain(acme, domain)
	}
	return nil
}

func (c *CLI) saveDomain(store *acmestore.Store, domain string) {
	cert, err := store.Find(domain)
	if err != nil {
		log.Printf("Domain %s not found in acme file", domain)
		return
	}
	certPEM, err := cert.CertPEM()
	if err != nil {
		log.Printf("Error extraing certificate for domain %s: %v", domain, err)
		return
	}
	keyPEM, err := cert.KeyPEM()
	if err != nil {
		log.Printf("Error extracting key for domain %s: %v", domain, err)
		return
	}

	certfile, err := os.Create(path.Join(c.config.CertDir, domain+c.certExtension))
	if err != nil {
		log.Printf("Error creating certificate file for domain %s: %v", domain, err)
		return
	}
	defer certfile.Close()
	_, err = certfile.Write(certPEM)
	if err != nil {
		log.Printf("Error saving certificate file for domain %s: %v", domain, err)
		return
	}

	keyfile, err := os.Create(path.Join(c.config.KeyDir, domain+c.keyExtension))
	if err != nil {
		log.Printf("Error creating key file for domain %s: %v", domain, err)
		return
	}
	defer keyfile.Close()
	_, err = keyfile.Write(keyPEM)
	if err != nil {
		log.Printf("Error saving key file for domain %s: %v", domain, err)
		return
	}
}
