package internal

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path"
	"runtime"
	"strings"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <role> [<command> <args...>]\n", os.Args[0])
	flag.PrintDefaults()
}

func defaultFormat() string {
	shell := os.Getenv("SHELL")

	switch runtime.GOOS {
	case "windows":
		if shell == "" {
			return "powershell"
		}
		fallthrough
	default:
		if strings.HasSuffix(shell, "fish") {
			return "fish"
		}
		return "bash"
	}
}

func ParseArgs() (string, int, *ini.File, string) {
	homeDir, _ := os.UserHomeDir()
	defaultCredentialsFilepath := path.Join(homeDir, ".aws", "credentials")
	defaultConfigFilepath := path.Join(homeDir, ".aws", "config")

	duration := flag.Int("duration", 3600, "The duration that the credentials will be valid for in seconds.")
	credentialsFilepath := flag.String("credentials-path", defaultCredentialsFilepath, "The absolute filepath to your aws credentials file.")
	configFilepath := flag.String("config-path", defaultConfigFilepath, "The absolute filepath to your aws config file.")
	format := flag.String("format", defaultFormat(), "Format can be \"bash\", \"fish\" or \"powershell\".")
	list := flag.Bool("list", false, "Show list of available roles.")

	flag.Parse()
	argv := flag.Args()

	validateFormat(*format)

	awsConfigFiles, err := ini.LooseLoad(*credentialsFilepath, *configFilepath)
	CheckError(err, "Could not load aws credentials or config file")

	if *list {
		fmt.Fprintf(os.Stderr, "Available AWS roles:\n")
		for _, name := range awsConfigFiles.SectionStrings() {
			fmt.Fprintf(os.Stderr, "%s\n", name)
		}
		os.Exit(0)
	}

	if len(argv) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	return argv[0], *duration, awsConfigFiles, *format
}

func validateFormat(format string) {
	validFormats := []string{"bash", "fish", "powershell"}

	for _, validFormat := range validFormats {
		if validFormat == format {
			return
		}
	}

	flag.Usage()
	os.Exit(1)
}
