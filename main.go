package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	TestResults   string `env:"test_result_path,required"`
	GenerateJUnit bool   `env:"verbose,required"`
	Verbose       bool   `env:"verbose,required"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()

	testResults := strings.Split(cfg.TestResults, "\n")
	log.SetEnableDebugLog(cfg.Verbose)

	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current directory, error: %s", err))
	}

	x := xcTestHTMLReport{
		verbose:           cfg.Verbose,
		generateJUnit:     cfg.GenerateJUnit,
		resultBundlePaths: testResults,
	}

	//
	// Install
	{
		cmd := x.installCmd()
		cmd.SetDir(dir).
			SetStdout(os.Stdout).
			SetStderr(os.Stderr)

		if err := cmd.Run(); err != nil {
			panic(fmt.Sprintf("Failed to install XCTestHTMLReport, error: %s", err))
		}
	}

	// Generate report
	//
	{
		cmd := x.convertToHTMReportCmd()
		cmd.SetDir(dir).
			SetStdout(os.Stdout).
			SetStderr(os.Stderr)

		if err := cmd.Run(); err != nil {
			panic(fmt.Sprintf("Failed to generate XCTestHTMLReport, error: %s", err))
		}
	}
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	log.Warnf("For more details you can enable the debug logs by turning on the verbose step input.")
	os.Exit(1)
}
