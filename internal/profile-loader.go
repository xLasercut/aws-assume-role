package internal

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

func removeAppendSections(profileSectionName string) string {
	removedProfileAppend := strings.Replace(profileSectionName, "profile ", "", 1)
	removedSsoSessionAppend := strings.Replace(removedProfileAppend, "sso-session ", "", 1)
	return removedSsoSessionAppend
}

func getAllProfiles(awsConfigFiles *ini.File) []string {
	allProfileSectionNames := awsConfigFiles.SectionStrings()

	var allProfileNames []string

	for _, profileSectionName := range allProfileSectionNames {
		profileName := removeAppendSections(profileSectionName)
		allProfileNames = append(allProfileNames, profileName)
	}

	return allProfileNames
}

func getProfileSectionName(awsConfigFiles *ini.File, profileName string) string {
	allProfileSections := awsConfigFiles.SectionStrings()

	for _, name := range allProfileSections {
		if name == profileName || name == fmt.Sprintf("profile %s", profileName) || name == fmt.Sprintf("sso-session %s", profileName) {
			return name
		}
	}

	err := errors.New(fmt.Sprintf("profile \"%s\" does not exist", profileName))
	CheckError(err, "Could not load profile information")
	return ""
}

func GetProfileChain(awsConfigFiles *ini.File, profileName string) ([]AwsProfile, AwsProfile) {
	var fullProfileChain []AwsProfile
	profileNext := profileName

	allProfileNames := getAllProfiles(awsConfigFiles)

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

	if len(fullProfileChain) < 1 {
		err := errors.New("invalid assume profile config")
		CheckError(err, "Could not load profile information")
	}

	if len(fullProfileChain) == 1 {
		return []AwsProfile{}, fullProfileChain[0]
	}

	if fullProfileChain[0].RoleArn != "" || fullProfileChain[0].SsoRoleName != "" {
		return fullProfileChain[1:], fullProfileChain[0]
	}

	profileChain := fullProfileChain[2:]
	baseProfileChain := fullProfileChain[:2]
	return profileChain, getBaseProfile(baseProfileChain)
}

func getAwsProfile(awsConfigFiles *ini.File, profileName string) AwsProfile {
	sectionName := getProfileSectionName(awsConfigFiles, profileName)
	credentialsFileSection := awsConfigFiles.Section(sectionName)
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
