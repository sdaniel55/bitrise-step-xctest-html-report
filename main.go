package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Config ...
type Config struct {
	// XCHTMLReport
	TestResults   string `env:"test_result_path,required"`
	GenerateJUnit bool   `env:"generate_junit,required"`

	// Common
	OutputDir string `env:"output_dir,dir"`

	// Log
	Verbose bool `env:"verbose,required"`
}

// exportReports will search for the generated html and junit report
// it will copy them to the provided outputDir
// export the reports' path to env via envman
func exportReports(pth, outputDir string, generateJUnit bool, errors *[]error) (string, string, error) {

	// Find the generated reports
	htmlReportPth := path.Join(pth, "index.html")
	junitPth := path.Join(pth, "report.junit")

	//
	// Check the report files
	{
		if exists, err := pathutil.IsPathExists(htmlReportPth); err != nil {
			return "", "", fmt.Errorf("Failed to check if path exists, error: %s", err)
		} else if !exists {
			*errors = append(*errors, fmt.Errorf("HTML report does not exists in: %s", htmlReportPth))
		}

		if generateJUnit {
			if exists, err := pathutil.IsPathExists(junitPth); err != nil {
				return "", "", fmt.Errorf("Failed to check if path exists, error: %s", err)
			} else if !exists {
				*errors = append(*errors, fmt.Errorf("JUnit report does not exists in: %s", junitPth))
			}
		}
	}

	//
	// Copy reports
	var exportedJUnitReportPth string
	exportedHTMLReportPth := copy(htmlReportPth, outputDir, errors)
	if generateJUnit {
		exportedJUnitReportPth = copy(junitPth, outputDir, errors)
	}

	//
	// Export reports
	{
		if err := tools.ExportEnvironmentWithEnvman("XC_HTML_REPORT", exportedHTMLReportPth); err != nil {
			return "", "", fmt.Errorf("Failed to generate output - %s", "XC_HTML_REPORT")
		}

		if generateJUnit {
			if err := tools.ExportEnvironmentWithEnvman("XC_JUNIT_REPORT", exportedJUnitReportPth); err != nil {
				return "", "", fmt.Errorf("Failed to generate output - %s", "XC_JUNIT_REPORT")
			}
		}
	}
	return exportedHTMLReportPth, exportedJUnitReportPth, nil
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()

	testResults := strings.Split(strings.TrimRight(cfg.TestResults, "\n"), "\n")
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

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("Failed to install XCTestHTMLReport, error: %s", err)
		}

		log.Successf("XCTestHTMLReport successfully installed")
		fmt.Println()
	}

	//
	// Generate reports
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

		log.Printf("$ %s", cmd.PrintableCommandArgs())

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

	//
	// Export generated reports
	log.Infof("Exporting generated reports")

	var errors []error
	htmlReport, junitReport, err := exportReports(x.resultBundlePaths[0], cfg.OutputDir, x.generateJUnit, &errors)
	if err != nil {
		failf("Failed to export the generated reports, error: %s", err)
	}

	// Log envs
	log.Successf("XC_HTML_REPORT => %s", htmlReport)
	if x.generateJUnit {
		log.Successf("XC_JUNIT_REPORT => %s", junitReport)
	}

	// Log errors
	if errors != nil {
		log.Warnf("Errors during the step:\n")
		for _, err := range errors {
			log.Errorf(err.Error())
		}
	}
}
