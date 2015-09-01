package freader

import (
    "testing"
    "github.com/bmizerany/assert"
    "runtime"
    "path/filepath"
)

func TestGetAbsPath(t *testing.T) {
    var dir string
    if runtime.GOOS == "windows" {
        dir = "c:/"
        assert.Equal(t, filepath.IsAbs(dir), true)
    } else {
        dir = "/etc"
        assert.Equal(t, filepath.IsAbs(dir), true)
    }
    assert.Equal(t, GetAbsPath(dir), dir)
}

func TestLookupFiles(t *testing.T) {
    files, _ := LookupFiles(".", "*.go")
    found := false
    for _, f := range files {
        if f == "xmain.go" {
            found = true
        }
    }
    assert.Equal(t, found, true)
}