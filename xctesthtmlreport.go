package main

import (
	"fmt"

	"github.com/bitrise-io/go-utils/command"
)

// InstallBranch is the selected source branch of the XCHTMLReport repository
type InstallBranch string

// enum SourceBranch
const (
	Develop InstallBranch = "develop"
	Master  InstallBranch = "master"
)

const xcHTMLReportRepository string = "https://raw.githubusercontent.com/TitouanVanBelle/XCTestHTMLReport/%s/xchtmlreport.rb"

type xcTestHTMLReport struct {
	verbose           bool
	generateJUnit     bool
	resultBundlePaths []string
}

//
// Reciever methods

func (xcTestHTMLReport) installCmd(branch InstallBranch) *command.Model {
	return command.New("brew", "install", fmt.Sprintf(xcHTMLReportRepository, branch))
}

func (x xcTestHTMLReport) convertToHTMReportCmd() *command.Model {
	return command.New("xchtmlreport", convertToHTMReportArgs(x)...)
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
