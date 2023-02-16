#!/bin/bash

set -e

_install() {
  if [[ -z $SHELL_FILE ]]; then
    echo "SHELL_FILE not set"
    return 1
  fi

  wget -O ~/__assume-role.sh "https://raw.githubusercontent.com/xLasercut/aws-assume-role/master/__assume-role.sh" || return 1

  if grep -Fxq "source ~/__assume-role.sh" "$SHELL_FILE"; then
    echo "source command already added to $SHELL_FILE. Nothing else to"
    return 0
  else
    echo "adding source command to $SHELL_FILE"
    echo "source ~/__assume-role.sh" >> "$SHELL_FILE"
    return 0
  fi
}

_install
