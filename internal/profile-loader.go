package internal

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
)

func GetProfileChain(awsConfigFiles *ini.File, profileName string) ([]AwsProfile, AwsProfile) {
	var fullProfileChain []AwsProfile
	profileNext := profileName

	allProfileNames := awsConfigFiles.SectionStrings()

	for {
		if !profileExists(allProfileNames, profileNext) {
			err := errors.New(fmt.Sprintf("profile \"%s\" does not exist", profileNext))
			CheckError(err, "Could not load profile information")
		}

		awsProfile := getAwsProfile(awsConfigFiles, profileNext)
		fullProfileChain = append([]AwsProfile{awsProfile}, fullProfileChain...)

		if awsProfile.SourceProfile == "" && awsProfile.SsoSession == "" {
			break
		}

		if len(fullProfileChain) > 20 {
			err := errors.New("profile chain length exceeded 20 items")
			CheckError(err, "Could not load profile information")
		}

		profileNext = getNextProfileName(awsProfile)
	}

	if len(fullProfileChain) < 2 {
		err := errors.New("invalid assume profile config")
		CheckError(err, "Could not load profile information")
	}

	profileChain := fullProfileChain[2:]
	baseProfileChain := fullProfileChain[:2]
	return profileChain, getBaseProfile(baseProfileChain)
}

func getAwsProfile(awsConfigFiles *ini.File, profileName string) AwsProfile {
	credentialsFileSection := awsConfigFiles.Section(profileName)
	return AwsProfile{
		Name:                  profileName,
		RoleArn:               credentialsFileSection.Key("role_arn").String(),
		SourceProfile:         credentialsFileSection.Key("source_profile").String(),
		MfaSerial:             credentialsFileSection.Key("mfa_serial").String(),
		Region:                credentialsFileSection.Key("region").String(),
		SsoRegion:             credentialsFileSection.Key("sso_region").String(),
		SsoRoleName:           credentialsFileSection.Key("sso_role_name").String(),
		SsoAccountId:          credentialsFileSection.Key("sso_account_id").String(),
		SsoSession:            credentialsFileSection.Key("sso_session").String(),
		SsoStartUrl:           credentialsFileSection.Key("sso_start_url").String(),
		SsoRegistrationScopes: credentialsFileSection.Key("sso_registration_scopes").String(),
		AwsAccessKeyId:        credentialsFileSection.Key("aws_access_key_id").String(),
		AwsSecretAccessKey:    credentialsFileSection.Key("aws_secret_access_key").String(),
	}
}

func profileExists(allProfileNames []string, profileName string) bool {
	for _, name := range allProfileNames {
		if name == profileName {
			return true
		}
	}
	return false
}

func getNextProfileName(profile AwsProfile) string {
	if profile.SourceProfile != "" {
		return profile.SourceProfile
	}
	return profile.SsoSession
}

func getBaseProfile(baseProfileChain []AwsProfile) AwsProfile {
	return AwsProfile{
		Name:                  baseProfileChain[1].Name,
		RoleArn:               baseProfileChain[1].RoleArn,
		SourceProfile:         baseProfileChain[1].SourceProfile,
		MfaSerial:             baseProfileChain[1].MfaSerial,
		Region:                baseProfileChain[1].Region,
		SsoRegion:             baseProfileChain[0].SsoRegion,
		SsoRoleName:           baseProfileChain[1].SsoRoleName,
		SsoAccountId:          baseProfileChain[1].SsoAccountId,
		SsoSession:            baseProfileChain[1].SsoSession,
		SsoStartUrl:           baseProfileChain[0].SsoStartUrl,
		SsoRegistrationScopes: baseProfileChain[0].SsoRegistrationScopes,
		AwsAccessKeyId:        baseProfileChain[0].AwsAccessKeyId,
		AwsSecretAccessKey:    baseProfileChain[0].AwsSecretAccessKey,
	}
}
