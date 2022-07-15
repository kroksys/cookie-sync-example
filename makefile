website:
	@clear
	@go run ./example/website/main.go

dsp:
	@clear
	@go run ./example/dsp/main.go

ssp1:
	@clear
	@go run ./example/ssp1/main.go

ssp2:
	@clear
	@go run ./example/ssp2/main.go

ssp3:
	@clear
	@go run ./example/ssp3/main.go

up:
	docker-compose -f example/docker-compose.yaml up --remove-orphans --build

down:
	docker-compose -f example/docker-compose.yaml down --remove-orphans

CERT_FOLDER=./example/server/cert/
V3EXT_FILE=./example/server/v3.ext
# Generate cert authority
ca:
	openssl req -x509 -nodes -new -sha512 -days 365 -newkey rsa:4096 \
		-keyout ${CERT_FOLDER}ca.key -out ${CERT_FOLDER}ca.pem -subj "/C=US/CN=MY-CA"
	openssl x509 -outform pem -in ${CERT_FOLDER}ca.pem -out ${CERT_FOLDER}ca.crt

define generate_cert
	openssl req -new -nodes -newkey rsa:4096 \
  		-keyout ${CERT_FOLDER}$(1).key -out ${CERT_FOLDER}$(1).csr \
  		-subj "/C=US/ST=State/L=City/O=Some-Organization-Name/CN=$(1)"
	openssl x509 -req -sha512 -days 365 \
		-extfile ${V3EXT_PATH} \
		-CA ${CERT_FOLDER}ca.crt -CAkey ${CERT_FOLDER}ca.key -CAcreateserial \
		-in ${CERT_FOLDER}$(1).csr \
		-out ${CERT_FOLDER}$(1).crt
endef
certs:
	$(call generate_cert,website.com)
	$(call generate_cert,dsp.com)
	$(call generate_cert,ssp1.com)
	$(call generate_cert,ssp2.com)
	$(call generate_cert,ssp3.com)

.PHONY: website dsp up down ca certs ssp1 ssp2 ssp3