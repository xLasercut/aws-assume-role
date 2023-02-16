#!/bin/bash

set -e

_uninstall() {
  if [[ -z $SHELL_FILE ]]; then
    echo "SHELL_FILE not set"
    return 1
  fi

  echo "deleting assume role script"
  rm -rf ~/__assume-role.sh

  echo "removing source command from $SHELL_FILE"
  sed -i '/source ~\/__assume-role.sh/d' "$SHELL_FILE"
}

_uninstall

