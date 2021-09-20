deps:
	go install github.com/bketelsen/crypt/bin/crypt@latest

set-config:
	crypt set -plaintext /config/user/development.properties.yml user/development.properties.yml