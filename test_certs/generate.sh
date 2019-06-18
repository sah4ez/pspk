CERT_NAME=localhost

openssl genrsa -des3 -out rootCA.key 4096
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 1024 -out rootCA.crt
openssl genrsa -out $CERT_NAME.key 2048
openssl req -new -sha256 -key $CERT_NAME.key -subj "/C=US/ST=CA/O=MyOrg, Inc./CN=$CERT_NAME" -out $CERT_NAME.csr
openssl x509 -req -in $CERT_NAME.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out $CERT_NAME.crt -days 500 -sha256
openssl x509 -in $CERT_NAME.crt -text -noout

