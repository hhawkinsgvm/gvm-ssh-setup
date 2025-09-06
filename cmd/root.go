package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// REAL_HOME lets the app write to the host home when running in Docker.
	RealHome = os.Getenv("REAL_HOME")
)

var rootCmd = &cobra.Command{
	Use:   "gvm-ssh",
	Short: "GVM SSH/Git per-account setup tool",
	Long:  "Set up SSH aliases, keys, and per-directory Git identity; includes checks and connectivity tests.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}