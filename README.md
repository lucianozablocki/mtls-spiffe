## mTLS SPIFFE connection

- First, we need to generate a signing certificate for our trust domain. This can be done with the following command:

`sudo openssl req -new -newkey rsa:2048 -x509 -nodes -days 365 -subj "/CN=localhost" -addext "subjectAltName=URI:spiffe://localhost/" -addext "keyUsage=critical,keyCertSign" -keyout ca.key -out ca.crt`

- Next up, client and server certificates (signed by this authority) must be generated. For this, we create a new key:

`openssl genrsa -out server.key 2048`

- And create a certificate signing request, using this key, with no subject

`sudo openssl req -new -key server.key -out server.csr -subj '/'`

- Lastly, we sign the certificate, using the key of the certificate authority that we created. Configuration file specifies the extensions used, that comply the SPIFFE standards.

`sudo openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 365 -out server.crt -extensions v3_ca -extfile ./server.cnf`

- Analogous steps are followed to obtain a client certificate signed by the same authority.

- The certificates must be located in */cert* folder 

- Run server with:

`sudo go run -v server/server.go`

- To make a request to the server, run the client:

`sudo go run -v client/client.go`

- On client side, you should see the classic *'Hello, world!'* message

