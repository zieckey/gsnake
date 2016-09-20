package gsnake

import (
    "testing"
    "os"
    "path/filepath"
    "strconv"
    "github.com/golang/glog"
    "time"
    "log"
    "flag"
)

type MyTestTextModule struct {
    index int
    fileCount int
    lineCount int
    dispatcher *Dispatcher
}

func (m *MyTestTextModule) OnRecord(buf []byte) {
    record := string(buf)
    expectedText := strconv.Itoa(m.index)
    if expectedText != record {
        glog.Fatalf("index=%d read data is [%v]. WRONG!!!", m.index, record)
    }
    m.index++

    if m.index == m.fileCount * m.lineCount {
        dispatcher.Stop()
    }
}

func GenerateTextFiles(testDir string, fileCount, lineCount int) {
    time.Sleep(time.Second)
    index := 0
    for f := 0; f < fileCount; f++ {
        path := filepath.Join(testDir, strconv.Itoa(f) + ".txt")
        fp, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR, 0644)
        if err != nil {
            log.Printf("dddddddddddddddddddddddd %v", err.Error())
            glog.Fatalf("OpenFile %v failed : %v", path, err.Error())
        }
        log.Printf("Create file %v OK", path)
        for line := 0; line < lineCount; line++ {
            s := strconv.Itoa(index)
            //log.Printf("Write [%v] to %v", s, path)
            index++
            fp.Write([]byte(s))
            fp.Write([]byte("\n"))
            time.Sleep(time.Millisecond)
        }
        fp.Close()
    }
}

func TestPtailReader(t *testing.T) {
    testDir := "test_data"
    os.RemoveAll(testDir)
    os.Mkdir(testDir, 0755)
    defer os.RemoveAll(testDir)

    *dir = testDir
    *statusFile = filepath.Join(testDir, "status.dat")
    *filePattern = "*.txt"
    *readerType = "PTailReader"
    flag.Set("stderrthreshold", "0")

    dispatcher, err := New()
    if err != nil {
        t.Error(err.Error())
    }

    m := &MyTestTextModule{
        index : 0,
        dispatcher: dispatcher,
        fileCount: 10,
        lineCount: 100,
    }
    go GenerateTextFiles(testDir, m.fileCount, m.lineCount)
    m.dispatcher.Register(m)
    m.dispatcher.Run()
    if m.index != m.fileCount * m.lineCount {
        t.Errorf("count ERROR")
    }
}
