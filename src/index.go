package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"os"
	"strings"
	"time"
)

type AwsProfile struct {
	name          string
	roleArn       string
	sourceProfile string
	mfaSerial     string
}

func init() {
	flag.Usage = usage
}

func main() {
	profileName, duration, awsConfigFiles, format := parseArgs()

	profileChain := getProfileChain(awsConfigFiles, profileName)

	baseProfileName := profileChain[0].sourceProfile

	ctx := context.TODO()
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(baseProfileName))
	var credentials aws.Credentials

	for _, profile := range profileChain {
		cfg.Credentials = assumeRole(cfg, profile, baseProfileName, duration)
		credentials = getCredentials(cfg, ctx)
	}

	switch format {
	case "powershell":
		printPowerShellCredentials(profileName, credentials)
	case "bash":
		printCredentials(profileName, credentials)
	case "fish":
		printFishCredentials(profileName, credentials)
	default:
		flag.Usage()
		os.Exit(1)
	}
	return
}

func assumeRole(cfg aws.Config, profile AwsProfile, baseProfileName string, duration int) *aws.CredentialsCache {
	stsClient := sts.NewFromConfig(cfg)

	fmt.Fprintf(os.Stderr, "Assuming Profile: %s\n", profile.name)
	fmt.Fprintf(os.Stderr, "RoleArn: %s\n", profile.roleArn)
	fmt.Fprintf(os.Stderr, "SourceProfile: %s\n", profile.sourceProfile)

	provider := stscreds.NewAssumeRoleProvider(stsClient, profile.roleArn, func(p *stscreds.AssumeRoleOptions) {
		p.RoleSessionName = fmt.Sprintf("%v-%v", profile.name, baseProfileName)
		p.TokenProvider = tokenProvider
		p.Duration = time.Duration(duration) * time.Second

		if profile.mfaSerial != "" {
			p.SerialNumber = &profile.mfaSerial
		}
	})

	return aws.NewCredentialsCache(provider)
}

func tokenProvider() (string, error) {
	var v string
	fmt.Fprint(os.Stderr, "Assume Role MFA token code: ")
	_, err := fmt.Scanln(&v)

	return v, err
}

func getCredentials(cfg aws.Config, ctx context.Context) aws.Credentials {
	credentials, err := cfg.Credentials.Retrieve(ctx)
	checkError(err, "Could not assume role")

	fmt.Fprintf(os.Stderr, "Success!\n")
	fmt.Fprintf(os.Stderr, "Expires: %v\n\n", credentials.Expires)

	return credentials
}

// printCredentials prints the credentials in a way that can easily be sourced
// with bash.
func printCredentials(role string, credentials aws.Credentials) {
	fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", credentials.AccessKeyID)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", credentials.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", credentials.SessionToken)
	fmt.Printf("export AWS_SECURITY_TOKEN=\"%s\"\n", credentials.SessionToken)
	fmt.Printf("export AWS_ASSUMED_ROLE=\"%s\"\n", role)
	fmt.Printf("export AWS_SESSION_EXPIRATION=\"%s\"\n", credentials.Expires.Format(time.RFC3339))
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# eval $(%s)\n", strings.Join(os.Args, " "))
}

// printFishCredentials prints the credentials in a way that can easily be sourced
// with fish.
func printFishCredentials(role string, credentials aws.Credentials) {
	fmt.Printf("set -gx AWS_ACCESS_KEY_ID \"%s\";\n", credentials.AccessKeyID)
	fmt.Printf("set -gx AWS_SECRET_ACCESS_KEY \"%s\";\n", credentials.SecretAccessKey)
	fmt.Printf("set -gx AWS_SESSION_TOKEN \"%s\";\n", credentials.SessionToken)
	fmt.Printf("set -gx AWS_SECURITY_TOKEN \"%s\";\n", credentials.SessionToken)
	fmt.Printf("set -gx AWS_ASSUMED_ROLE \"%s\";\n", role)
	fmt.Printf("set -gx AWS_SESSION_EXPIRATION \"%s\";\n", credentials.Expires.Format(time.RFC3339))
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# eval (%s)\n", strings.Join(os.Args, " "))
}

// printPowerShellCredentials prints the credentials in a way that can easily be sourced
// with Windows powershell using Invoke-Expression.
func printPowerShellCredentials(role string, credentials aws.Credentials) {
	fmt.Printf("$env:AWS_ACCESS_KEY_ID=\"%s\"\n", credentials.AccessKeyID)
	fmt.Printf("$env:AWS_SECRET_ACCESS_KEY=\"%s\"\n", credentials.SecretAccessKey)
	fmt.Printf("$env:AWS_SESSION_TOKEN=\"%s\"\n", credentials.SessionToken)
	fmt.Printf("$env:AWS_SECURITY_TOKEN=\"%s\"\n", credentials.SessionToken)
	fmt.Printf("$env:AWS_ASSUMED_ROLE=\"%s\"\n", role)
	fmt.Printf("$env:AWS_SESSION_EXPIRATION=\"%s\"\n", credentials.Expires.Format(time.RFC3339))
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# %s | Invoke-Expression \n", strings.Join(os.Args, " "))
}
