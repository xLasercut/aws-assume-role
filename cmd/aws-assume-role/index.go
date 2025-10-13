package main

import (
	"flag"

	. "github.com/xLasercut/aws-assume-role/internal"
)

var AppVersion string = "dev"

func init() {
	flag.Usage = Usage
}

func main() {
	profileName, duration, awsConfigFiles, format := ParseArgs(AppVersion)
	profileChain, baseProfile := GetProfileChain(awsConfigFiles, profileName)

	var awsCreds AwsCredentials

	if baseProfile.SsoStartUrl != "" {
		awsCreds = AssumeBaseRoleSso(baseProfile)
	} else {
		awsCreds = AssumeBaseRoleSts(baseProfile, duration)
	}

	for _, profile := range profileChain {
		awsCreds = AssumeRoleSts(profile, awsCreds, duration)
	}

	OutputAwsCredentials(profileName, awsCreds, format)
	return
}
