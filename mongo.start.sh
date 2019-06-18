#/!/bin/sh 

CERT_NAME=localhost

docker run -it --rm -v $(pwd)/test_certs/:/etc/certs -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=root --name mongo -p 0.0.0.0:27017:27017 mongo --auth --sslMode requireSSL --sslPEMKeyFile /etc/certs/${CERT_NAME}.pem --sslCAFile /etc/certs/rootCA.crt --sslAllowInvalidHostnames

