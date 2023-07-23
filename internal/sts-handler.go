package internal

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"os"
	"time"
)

func AssumeBaseRoleSts(profile AwsProfile, duration int) AwsCredentials {
	creds := AwsCredentials{
		AccessKeyId:     profile.AwsAccessKeyId,
		SecretAccessKey: profile.AwsSecretAccessKey,
	}
	return AssumeRoleSts(profile, creds, duration)
}

func AssumeRoleSts(profile AwsProfile, currentCreds AwsCredentials, duration int) AwsCredentials {
	StartAssumeRoleMessage(profile)
	stsClient := initClient(currentCreds)
	provider := stscreds.NewAssumeRoleProvider(stsClient, profile.RoleArn, func(p *stscreds.AssumeRoleOptions) {
		p.TokenProvider = tokenProvider
		p.Duration = time.Duration(duration) * time.Second
		p.RoleSessionName = fmt.Sprintf("%v-%v", profile.Name, profile.SourceProfile)

		if profile.MfaSerial != "" {
			p.SerialNumber = &profile.MfaSerial
		}
	})
	credsCache := aws.NewCredentialsCache(provider)
	cacheCreds := getCacheCredentials(credsCache)
	awsCreds := formatCacheCreds(cacheCreds)
	SuccessAssumeRoleMessage(awsCreds)
	return awsCreds
}

func initClient(currentCreds AwsCredentials) *sts.Client {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				currentCreds.AccessKeyId,
				currentCreds.SecretAccessKey,
				currentCreds.SessionToken,
			),
		),
	)
	CheckError(err, "Could not load sts client config")
	return sts.NewFromConfig(cfg)
}

func getCacheCredentials(credsCache *aws.CredentialsCache) aws.Credentials {
	ctx := context.TODO()
	creds, err := credsCache.Retrieve(ctx)
	CheckError(err, "Could not assume role")
	return creds
}

func formatCacheCreds(cacheCreds aws.Credentials) AwsCredentials {
	return AwsCredentials{
		AccessKeyId:     cacheCreds.AccessKeyID,
		SecretAccessKey: cacheCreds.SecretAccessKey,
		ExpiresAt:       cacheCreds.Expires,
		SessionToken:    cacheCreds.SessionToken,
	}
}

func tokenProvider() (string, error) {
	var v string
	fmt.Fprint(os.Stderr, "Assume Role MFA token code: ")
	_, err := fmt.Scanln(&v)

	return v, err
}
