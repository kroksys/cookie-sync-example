# Cookie sync services

## Overview
1. User enters a shopping website [https://website.com] <br />
that has iframe for DSP [iframe.src="https://dsp.com"] <br />
2. By request iframe from dps.com a 3rd party cookie (user_id) is set for dsp.com.
3. Iframe from dsp.com responds with all partner "pixels" to be loaded with the create user_id and redirect_back url.
4. Partner receives request. Creates its own user_id and put the data in matching table. If redirect_url is passed redirects back to the url provided with its user_id.


## Self signed certificates
### 1. create certificate authority
```bash
make ca
```
### 2. generate certificates for websites
```bash 
make certs
```

### 3. make created CA  to be trusted by your OS
Google how to do that for your OS.  
This is the certificate file:
```bash
./nginx/cert/ca.crt
```

### 4. add records to your OS host file
```bash
# /etc/hosts # for macos
127.0.0.1 website.com
127.0.0.1 dsp.com
127.0.0.1 ssp1.com
127.0.0.1 ssp2.com
127.0.0.1 ssp3.com
```

## Starting
```bash
make up # docker-compose -f example/docker-compose.yaml up --remove-orphans --build
```

