package internal

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func StartAssumeRoleMessage(profile AwsProfile) {
	fmt.Fprintf(os.Stderr, "Assuming Profile: %s\n", profile.Name)
	fmt.Fprintf(os.Stderr, "Role: %s\n", getRole(profile))
	fmt.Fprintf(os.Stderr, "SourceProfile: %s\n", getSourceProfile(profile))
}

func SuccessAssumeRoleMessage(creds AwsCredentials) {
	fmt.Fprintf(os.Stderr, "Success!\n")
	fmt.Fprintf(os.Stderr, "Expires: %s\n\n", formatExpiresAt(creds.ExpiresAt))
}

func OutputAwsCredentials(profileName string, creds AwsCredentials, format string) {
	switch format {
	case "powershell":
		printPowerShellCredentials(profileName, creds)
	case "bash":
		printCredentials(profileName, creds)
	case "fish":
		printFishCredentials(profileName, creds)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func getRole(profile AwsProfile) string {
	if profile.RoleArn != "" {
		return profile.RoleArn
	}
	return profile.SsoRoleName
}

func getSourceProfile(profile AwsProfile) string {
	if profile.SourceProfile != "" {
		return profile.SourceProfile
	}
	return profile.SsoSession
}

// printCredentials prints the credentials in a way that can easily be sourced
// with bash.
func printCredentials(profileName string, creds AwsCredentials) {
	fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyId)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("export AWS_SECURITY_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("export AWS_ASSUMED_ROLE=\"%s\"\n", profileName)
	fmt.Printf("export AWS_SESSION_EXPIRATION=\"%s\"\n", formatExpiresAt(creds.ExpiresAt))
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# eval $(%s)\n", strings.Join(os.Args, " "))
}

// printFishCredentials prints the credentials in a way that can easily be sourced
// with fish.
func printFishCredentials(profileName string, creds AwsCredentials) {
	fmt.Printf("set -gx AWS_ACCESS_KEY_ID \"%s\";\n", creds.AccessKeyId)
	fmt.Printf("set -gx AWS_SECRET_ACCESS_KEY \"%s\";\n", creds.SecretAccessKey)
	fmt.Printf("set -gx AWS_SESSION_TOKEN \"%s\";\n", creds.SessionToken)
	fmt.Printf("set -gx AWS_SECURITY_TOKEN \"%s\";\n", creds.SessionToken)
	fmt.Printf("set -gx AWS_ASSUMED_ROLE \"%s\";\n", profileName)
	fmt.Printf("set -gx AWS_SESSION_EXPIRATION \"%s\";\n", formatExpiresAt(creds.ExpiresAt))
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# eval (%s)\n", strings.Join(os.Args, " "))
}

// printPowerShellCredentials prints the credentials in a way that can easily be sourced
// with Windows powershell using Invoke-Expression.
func printPowerShellCredentials(profileName string, creds AwsCredentials) {
	fmt.Printf("$env:AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyId)
	fmt.Printf("$env:AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
	fmt.Printf("$env:AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("$env:AWS_SECURITY_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("$env:AWS_ASSUMED_ROLE=\"%s\"\n", profileName)
	fmt.Printf("$env:AWS_SESSION_EXPIRATION=\"%s\"\n", formatExpiresAt(creds.ExpiresAt))
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# %s | Invoke-Expression \n", strings.Join(os.Args, " "))
}

func formatExpiresAt(expires time.Time) string {
	loc, _ := time.LoadLocation("Europe/London")
	return expires.In(loc).Format(time.RFC3339)
}
