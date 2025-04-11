#!/bin/bash

if [ ! -z "${JWT_SECRET_KEY}" ] && [ ! -z "${JWT_PUBLIC_KEY}" ]; then
  echo "Using JWT secret and public keys from environment variables..."
  exit 0
fi

if [ ! -z "${JWT_SECRET_KEY_PATH}" ] && [ ! -z "${JWT_PUBLIC_KEY_PATH}" ]; then
  if [ ! -f "${JWT_SECRET_KEY_PATH}" ] || [ ! -f "${JWT_PUBLIC_KEY_PATH}" ]; then
    echo "JWT secret or public key file does not exist..."
    exit 1
  fi

  echo "Using JWT secret and public keys from files..."
  exit 0
fi

if [ -f /var/lib/powerdns-auth-proxy/data/jwt.privkey.pem ] && [ -f /var/lib/powerdns-auth-proxy/data/jwt.pubkey.pem ]; then
  echo "Using existing JWT secret and public keys..."
  exit 0
fi

echo "Generating JWT secret and public keys..."
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out /var/lib/powerdns-auth-proxy/data/jwt.privkey.pem
openssl ec -in /var/lib/powerdns-auth-proxy/data/jwt.privkey.pem -pubout -out /var/lib/powerdns-auth-proxy/data/jwt.pubkey.pem
chown -R apiuser:apiuser /var/lib/powerdns-auth-proxy/data/jwt.*
