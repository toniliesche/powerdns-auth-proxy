#!/bin/bash

COMMAND="powerdns-auth-proxy --config /etc/powerdns-auth-proxy/config.yaml migrate-db --init-file /var/lib/powerdns-auth-proxy/res/init.yaml"

if [ ! -z ${IMPORTFILE} ]; then
  COMMAND="${COMMAND} --import-file ${IMPORT_FILE}"
fi

su apiuser -s /bin/bash -c "${COMMAND}"
