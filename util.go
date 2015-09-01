package gsnake
import (
    "path"
    "os"
    "path/filepath"
    "os/exec"
    "strings"
)

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

// GetAbsPath gets the absolute path of the giving path p
func GetAbsPath(p string) string {
    if filepath.IsAbs(p) {
        return p
    }

    file, _ := exec.LookPath(os.Args[0])
    exePath, _ := filepath.Abs(file)
    dir := filepath.Dir(exePath)
    fullPath := path.Join(dir, p)
    fullPath, _ = filepath.Abs(fullPath)
    return strings.TrimRight(fullPath, "/\\")
}

// IsDir returns true if given path is a directory,
// or returns false when it's a file or does not exist.
func IsDir(dir string) bool {
    f, e := os.Stat(dir)
    if e != nil {
        return false
    }
    return f.IsDir()
}

func LookupFiles(dir string, pattern string) ([]string, error) {
    var files []string = make([]string, 0, 5)

    err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
        if f == nil {
            return err
        }

        if f.IsDir() {
            return nil
        }

        if ok, err := filepath.Match(pattern, f.Name()); err != nil {
            return err
        } else if ok {
            files = append(files, path)
        }
        return nil
    })

    return files, err
}
