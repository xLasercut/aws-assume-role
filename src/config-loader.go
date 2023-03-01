package main

import (
    "errors"
    "fmt"
    "gopkg.in/ini.v1"
)

func getProfileChain(awsConfigFiles *ini.File, profileName string) []AwsProfile {
    var profileChain []AwsProfile
    profileNext := profileName

    allProfileNames := awsConfigFiles.SectionStrings()

    for {
        if !profileExists(allProfileNames, profileNext) {
            checkError(errors.New(fmt.Sprintf("profile \"%v\" does not exist", profileNext)), "Could not load profile information")
        }

        awsProfile := getAwsProfile(awsConfigFiles, profileNext)
        if awsProfile.sourceProfile == "" {
            break
        }

        profileChain = append([]AwsProfile{awsProfile}, profileChain...)
        profileNext = awsProfile.sourceProfile
    }

    return profileChain
}

func getAwsProfile(awsConfigFiles *ini.File, profileName string) AwsProfile {
    credentialsFileSection := awsConfigFiles.Section(profileName)
    roleArn := credentialsFileSection.Key("role_arn").String()
    sourceProfile := credentialsFileSection.Key("source_profile").String()
    mfaSerial := credentialsFileSection.Key("mfa_serial").String()
    return AwsProfile{profileName, roleArn, sourceProfile, mfaSerial}
}

func profileExists(allProfileNames []string, profileName string) bool {
    for _, name := range allProfileNames {
        if name == profileName {
            return true
        }
    }
    return false
}
