package internal

type AwsProfile struct {
	Name          string
	RoleArn       string
	SourceProfile string
	MfaSerial     string
}
