package main

import (
	"fmt"
	"os"
	"os/exec"
)

func checkError(err error, msg string) {
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Errors are already on Stderr.
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "%v: %v\n", msg, err)
		os.Exit(1)
	}
}
