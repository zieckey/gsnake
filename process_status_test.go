package freader
import (
    "testing"
    "fmt"
    "github.com/bmizerany/assert"
    "time"
    "strconv"
    "os/exec"
)


func TestSscanfAndTime(t *testing.T) {
    var start,end,path string
    line := "2015/08/28-20:42:12.1231 2015/08/28-20:43:23.3123 /home/s/data/log/xxx.log"
    fmt.Sscanf(line, "%s %s %s", &start, &end, &path)
    assert.Equal(t, start, "2015/08/28-20:42:12.1231")
    assert.Equal(t, end, "2015/08/28-20:43:23.3123")
    assert.Equal(t, path, "/home/s/data/log/xxx.log")

    st, _ := time.Parse("2006/01/02-15:04:05.99999", start)

    assert.Equal(t, st.Year(), 2015)
    assert.Equal(t, st.Month(), time.August)
    assert.Equal(t, st.Day(), 28)
    assert.Equal(t, st.Hour(), 20)
    assert.Equal(t, st.Minute(), 42)
    assert.Equal(t, st.Second(), 12)

    s := st.Format("2006/01/02-15:04:05.99999")
    assert.Equal(t, s, start)
}

func TestNewProcessStatus(t *testing.T) {
    t1, err := time.Parse("2006/01/02-15:04:05.99999", "2015/08/28-23:48:34.7161")
    t2, err := time.Parse("2006/01/02-15:04:05.99999", "2015/08/28-23:48:36.7161")
    f1 := "1.txt"
    f2 := "2.txt"
    statusFile := ".status." + strconv.FormatInt(time.Now().UnixNano(), 10)
    ps, err := NewProcessStatus(statusFile)
    assert.Equal(t, err, nil)
    ps.OnFileProcessingFinished(f1, t1)
    time.Sleep(time.Second)
    ps.OnFileProcessingFinished(f2, t2)
    ps.Close()

    ps, err = NewProcessStatus(statusFile)
    assert.Equal(t, err, nil)
    files := ps.GetProcessedFiles()
    ts, ok := files[f1]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t1)
    ts, ok = files[f2]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t2)
    ps.Close()

    t3, err := time.Parse("2006/01/02-15:04:05.99999", "2015/08/28-23:48:38.7161")
    f3 := "3.txt"
    ps, err = NewProcessStatus(statusFile)
    assert.Equal(t, err, nil)
    ps.OnFileProcessingFinished(f3, t3)
    files = ps.GetProcessedFiles()
    ts, ok = files[f1]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t1)
    ts, ok = files[f2]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t2)
    ts, ok = files[f3]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t3)
    ps.Close()

    t4, err := time.Parse("2006/01/02-15:04:05.99999", "2015/08/28-23:48:39.1000")
    f4 := "4.txt"
    ps, err = NewProcessStatus(statusFile)
    assert.Equal(t, err, nil)
    ps.OnFileProcessingFinished(f4, t4)
    files = ps.GetProcessedFiles()
    ts, ok = files[f1]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t1)
    ts, ok = files[f2]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t2)
    ts, ok = files[f3]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t3)
    ts, ok = files[f4]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t4)
    ps.OnFileDeleted(f1)
    ps.Close()

    ps, err = NewProcessStatus(statusFile)
    assert.Equal(t, err, nil)
    ps.OnFileProcessingFinished(f4, t4)
    files = ps.GetProcessedFiles()
    ts, ok = files[f1]
    assert.Equal(t, ok, false)
    ts, ok = files[f2]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t2)
    ts, ok = files[f3]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t3)
    ts, ok = files[f4]
    assert.Equal(t, ok, true)
    ts.Start.Equal(t4)
    ps.Close()

    cmd := exec.Command("rm", "-f", ".status.*")
    cmd.Run()
}
