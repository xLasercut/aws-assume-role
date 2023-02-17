package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

func usage() {
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

func parseArgs() (string, time.Duration, string, string) {
	homeDir, _ := os.UserHomeDir()
	defaultCredentialsFilepath := path.Join(homeDir, ".aws", "credentials")

	duration := flag.Duration("duration", 3600, "The duration that the credentials will be valid for in seconds.")
	credentialsFilepath := flag.String("credentials-path", defaultCredentialsFilepath, "The absolute filepath to your aws credentials file.")
	format := flag.String("format", defaultFormat(), "Format can be \"bash\", \"fish\" or \"powershell\".")

	flag.Parse()
	argv := flag.Args()

	if len(argv) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	validateFormat(*format)

	return argv[0], *duration, *credentialsFilepath, *format
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
