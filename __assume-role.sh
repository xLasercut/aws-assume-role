#!/bin/bash

function __assume_role_help() {
  echo
  echo "assume-role <command>"
  echo
  echo "commands:"
  echo "  help                        - this help screen"
  echo "  <profile>                   - assume role involved with profile"
  echo
  return 1
}

function assume-role() {
  local command=$1

  case $command in
  "help") __assume_role_help ;;
  *) __assume_role "${@:1}" ;;
  esac
}

function __assume_role() {
  local profile_name=$1
  if [[ -z ${profile_name} ]]; then
    echo "Profile name is required." >&2
    return 1
  fi

  local duration_seconds
  local region
  local role_arn
  local mfa_serial
  local source_profile

  region=$(__get_aws_info "${profile_name}" region) || return 1
  role_arn=$(__get_aws_info "${profile_name}" role_arn) || return 1
  source_profile=$(__get_aws_info "${profile_name}" source_profile || echo "default") || return 1
  duration_seconds=$(__get_aws_info "${profile_name}" duration_seconds || echo "3600") || return 1
  mfa_serial=$(__get_aws_info "${profile_name}" mfa_serial || echo "") || return 1

  current_user=$(__get_aws_current_user "${source_profile}") || return 1
  local session_name="${profile_name}-${current_user}"

  local token_code
  if [[ $mfa_serial != "" ]]; then
    echo "MFA token code:"
    read -r token_code
  fi

  local session_token_json
  if [[ $mfa_serial != "" ]]; then
    session_token_json=$(aws sts assume-role \
      --role-arn "${role_arn}" \
      --role-session-name "${session_name}" \
      --region "${region}" \
      --duration-seconds "${duration_seconds}" \
      --serial-number "${mfa_serial}" \
      --token-code "${token_code}" \
      --profile "${source_profile}" \
      --query Credentials) || return 1
  else
    session_token_json=$(aws sts assume-role \
      --role-arn "${role_arn}" \
      --role-session-name "${session_name}" \
      --region "${region}" \
      --duration-seconds "${duration_seconds}" \
      --profile "${source_profile}" \
      --query Credentials) || return 1
  fi

  export AWS_ACCESS_KEY_ID=$(echo "$session_token_json" | jq -r .AccessKeyId) || return 1
  export AWS_SECRET_ACCESS_KEY=$(echo "$session_token_json" | jq -r .SecretAccessKey) || return 1
  export AWS_SESSION_TOKEN=$(echo "$session_token_json" | jq -r .SessionToken) || return 1
  export AWS_SESSION_EXPIRY=$(echo "$session_token_json" | jq -r .Expiration) || return 1

  if [[ -n "${AWS_ACCESS_KEY_ID}" && -n "${AWS_SECRET_ACCESS_KEY}" && -n "${AWS_SESSION_TOKEN}" ]]; then
    export AWS_ROLE=$(echo "${role_arn}" | cut -d'/' -f 2) || return 1
    export AWS_ROLE_ARN="${role_arn}" || return 1
    export AWS_ACCOUNT_ID=$(echo "${role_arn}" | cut -d':' -f 5) || return 1

    echo "Successfully assumed the role with ARN ${role_arn}."
    echo "Access keys valid until ${AWS_SESSION_EXPIRY}."

    return 0
  fi

  return 1
}

function __get_aws_current_user() {
  local source_profile
  source_profile=$(__get_first_source_profile_in_chain "$1")

  local caller_identity_json
  caller_identity_json=$(aws sts get-caller-identity --profile "${source_profile}")
  if ! [[ "${?}" -eq 0 ]]; then
    echo "No AWS credentials found" >&2
    return 1
  fi

  local current_user
  local current_principal_arn
  current_principal_arn=$(echo "${caller_identity_json}" | jq -r .Arn) || return 1
  echo "${current_principal_arn}" | cut -d'/' -f 2
}

function __get_first_source_profile_in_chain() {
  local source_profile=$1
  local chained_source_profile=$1

  while true; do
    chained_source_profile=$(__get_aws_info "${source_profile}" source_profile || echo "")
    if [[ $chained_source_profile == "" ]]; then
      break
    else
      source_profile=$chained_source_profile
    fi
  done

  echo "$source_profile"
}

function __get_aws_info() {
  local profile_name=$1
  local config_key=$2

  aws configure get "$config_key" --profile "$profile_name"
}
