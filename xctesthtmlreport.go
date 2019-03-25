package main

import "github.com/bitrise-io/go-utils/command"

const xcHTMLReportRepository string = "https://raw.githubusercontent.com/TitouanVanBelle/XCTestHTMLReport/develop/xchtmlreport.rb"

type xcTestHTMLReport struct {
	verbose           bool
	generateJUnit     bool
	resultBundlePaths []string
}

//
// Reciever methods

func (xcTestHTMLReport) installCmd() *command.Model {
	return command.New("brew", "install", xcHTMLReportRepository)
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
