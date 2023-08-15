#!/bin/bash

BASE_DIR="/tmp/k8s-webhook-server/serving-certs"

if [ ! -d ${BASE_DIR} ]; then
  mkdir -p ${BASE_DIR}
  openssl req \
    -x509 \
    -newkey rsa:4096 \
    -keyout ${BASE_DIR}/tls.key \
    -out ${BASE_DIR}/tls.crt \
    -sha256 \
    -days 365 \
    -nodes \
    -subj "/C=DE/ST=Berlin/L=Berlin/O=Verlag der Tagesspiegel GmbH/OU=Org/CN=*.localhost"
fi
