package internal

import "time"

type AwsProfile struct {
	Name                  string
	RoleArn               string
	SourceProfile         string
	MfaSerial             string
	Region                string
	SsoRegion             string
	SsoRoleName           string
	SsoAccountId          string
	SsoSession            string
	SsoStartUrl           string
	SsoRegistrationScopes string
	AwsAccessKeyId        string
	AwsSecretAccessKey    string
}

type SsoClientState struct {
	AccessToken             string
	ClientId                string
	ClientSecret            string
	DeviceCode              string
	VerificationUriComplete string
	StartUrl                string
}

type AwsCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
	ExpiresAt       time.Time
	SessionToken    string
}
