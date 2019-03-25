package main

import (
	"reflect"
	"testing"
)

func Test_convertToHTMReportArgs(t *testing.T) {
	tests := []struct {
		name string
		x    xcTestHTMLReport
		want []string
	}{
		//
		// One test result
		{
			name: "Generate one html report",
			x: xcTestHTMLReport{
				verbose:       false,
				generateJUnit: false,
				resultBundlePaths: []string{
					"./Test.xcresult",
				},
			},
			want: []string{
				"-r",
				"./Test.xcresult",
			},
		},
		{
			name: "Generate one html report with verbose option",
			x: xcTestHTMLReport{
				verbose:       true,
				generateJUnit: false,
				resultBundlePaths: []string{
					"./Test.xcresult",
				},
			},
			want: []string{
				"-r",
				"./Test.xcresult",
				"-v",
			},
		},
		{
			name: "Generate one html & JUnit report",
			x: xcTestHTMLReport{
				verbose:       false,
				generateJUnit: true,
				resultBundlePaths: []string{
					"./Test.xcresult",
				},
			},
			want: []string{
				"-r",
				"./Test.xcresult",
				"-j",
			},
		},
		{
			name: "Generate one html & JUnit report with verbose option",
			x: xcTestHTMLReport{
				verbose:       true,
				generateJUnit: true,
				resultBundlePaths: []string{
					"./Test.xcresult",
				},
			},
			want: []string{
				"-r",
				"./Test.xcresult",
				"-j",
				"-v",
			},
		},

		//
		// Multiple test results
		{
			name: "Generate multiple html reports",
			x: xcTestHTMLReport{
				verbose:       false,
				generateJUnit: false,
				resultBundlePaths: []string{
					"./Test.xcresult",
					"./Test_2.xcresult",
					"./Test_3.xcresult",
				},
			},
			want: []string{
				"-r",
				"./Test.xcresult",
				"-r",
				"./Test_2.xcresult",
				"-r",
				"./Test_3.xcresult",
			},
		},
		{
			name: "Generate multiple html reports with verbose option",
			x: xcTestHTMLReport{
				verbose:       true,
				generateJUnit: false,
				resultBundlePaths: []string{
					"./Test.xcresult",
					"./Test_2.xcresult",
					"./Test_3.xcresult",
				},
			},
			want: []string{
				"-r",
				"./Test.xcresult",
				"-r",
				"./Test_2.xcresult",
				"-r",
				"./Test_3.xcresult",
				"-v",
			},
		},
		{
			name: "Generate multiple html & JUnit reports",
			x: xcTestHTMLReport{
				verbose:       false,
				generateJUnit: true,
				resultBundlePaths: []string{
					"./Test.xcresult",
					"./Test_2.xcresult",
					"./Test_3.xcresult",
				},
			},
			want: []string{
				"-r",
				"./Test.xcresult",
				"-r",
				"./Test_2.xcresult",
				"-r",
				"./Test_3.xcresult",
				"-j",
			},
		},
		{
			name: "Generate multiple html & JUnit reports with verbose option",
			x: xcTestHTMLReport{
				verbose:       true,
				generateJUnit: true,
				resultBundlePaths: []string{
					"./Test.xcresult",
					"./Test_2.xcresult",
					"./Test_3.xcresult",
				},
			},
			want: []string{
				"-r",
				"./Test.xcresult",
				"-r",
				"./Test_2.xcresult",
				"-r",
				"./Test_3.xcresult",
				"-j",
				"-v",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToHTMReportArgs(tt.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToHTMReportArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
