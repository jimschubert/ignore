package ignore

import (
	`fmt`
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jimschubert/ignore/test"
)

type AllowTestCondition struct {
	File       string
	Allows     bool
	WantErr    bool
	TestReason string
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
				{File: "cmake-build-anything", Allows: true, TestReason: "should not exclude a file matching a directory rule"},
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
		{
			name:       "escaped special characters",
			ignoreFile: "gitignore_escaped_chars",
			conditions: []AllowTestCondition{
				{File: "important.txt", Allows: true, TestReason: `Doesn't match "\*important.txt" (escaped *)`},
				{File: "*important.txt", Allows: false, TestReason: `Matches \*important.txt`},
				{File: "file?.md", Allows: false, TestReason: `Filename includes a ?, not ? match pattern (file\?.md)`},
				{File: "fileX.md", Allows: true, TestReason: `Doesn't match file\?.md (i.e. ? not wildcard)`},
			},
			wantErr: false,
		},
		{
			name:       "trailing spaces handling",
			ignoreFile: "gitignore_trailing_spaces",
			conditions: []AllowTestCondition{
				{File: "trimmed.txt", Allows: false},
				{File: "spaced.txt ", Allows: false, TestReason: "unlikely, but I've seen weird whitespace handling in filenames"},
			},
			wantErr: false,
		},
		{
			name:       "negation patterns",
			ignoreFile: "gitignore_negation",
			conditions: []AllowTestCondition{
				{File: "debug.log", Allows: false, TestReason: "excluded by *.log"},
				{File: "important.log", Allows: true, TestReason: "re-included by !important.log"},
				{File: "build/output.txt", Allows: false, TestReason: "excluded by build/"},
				{File: "build/keep.txt", Allows: true, TestReason: "Actually re-included (current impl doesn't check parent)"},
			},
			wantErr: false,
		},
		{
			name:       "double asterisk patterns",
			ignoreFile: "gitignore_double_asterisk",
			conditions: []AllowTestCondition{
				{File: "test.tmp", Allows: false, TestReason: "matches **/*.tmp (** means 0..n)"},
				{File: "foo/test.tmp", Allows: false, TestReason: "matches **/*.tmp"},
				{File: "foo/bar/test.tmp", Allows: false, TestReason: "matches **/*.tmp"},
				{File: "cache/file.txt", Allows: false, TestReason: "matches cache/**"},
				{File: "cache", Allows: true, TestReason: "should not exclude a file matching a directory rule"},
				{File: "cache.md", Allows: true, TestReason: "should not exclude a file matching unrelated pattern"},
				{File: "src/cache/file.txt", Allows: true, TestReason: "should not exclude a file matching a directory rule if not at root"},
				{File: "src/cache", Allows: true, TestReason: "should not exclude a file matching a directory rule if not at root"},
				{File: "src/cache.md", Allows: true, TestReason: "should not exclude a file matching unrelated pattern if not at root"},
			},
			wantErr: false,
		},
		{
			name:       "character ranges and wildcards",
			ignoreFile: "gitignore_ranges_wildcards",
			conditions: []AllowTestCondition{
				{File: "file0.txt", Allows: false, TestReason: "matches file[0-9].txt"},
				{File: "file5.txt", Allows: false, TestReason: "matches file[0-9].txt"},
				{File: "fileA.txt", Allows: true, TestReason: "doesn't match file[0-9].txt"},
				{File: "testX.txt", Allows: false, TestReason: "matches test?.txt (? matches single char)"},
				{File: "test/.txt", Allows: true, TestReason: "? doesn't match /"},
			},
			wantErr: false,
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
					reason := ""
					if condition.TestReason != "" {
						reason = fmt.Sprintf(" TestReason=[%s]", condition.TestReason)
					}
					t.Errorf("NewProcessor() condition for path '%s': expected Allows=%v but got %v%s (ignore file %s)", condition.File, condition.Allows, isAllowed, reason, tt.ignoreFile)
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
