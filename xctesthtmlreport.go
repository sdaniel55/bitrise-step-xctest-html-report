package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
)

// InstallBranch is the selected source branch of the XCHTMLReport repository
type InstallBranch string

// enum SourceBranch
const (
	Develop InstallBranch = "develop"
	Master  InstallBranch = "master"
)

const xcHTMLReportRepository string = "https://raw.githubusercontent.com/TitouanVanBelle/XCTestHTMLReport/%s/xchtmlreport.rb"
const xcHTMLReportInstallScriptURL string = "https://raw.githubusercontent.com/TitouanVanBelle/XCTestHTMLReport/master/install.sh"
const xcHTMLReportGithubOrg string = "TitouanVanBelle"
const xcHTMLReportGithubRepo string = "XCTestHTMLReport"

type xcTestHTMLReport struct {
	verbose           bool
	generateJUnit     bool
	resultBundlePaths []string
	version           string
}

//
// Reciever methods
// Deprecated:
func (xcTestHTMLReport) installCmd(branch InstallBranch) *command.Model {
	return command.New("brew", "install", fmt.Sprintf(xcHTMLReportRepository, branch))
}

func (xcTestHTMLReport) installViaScriptCmd(version string) *command.Model {
	return command.New("/bin/sh", []string{"install.sh", version}...)
}

func (x xcTestHTMLReport) convertToHTMReportCmd() *command.Model {
	return command.New("xchtmlreport", convertToHTMReportArgs(x)...)
}

// installScript returns the install script located on the master branch of the XCTestHTMLReport repository
// https://raw.githubusercontent.com/TitouanVanBelle/XCTestHTMLReport/master/install.sh
func (x xcTestHTMLReport) installScript() (string, error) {
	resp, err := http.Get(xcHTMLReportInstallScriptURL)
	if err != nil {
		return "", fmt.Errorf("failed to call the %s, error: %v", xcHTMLReportInstallScriptURL, err)
	}
	log.Debugf("Response status: %s", resp.Status)

	defer func() {
		if cerr := resp.Body.Close(); err != nil {
			log.Warnf("Failed to close response body of %s, error: %v", xcHTMLReportInstallScriptURL, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status %v", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return "", fmt.Errorf("failed to copy the response body to a buffer, error: %v", err)
	}
	return buf.String(), nil
}

//
// Private methods

func convertToHTMReportArgs(x xcTestHTMLReport) []string {
	var args []string
	for _, path := range x.resultBundlePaths {
		args = append(args, []string{"-r", path}...)
	}

	if x.generateJUnit {
		args = append(args, "-j")
	}

	if x.verbose {
		args = append(args, "-v")
	}

	return args
}
