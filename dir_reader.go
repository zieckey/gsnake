package freader
import (
    "sync"
    "container/list"
    "io"
    "github.com/golang/glog"
    "time"
)

type DirReader struct {
    dir string


    fr FileReader

    waiting bool    //FIXME use atomic variable ??
    wakeup chan int

    currentReadingFile string

    mutex sync.Mutex
    files *list.List // The files to be reading
}

func NewPathReader(dir string) (*DirReader, error) {
    r := &DirReader{}
    r.dir = dir
    r.waiting = false
    r.files = list.New()
    r.wakeup = make(chan int)
    r.fr = createReader()
    return r, nil
}

//func (r *PathReader) ReadLine() (line string, err error)  {
//
//    return line, err
//}
//
//func (r *PathReader) Read(p []byte) (n int, err error) {
//    return 0, nil
//}

func (r *DirReader) add(file string) (err error) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    r.files.PushBack(file)
    return err
}

const (
    kModify int = 1
    kCreate int = 2
)

func (r *DirReader) OnFileModified(file string) (err error) {
    if r.currentReadingFile == file && r.waiting {
        glog.Infof("send kModify signal")
        r.wakeup <- kModify
    } else {
        glog.Infof("do not need to send kModify signal")
    }
    return nil
}

func (r *DirReader) OnFileCreated(file string) (err error) {
    r.add(file)
    if r.waiting && r.files.Len() == 1 {
        /*
        r.waiting : we will send a signal only if the goroutine is waiting
        r.files.Len() == 0 : when we create more than 2 files in the same time, the waiting goroutine may be still waiting when we try to send the second signal
         */
        glog.Infof("send kCreate signal")
        r.wakeup <- kCreate
    } else {
        glog.Infof("do not need to send kCreate signal")
    }
    return nil
}

func createReader() FileReader {
    //GzipReader, PTailReader
    if *reader_type == "PTailReader" {
        return NewPTailFileReader()
    } else if *reader_type == "GzipReader" {
        return NewGzipFileReader()
    } else if *reader_type == "PcapReader" {
        return NewPcapFileReader()
    }

    return nil
}

func (r *DirReader) StartToRead() (err error) {
    glog.Infof("Starting to read files ...")
    startTime := time.Now()
    var tfr TextFileReader
    if *reader_type == "PTailReader" || *reader_type == "GzipReader" {
        var ok bool
        if tfr, ok = r.fr.(TextFileReader); !ok {
            glog.Errorf("Dynamic cast to TextFileReader failed")
            panic("ERROR")
        }
    } else if *reader_type == "PcapReader" {

    }

    for {
        if r.files.Len() == 0 {
            glog.Infof("No files. Waiting ...")
            r.Wait()
            if r.files.Len() == 0 {
                glog.Errorf("This is a logic ERROR, but we ignore it right now and lately we should review this code logic.")
                continue
            }
        }

        file := r.nextFile()
        if len(file) == 0 {
            continue
        }

        r.fr.LoadFile(file, 0)
        if len(r.currentReadingFile) > 0 {
            glog.Infof("Finished to process file %v", r.currentReadingFile)
            dispatcher.status.OnFileProcessingFinished(r.currentReadingFile, startTime)
        }
        startTime = time.Now()
        r.currentReadingFile = file
        glog.Infof("Begin to process file %v", file)

        //FIXME 这里不太符合面向对象的类/接口继承关系，但也仅仅就两个分支，这里做了简单化处理。后续如果分支太多，再好好设计继承关系重构一下。
        if tfr != nil {
            r.readTextFile(tfr, file)
        } else {
            //glog.Infof("Pcap Reader does not need do anything")
        }
    }
}

func (r *DirReader) nextFile() string {
    r.mutex.Lock()
    e := r.files.Front()
    r.files.Remove(e)
    r.mutex.Unlock()

    file, ok := e.Value.(string)
    if !ok {
        glog.Errorf("Get element from file List failed.")
        return ""
    }
    return file
}

func (r *DirReader) readTextFile(tfr TextFileReader, file string) {
    var lastLine []byte
    for {
        line, err := tfr.ReadLine()
        //glog.Infof("ReadLine: lastLine=<%s> current-read=<%s> <%v>", string(lastLine), string(line), err)
        if len(lastLine) > 0 {
            line = append(lastLine, line...)
        }

        if err == io.EOF {
            if len(line) > 0 {
                lastLine = line
            }

            // there are still files which are ready to be processed
            if r.files.Len() > 0 {
                break
            }

            // no more files. we wait this file to be updated or wait new file created
            glog.Infof("no more files, we wait this file <%v> to be updated. Waiting ...", file)
            r.Wait()
            continue
        } else if err != nil {
            glog.Errorf("Read data from <%s> failed : %v", file, err.Error())
            break
        } else {
            lastLine = []byte{}
        }

        dispatcher.textModule.OnRecord(line)
    }
}

func (r *DirReader) Wait() int {
    r.waiting = true
    event := <-r.wakeup
    r.waiting = false
    return event
}