package test

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Data(t *testing.T, file string) []byte {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal(errors.New("[helpers_test] test.Data: unable to determine testdata location"))
	}

	p := path.Join(path.Dir(filename), "..", "testdata", strings.TrimPrefix(file, "testdata"))
	contents, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func CopyToTempLocation(t *testing.T, data []byte) (fileLocation string, cleanup func()) {
	t.Helper()
	tempDir := os.TempDir()
	testHash := hash(t.Name())
	testFile := fmt.Sprintf("file-%s", testHash)
	filePath := filepath.Join(tempDir, testFile)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Fatal(err)
	}
	return filePath, func() { _ = os.RemoveAll(filePath) }
}

func hash(s string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
