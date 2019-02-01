package main

import (
	"fmt"
	"os"
)

var resultBundlePaths = []string{"/Users/akosbirmacher/Develop/Bitrise/github/BirmacherAkos/bitrise-samples/apps/ios/xcode-10/default/Xcode-10_default/ddata/Test.xcresult"}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current directory, error: %s", err))
	}

	x := xcTestHTMLReport{
		verbose:           true,
		generateJUnit:     true,
		resultBundlePaths: resultBundlePaths,
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
