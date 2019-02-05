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
		failf("Failed to get current directory, error: %s", err)
	}

	x := xcTestHTMLReport{
		verbose:           cfg.Verbose,
		generateJUnit:     cfg.GenerateJUnit,
		resultBundlePaths: testResults,
	}

	//
	// Install
	{
		log.Infof("Install XCTestHTMLReport via brew")

		cmd := x.installCmd()
		cmd.SetDir(dir).
			SetStdout(os.Stdout).
			SetStderr(os.Stderr)

		if err := cmd.Run(); err != nil {
			failf("Failed to install XCTestHTMLReport, error: %s", err)
		}

		log.Successf("XCTestHTMLReport successfully installed")
		fmt.Println()
	}

	// Generate report
	//
	{
		info := "Generating html report"
		if cfg.GenerateJUnit {
			info = "Generating html and JUnit report"
		}
		log.Infof(info)

		cmd := x.convertToHTMReportCmd()
		cmd.SetDir(dir).
			SetStdout(os.Stdout).
			SetStderr(os.Stderr)

		if err := cmd.Run(); err != nil {
			failf("Failed to generate XCTestHTMLReport, error: %s", err)
		}

		info = "Html report successfully generated"
		if cfg.GenerateJUnit {
			info = "Html and JUnit reports successfully generated"
		}

		log.Successf(info)
		fmt.Println()
	}
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	log.Warnf("For more details you can enable the debug logs by turning on the verbose step input.")
	os.Exit(1)
}
