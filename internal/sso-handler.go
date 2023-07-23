package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func AssumeBaseRoleSso(profile AwsProfile) AwsCredentials {
	StartAssumeRoleMessage(profile)
	clientState := &SsoClientState{
		AccessToken:             "",
		ClientId:                "",
		ClientSecret:            "",
		DeviceCode:              "",
		VerificationUriComplete: "",
		StartUrl:                profile.SsoStartUrl,
	}
	oidcClient, ssoClient := initClients(profile.SsoRegion)
	registerClient(oidcClient, clientState)
	startDeviceAuthorization(oidcClient, clientState)
	getAccessToken(oidcClient, clientState)
	roleCreds := getRoleCreds(ssoClient, clientState, profile)
	awsCreds := formatRoleCreds(roleCreds)
	SuccessAssumeRoleMessage(awsCreds)
	return awsCreds
}

func initClients(region string) (*ssooidc.Client, *sso.Client) {
	oidcClientOptions := ssooidc.Options{
		Region: region,
	}
	oidcClient := ssooidc.New(oidcClientOptions)

	ssoClientOptions := sso.Options{
		Region: region,
	}
	ssoClient := sso.New(ssoClientOptions)

	return oidcClient, ssoClient
}

func registerClient(client *ssooidc.Client, clientState *SsoClientState) {
	ctx := context.TODO()
	clientName := "aws-assume-role"
	clientType := "public"

	registerClientInput := ssooidc.RegisterClientInput{
		ClientName: &clientName,
		ClientType: &clientType,
	}
	registerClientOutput, err := client.RegisterClient(ctx, &registerClientInput)
	CheckError(err, "Unable to register sso client")
	clientState.ClientId = *registerClientOutput.ClientId
	clientState.ClientSecret = *registerClientOutput.ClientSecret
}

func startDeviceAuthorization(client *ssooidc.Client, clientState *SsoClientState) {
	ctx := context.TODO()
	deviceAuthorizationInput := ssooidc.StartDeviceAuthorizationInput{
		ClientId:     &clientState.ClientId,
		ClientSecret: &clientState.ClientSecret,
		StartUrl:     &clientState.StartUrl,
	}
	deviceAuthorizationOutput, err := client.StartDeviceAuthorization(ctx, &deviceAuthorizationInput)
	CheckError(err, "Unable to start device authorization")
	openUrlInBrowser(*deviceAuthorizationOutput.VerificationUriComplete)
	clientState.DeviceCode = *deviceAuthorizationOutput.DeviceCode
	clientState.VerificationUriComplete = *deviceAuthorizationOutput.VerificationUriComplete
}

func getAccessToken(client *ssooidc.Client, clientState *SsoClientState) {
	var timeoutCount int
	ctx := context.TODO()
	grantType := "urn:ietf:params:oauth:grant-type:device_code"
	tokenInput := ssooidc.CreateTokenInput{
		ClientId:     &clientState.ClientId,
		ClientSecret: &clientState.ClientSecret,
		GrantType:    &grantType,
		DeviceCode:   &clientState.DeviceCode,
	}
	timeoutCount = 0
	for {
		if timeoutCount > 20 {
			err := errors.New("timed out waiting for authorization")
			CheckError(err, "Could not fetch SSO access token")
		}
		createTokenOutput, err := client.CreateToken(ctx, &tokenInput)
		if isWaitingForAuthorization(err) {
			fmt.Fprint(os.Stderr, "Still waiting for authorization...\n")
			time.Sleep(3 * time.Second)
			timeoutCount++
			continue
		}
		clientState.AccessToken = *createTokenOutput.AccessToken
		break
	}
}

func isWaitingForAuthorization(err error) bool {
	if err != nil && strings.Contains(err.Error(), "AuthorizationPendingException") {
		return true
	}
	CheckError(err, "Could not fetch SSO access token")
	return false
}

func formatRoleCreds(roleCreds types.RoleCredentials) AwsCredentials {
	return AwsCredentials{
		AccessKeyId:     *roleCreds.AccessKeyId,
		SecretAccessKey: *roleCreds.SecretAccessKey,
		ExpiresAt:       time.Unix(roleCreds.Expiration/1000, 0),
		SessionToken:    *roleCreds.SessionToken,
	}
}

func getRoleCreds(client *sso.Client, clientState *SsoClientState, profile AwsProfile) types.RoleCredentials {
	ctx := context.TODO()
	roleCredentialsInput := sso.GetRoleCredentialsInput{
		AccessToken: &clientState.AccessToken,
		AccountId:   &profile.SsoAccountId,
		RoleName:    &profile.SsoRoleName,
	}
	output, err := client.GetRoleCredentials(ctx, &roleCredentialsInput)
	CheckError(err, "Could not assume role via sso")
	return *output.RoleCredentials
}

func openUrlInBrowser(url string) {
	fmt.Fprint(os.Stderr, "Please follow instructions for SSO on the browser\n")
	var err error
	osName := determineOsName()

	switch osName {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "wsl":
		err = exec.Command("wslview", url).Start()
	default:
		err = fmt.Errorf("unsupported platform. Please open the URL manually: %s", url)
	}
	CheckError(err, "Could not open browser")
}

func determineOsName() string {
	if isWindowsSubsystemForLinuxOS() {
		return "wsl"
	}
	return runtime.GOOS
}

// isWindowsSubsystemForLinuxOS determines if the program is running on WSL
// Returns true if the OS is running in WSL, false if not.
// see https://github.com/microsoft/WSL/issues/423#issuecomment-844418910
func isWindowsSubsystemForLinuxOS() bool {
	bytes, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err == nil {
		osInfo := strings.ToLower(string(bytes))
		return strings.Contains(osInfo, "wsl")
	}
	return false
}
