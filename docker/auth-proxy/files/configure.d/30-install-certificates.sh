#!/bin/bash

if [ -z "${CA_CERTIFICATES_PATH}" ]; then
  echo "CA_CERTIFICATES_PATH is not set. Skipping."
  exit 0
fi

if [ -d ${CA_CERTIFICATES_PATH} ]; then
  echo "Installing certificates from ${CA_CERTIFICATES_PATH}..."
  find ${CA_CERTIFICATES_PATH} -iname "*.crt" -exec cp -v {} /usr/local/share/ca-certificates/ {} \;
elif [ -f ${CA_CERTIFICATES_PATH} ]; then
  echo "Installing certificate ${CA_CERTIFICATES_PATH}..."
  cp -v ${CA_CERTIFICATES_PATH} /usr/local/share/ca-certificates/
fi

update-ca-certificates
