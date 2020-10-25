# tcsaver

## Description

tcsaver extracts certificates and private keys from a Traefik-generated acme.json file so that they can be reused in other services. The acme.json file is monitored for changes, so it is suitable to leave running continuously.

A configuration file is used to list the domains for which the certificates should be extracted. When searching the acme.json file for the domain names, both the main domain and the SANs will be searched. The extracted certificates and key files are saved with the name of the domain listed in the configuration file (not the main domain from the acme.json file). For example, if the domain name is `www.example.com`, the certificate will be saved as `www.example.com.pem` in the `/certs` directory. The extension can be modified by using the `TCSAVER_CERT_EXTENSION` and `TCSAVER_KEY_EXTENSION` environment variables.

The configuration file can optionally be watched for changes as well, so domains can be added and removed from the list without restarting the utility.

## Environment variables

| Variable | Default | Description |
|---|---|---|
| TCSAVER_WATCH | unset | If set, tcsaver will watch the configuration file for changes. |
| TCSAVER_CERT_EXTENSION | `.pem` | Extension to add to certificate files. |
| TCSAVER_KEY_EXTENSION | `.pem` | Extension to add to private key files. |

## Volumes

| Volume | Description |
|---|---|
| `/config.yaml` | Configuration file |
| `/acme.json` | The acme.json file produced by Traefik. |
| `/certs` | Certificates will be saved in this directory. |
| `/private` | Private keys will be saved in this directory. |

## Configuration file

The configuration file is a YAML file with an array of domains to search for.

```yaml
domains:
  - www.example.com
  - subdomain.example.com
```

## Sample docker-compose file

```yaml
version: '3.7'

services: 
  tcsaver:
    image: jcrummy/tcsaver:latest
    environment:
      - TCSAVER_WATCH=true
      - TCSAVER_CERT_EXTENSION=.crt
      - TCSAVER_KEY_EXTENSION=.key
    volumes:
      - /path/to/config.yaml:/config.yaml
      - /path/to/acme.json:/acme.json
      - /tls/certs:/certs
      - /tls/private:/private
```

## Issues

To file an issue for this utility, please do so at the github repository: github.com/jcrummy/tcsaver.
