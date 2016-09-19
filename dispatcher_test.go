package gsnake

import (
    "testing"
    "os"
    "path/filepath"
    "strconv"
    "github.com/golang/glog"
    "time"
    "log"
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
    log.Printf("Read one line [%v]", record)
    if expectedText != record {
        glog.Fatalf("index=%d read data is [%v]. WRONG!!!", m.index, record)
    }
    m.index++

    if m.index == m.fileCount * m.lineCount {
        dispatcher.Stop()
    }
}

func GenerateTextFiles(testDir string, fileCount, lineCount int) {
    log.Printf("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
    index := 0
    for f := 0; f < fileCount; f++ {
        log.Printf("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
        path := filepath.Join(testDir, strconv.Itoa(f) + ".txt")
        log.Printf("ccccccccccccccccccccccccccc")
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

    log.Printf("1111111111111111111")
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
    log.Printf("222222222222222222222222")
    go GenerateTextFiles(testDir, m.fileCount, m.lineCount)
    log.Printf("333333333333333333333")
    m.dispatcher.Register(m)
    log.Printf("4444444444444444444")
    m.dispatcher.Run()
    log.Printf("5555555555555555555555555555")
    if m.index != m.fileCount * m.lineCount {
        t.Errorf("count ERROR")
    }
}
