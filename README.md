# AWS assume role

This is a wrapper for aws cli to make assuming AWS roles much easier

## Credentials setup

Update aws credentials file `~/.aws/credentials` and/or `~/.aws/config`
```text
[user]
region = eu-west-2
aws_access_key_id = ABCSDFAAWERA
aws_secret_access_key = ASDFADFASDFASFWERQWER

[profile dev]
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

[profile dev]
role_arn = arn:aws:iam::1234567890:role/developer
region = eu-west-2
source_profile = user
mfa_serial = arn:aws:iam::1234567890:mfa/user

[profile dev-testing]
role_arn = arn:aws:iam::1234567890:role/developer-testing
region = eu-west-2
source_profile = dev
```

For SSO credentials and chaining
```text
[sso-session my-sso]
sso_region = eu-west-2
sso_start_url = https://abc.awsapps.com/start
sso_registration_scopes = sso:account:access

[profile admin]
sso_session = my-sso
sso_account_id = 1234567890
sso_role_name = Admin
region = eu-west-2

[profile dev]
role_arn = arn:aws:iam::1234567890:role/developer
source_profile = admin
region = eu-west-2

```

You can also config single roles likes this:
```text
[admin]
sso_region = eu-west-2
sso_start_url = https://abc.awsapps.com/start
sso_registration_scopes = sso:account:access
sso_account_id = 1234567890
sso_role_name = Admin
region = eu-west-2

[dev]
aws_access_key_id = ABCSDFAAWERA
aws_secret_access_key = ASDFADFASDFASFWERQWER
role_arn = arn:aws:iam::1234567890:role/developer
region = eu-west-2
mfa_serial = arn:aws:iam::1234567890:mfa/user
```

## Usage

simply run the following command to assume the corresponding roles. If mfa_serial is configured, the script should ask for your MFA token.
```shell
assume-role dev
```

If the role requires MFA, you will be asked for the token

```shell
$ assume-role dev
Assuming Profile: dev
RoleArn: arn:aws:iam::1234567890:role/developer
SourceProfile: user
Assume Role MFA token code: 123456
Success!
Expires: 2023-01-01 10:31:06 +0000 UTC
```
If no command is provided, `assume-role` will output the temporary security credentials:

```shell
$ assume-role dev
export AWS_ACCESS_KEY_ID="AAAA....AAAA"
export AWS_SECRET_ACCESS_KEY="BVS...a1Sfd"
export AWS_SESSION_TOKEN="AQ...1SDF=="
export AWS_SECURITY_TOKEN="AQ...1SDF=="
export AWS_ASSUMED_ROLE="dev"
export AWS_SESSION_EXPIRATION="2023-01-01T00:00:00Z"
# Run this to configure your shell:
# eval $(assume-role dev)
```
If you use `eval $(assume-role)` frequently, you may want to create a alias for it:

zsh

```shell
alias assume-role='function(){eval $(command assume-role $@);}'
```

bash
```shell
function assume-role { eval $( $(which assume-role) $@); }
```
