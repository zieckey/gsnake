package freader
import (
    "os"
    "bufio"
    "io"
    "time"
    "fmt"
    "strings"
    "sort"
    "strconv"
    "log"
)

type FileProcessingTime struct {
    Start time.Time
    End time.Time
}

type ProcessStatus struct {
    processedFiles map[string]FileProcessingTime    // The processed files and the time when starting to process and end

    // The content format of the status file :
    //  It is a text file. Every line represents a processed file.
    //  The line has 3 part
    //      1. start processing date time
    //      2. end of processing date time
    //      3. the name of the file
    // For example: 2015/08/28-20:42:12.1231 2015/08/28-20:43:23.3123 /home/s/data/log/xxx.log
    statusFile string       // The path of the status file which used to store the status information of all processed files
    statusFileFp *os.File   // The file pointer to the status file

}

func NewProcessStatus(statusFile string) (ps *ProcessStatus, err error) {
    ps = &ProcessStatus{}
    ps.statusFile = statusFile
    ps.processedFiles = make(map[string]FileProcessingTime)

    if IsExist(statusFile) {
        ps.statusFileFp, err = os.OpenFile(statusFile, os.O_RDWR, 0755)
        if err != nil {
            fmt.Printf("open status file <%v> failed : %v\n", err.Error())
            return nil, err
        }
        if err = ps.parse(); err != nil {
            return nil, err
        }
        ps.statusFileFp.Seek(0, os.SEEK_END)
    } else {
        ps.statusFileFp, err = os.OpenFile(statusFile, os.O_CREATE | os.O_RDWR, 0755)
        if err != nil {
            fmt.Printf("open status file <%v> failed : %v\n", err.Error())
            return nil, err
        }
    }

    return ps, nil
}

func (ps *ProcessStatus) IsProcessed(file string) bool {
    _, ok := ps.processedFiles[file]
    return ok
}

func (ps *ProcessStatus) GetProcessedFiles() map[string]FileProcessingTime {
    return ps.processedFiles
}

func (ps *ProcessStatus) OnFileProcessingFinished(path string, startProcessing time.Time) {
    var t FileProcessingTime
    t.Start = startProcessing
    t.End = time.Now()
    ps.processedFiles[path] = t

    w := ps.statusFileFp
    w.WriteString(t.Start.Format("2006/01/02-15:04:05.9999 "))
    w.WriteString(t.End.Format("2006/01/02-15:04:05.9999 "))
    w.WriteString(path)
    w.WriteString("\n")
    w.Sync()
}

func (ps *ProcessStatus) OnFileDeleted(path string) {
    delete(ps.processedFiles, path)
}

func (ps *ProcessStatus) Close()  {
    defer ps.statusFileFp.Close()
    if err := ps.saveAll(); err != nil {
        panic(err.Error())
    } // flush all data to files
}

func (ps *ProcessStatus) parse() error {
    r := bufio.NewReader(ps.statusFileFp)
    for {
        line, err := r.ReadString('\n')

        if err == io.EOF {
            break
        }

        if len(line) == 0 {
            continue
        }

        line = strings.TrimSpace(line)
        var start,end,path string
        fmt.Sscanf(line, "%s %s %s", &start, &end, &path)
        var t FileProcessingTime
        t.Start, err = time.Parse("2006/01/02-15:04:05.9999", start)
        if err != nil {
            return fmt.Errorf("ERROR line <%v> %v", line, err.Error())
        }
        t.End, err = time.Parse("2006/01/02-15:04:05.9999", end)
        if err != nil {
            return fmt.Errorf("ERROR line <%v> %v", line, err.Error())
        }
        if len(path) == 0 {
            return fmt.Errorf("ERROR line <%v>, path empty", line)
        }
        ps.processedFiles[path] = t
    }
    return nil
}

type StringArray []string

func (ss StringArray) Len() int {
    return len(ss)
}

func (ss StringArray) Less(i, j int) bool {
    return ss[i] < ss[j]
}

func (ss StringArray) Swap(i, j int) {
    ss[i], ss[j] = ss[j], ss[i]
}


func (ps *ProcessStatus) saveAll() error {
    bakFilePath := ps.statusFile + ".bak." + strconv.FormatInt(time.Now().UnixNano(), 10)
    fp, err := os.OpenFile(bakFilePath, os.O_CREATE | os.O_RDWR, 0755)
    if err != nil {
        return err
    }
    ps.statusFileFp.Seek(0, os.SEEK_SET)
    io.Copy(fp, ps.statusFileFp)
    fp.Sync()
    fp.Close()

    _, err = ps.statusFileFp.Seek(0, os.SEEK_SET)
    if err != nil {
        log.Printf("Seek <%s> failed : %v\n", ps.statusFile, err.Error())
        return err
    }
    err = ps.statusFileFp.Truncate(0)
    if err != nil {
        log.Printf("Truncate <%s> failed : %v\n", ps.statusFile, err.Error())
        return err
    }
    //stat, err := ps.statusFileFp.Stat()
    //log.Printf("%v len=%v", stat.Name(), stat.Size())
    var files StringArray
    for k, _ := range ps.processedFiles {
        files = append(files, k)
    }
    sort.Sort(files)
    //log.Print(files)

    w := ps.statusFileFp
    for _, f := range files {
        if t, ok := ps.processedFiles[f]; ok {
            w.WriteString(t.Start.Format("2006/01/02-15:04:05.9999 "))
            w.WriteString(t.End.Format("2006/01/02-15:04:05.9999 "))
            w.WriteString(f)
            w.WriteString("\n")
            //log.Printf("Write <%v> to status file\n", f)
        }
    }
    w.Sync()
    return nil
}
