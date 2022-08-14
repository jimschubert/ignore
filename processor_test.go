package ignore

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jimschubert/ignore/test"
)

type AllowTestCondition struct {
	File    string
	Allows  bool
	WantErr bool
}

func TestNewProcessor(t *testing.T) {
	tests := []struct {
		name       string
		ignoreFile string
		conditions []AllowTestCondition
		wantErr    bool
	}{
		{
			name:       "allows empty gitignore",
			ignoreFile: "gitignore_empty",
			conditions: []AllowTestCondition{
				{File: "a/b/c.txt", Allows: true},
			},
			wantErr: false,
		},
		{
			name:       "check simple excludes by file",
			ignoreFile: "go_jetbrains_windows",
			conditions: []AllowTestCondition{
				{File: "go.work", Allows: false},
				{File: "Thumbs.db", Allows: false},
				{File: "n", Allows: false},
				{File: ".idea/replstate.xml", Allows: false},
				{File: "prog.dll", Allows: false},
				{File: ".idea/nested/gradle.xml", Allows: false},
				{File: "cmake-build-anything/", Allows: false},
				{File: "cmake-build-anything", Allows: true}, // should not exclude a file matching a directory rule
				{File: "should_be_included", Allows: true},
				{File: "should/be/included", Allows: true},
			},
			wantErr: false,
		},
		{
			name:       "check simple excludes by directory",
			ignoreFile: "go_jetbrains_windows",
			conditions: []AllowTestCondition{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ignoreContents := test.Data(t, tt.ignoreFile)
			location, cleanup := test.CopyToTempLocation(t, ignoreContents)
			defer cleanup()
			processor, err := NewProcessor(
				WithGitignoreStrategy(),
				WithIgnoreFilePath(location),
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, condition := range tt.conditions {
				isAllowed, e := processor.AllowsFile(condition.File)
				if (e != nil) != condition.WantErr {
					t.Errorf("NewProcessor() condition for path '%s' error = %v, wantErr %v  (ignore file %s)", condition.File, e, condition.WantErr, tt.ignoreFile)
					return
				}

				if isAllowed != condition.Allows {
					t.Errorf("NewProcessor() condition for path '%s' did not allow as expected (ignore file %s)", condition.File, tt.ignoreFile)
					return
				}
			}
		})
	}
}

func ls(dir string) []string {
	contents := make([]string, 0)
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		contents = append(contents, path)
		return nil
	})
	sort.Strings(contents)
	for _, path := range contents {
		stdout(strings.ReplaceAll(path, `\`, `/`))
	}
	return contents
}

func stdout(msg string) {
	fmt.Printf("%s\n", msg)
}

func dump(path string) {
	contents, _ := os.ReadFile(path)
	stdout(string(contents))
}

func ExampleNewProcessor() {
	stdout("Given:")
	paths := ls("example/")
	stdout("\nWith ignore file:")
	dump("example/.ignore")
	stdout("Results:")

	processor, _ := NewProcessor(
		WithGitignoreStrategy(),
		WithIgnoreFilePath("example/.ignore"),
	)

	for _, path := range paths {
		if allowed, _ := processor.AllowsFile(path); allowed {
			stdout(fmt.Sprintf("✓ %s", strings.ReplaceAll(path, `\`, `/`)))
		} else {
			stdout(fmt.Sprintf("✘ %s", strings.ReplaceAll(path, `\`, `/`)))
		}
	}

	// Output:
	// Given:
	// example/
	// example/.ignore
	// example/fileA.txt
	// example/fileB.txt
	// example/first
	// example/first/contents.md
	// example/other.txt
	// example/second
	// example/second/contents.md
	// example/third
	// example/third/contents.md
	//
	// With ignore file:
	// **/file*.txt
	// !**/fileB.txt
	// **/contents.md
	// !**/second/contents.md
	//
	// Results:
	// ✓ example/
	// ✓ example/.ignore
	// ✘ example/fileA.txt
	// ✓ example/fileB.txt
	// ✓ example/first
	// ✘ example/first/contents.md
	// ✓ example/other.txt
	// ✓ example/second
	// ✓ example/second/contents.md
	// ✓ example/third
	// ✘ example/third/contents.md
}
