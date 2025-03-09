#!/bin/sh

if [ ! -f /.configured ] && [ -d /configure.d ]; then
  for file in /configure.d/*.sh; do
    if [ -f ${file} ]; then
      ${file}

      if [ $? -ne 0 ]; then
        echo "failed configuring container ..."
        echo "failed executing script ${file} ..."

        exit 1
      fi
    fi
  done

  touch /.configured
fi

if [ -d /startup.d ]; then
  for file in /startup.d/*.sh; do
    if [ -f ${file} ]; then
      ${file}

      if [ $? -ne 0 ]; then
        echo "failed starting up container ..."
        echo "failed executing script ${file} ..."

        exit 1
      fi
    fi
  done
fi

exec "$@"
