package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/bitrise-tools/go-steputils/tools"
)

// Config ...
type Config struct {
	// XCHTMLReport
	TestResults   string `env:"test_result_path,required"`
	GenerateJUnit bool   `env:"verbose,required"`

	// Common
	OutputDir string `env:"output_dir,dir"`

	// Log
	Verbose bool `env:"verbose,required"`
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

	//
	// Generate report
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

	//
	// Find the generated reports
	htmlReportPth := path.Join(x.resultBundlePaths[0], "index.html")
	junitPth := path.Join(x.resultBundlePaths[0], "report.junit")

	//
	// Check the report files
	var errors []error
	{
		if exists, err := pathutil.IsPathExists(htmlReportPth); err != nil {
			failf("Failed to check if path exists, error: %s", err)
		} else if !exists {
			errors = append(errors, fmt.Errorf("HTML report does not exists in: %s", htmlReportPth))
		}

		if x.generateJUnit {
			if exists, err := pathutil.IsPathExists(junitPth); err != nil {
				failf("Failed to check if path exists, error: %s", err)
			} else if !exists {
				errors = append(errors, fmt.Errorf("JUnit report does not exists in: %s", junitPth))
			}
		}
	}

	//
	// Copy reports
	var exportedJunitReportPth string
	exportedHTMLReportPth := copy(htmlReportPth, cfg.OutputDir, &errors)
	if x.generateJUnit {
		exportedJunitReportPth = copy(junitPth, cfg.OutputDir, &errors)
	}

	//
	// Export reports
	if err := tools.ExportEnvironmentWithEnvman("XC_HTML_Report", exportedHTMLReportPth); err != nil {
		failf("Failed to generate output - %s", "XC_HTML_Report")
	}

	if x.generateJUnit {
		if err := tools.ExportEnvironmentWithEnvman("XC_JUnit_Report", exportedJunitReportPth); err != nil {
			failf("Failed to generate output - %s", "XC_JUnit_Report")
		}
	}

	log.Successf("XC_HTML_Report => %s", exportedHTMLReportPth)
	if x.generateJUnit {
		log.Successf("XC_JUnit_Report => %s", exportedJunitReportPth)
	}

	// Log errors
	if errors != nil {
		log.Warnf("Errors during the step:\n")
		for _, err := range errors {
			log.Errorf(err.Error())
		}
	}
}

func copy(sourcePath, outputDir string, errors *[]error) string {
	source, err := os.Open(sourcePath)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("Failed to open file, error: %s", err))
	}

	defer func() {
		if cerr := source.Close(); cerr != nil {
			*errors = append(*errors, fmt.Errorf("Failed to close file, error: %s", cerr))
		}
	}()

	destinationPath := path.Join(outputDir, path.Base(sourcePath))
	destination, err := os.Create(destinationPath)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("Failed to open file, error: %s", err))
	}

	defer func() {
		if cerr := destination.Close(); cerr != nil {
			*errors = append(*errors, fmt.Errorf("Failed to close file, error: %s", cerr))
		}
	}()

	if _, err := io.Copy(destination, source); err != nil {
		*errors = append(*errors, fmt.Errorf("Failed to copy file, error: %s", err))
	}

	return destinationPath
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	log.Warnf("For more details you can enable the debug logs by turning on the verbose step input.")
	os.Exit(1)
}
