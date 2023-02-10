package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/bilus/scenarigo"
	"github.com/bilus/scenarigo/cmd/scenarigo/cmd/config"
	"github.com/bilus/scenarigo/context"
	"github.com/bilus/scenarigo/reporter"
)

// ErrTestFailed is the error returned when the test failed.
var ErrTestFailed = errors.New("test failed")

var verbose bool

func init() {
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print verbose log")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:           "run",
	Short:         "run test scenarios",
	Long:          "Runs test scenarios.",
	RunE:          run,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func run(cmd *cobra.Command, args []string) error {
	opts := []func(*scenarigo.Runner) error{}
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if cfg != nil {
		if len(args) > 0 {
			cfg.Scenarios = nil
		}
		opts = append(opts, scenarigo.WithConfig(cfg))
	}
	if len(args) > 0 {
		opts = append(opts, scenarigo.WithScenarios(args...))
	}
	r, err := scenarigo.NewRunner(opts...)
	if err != nil {
		return err
	}

	reporterOpts := []reporter.Option{
		reporter.WithWriter(cmd.OutOrStdout()),
	}
	if (cfg != nil && cfg.Output.Verbose) || verbose {
		reporterOpts = append(reporterOpts, reporter.WithVerboseLog())
	}

	noColor := color.NoColor
	if cfg != nil && cfg.Output.Colored != nil {
		noColor = !*cfg.Output.Colored
	}
	if noColor {
		reporterOpts = append(reporterOpts, reporter.WithNoColor())
	}

	var reportErr error
	success := reporter.Run(
		func(rptr reporter.Reporter) {
			r.Run(context.New(rptr))
			reportErr = r.CreateTestReport(rptr)
		},
		reporterOpts...,
	)
	if reportErr != nil {
		return fmt.Errorf("failed to create test reports: %w", reportErr)
	}
	if !success {
		return ErrTestFailed
	}
	return nil
}
