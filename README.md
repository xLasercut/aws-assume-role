# AWS assume role

This is a wrapper for aws cli to make assuming AWS roles much easier

## Requirements

- awscli
- jq

## Installation

## Credentials setup

Update aws credentials file `~/.aws/credentials`
```text
[user]
region = eu-west-2
aws_access_key_id = ABCSDFAAWERA
aws_secret_access_key = ASDFADFASDFASFWERQWER

[dev]
role_arn = arn:aws:iam::1234567890:role/developer
region = eu-west-2
source_profile = user
mfa_serial = arn:aws:iam::1234567890:mfa/user
```

You can also chain roles like this
```text
[user]
region = eu-west-2
aws_access_key_id = ABCSDFAAWERA
aws_secret_access_key = ASDFADFASDFASFWERQWER

[dev]
role_arn = arn:aws:iam::1234567890:role/developer
region = eu-west-2
source_profile = user
mfa_serial = arn:aws:iam::1234567890:mfa/user

[dev-testing]
role_arn = arn:aws:iam::1234567890:role/developer-testing
region = eu-west-2
source_profile = dev
```

## Usage

simply run the following command to assume the corresponding roles. If mfa_serial is configured, the script should ask for your MFA token.
```shell
assume-role dev
```