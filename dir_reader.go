package gsnake
import (
    "sync"
    "container/list"
    "github.com/golang/glog"
    "time"
    "sync/atomic"
)

type DirReader struct {
    dir string


    fr FileReader

    waiting int32
    wakeup chan int
    currentReadingFile string

    mutex sync.Mutex
    files *list.List // The files to be reading
}

func NewPathReader(dir string) (*DirReader, error) {
    r := &DirReader{}
    r.dir = dir
    r.waiting = 0
    r.files = list.New()
    r.wakeup = make(chan int)
    r.fr = r.createReader()
    return r, nil
}

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
    if r.currentReadingFile == file && atomic.LoadInt32(&r.waiting) > 0 {
        glog.Infof("send kModify signal")
        r.wakeup <- kModify
    } else {
        glog.Infof("do not need to send kModify signal")
    }
    return nil
}

func (r *DirReader) OnFileCreated(file string) (err error) {
    r.add(file)
    if atomic.LoadInt32(&r.waiting) > 0 && r.files.Len() == 1 {
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

func (r *DirReader) createReader() FileReader {
    if *reader_type == "PTailReader" || *reader_type == "GzipReader" {
        return NewTextFileTailReader(r)
    } else if *reader_type == "PcapReader" {
        return NewPcapFileReader()
    }

    return nil
}

func (r *DirReader) Read() (err error) {
    glog.Infof("Starting to read files ...")
    for {
        if r.files.Len() == 0 {
            glog.Infof("No more files. Waiting ...")
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

        r.currentReadingFile = file
        startTime := time.Now()
        glog.Infof("Begin to process file %v", file)
        r.fr.ReadFile(file, 0)
        glog.Infof("Finished to process file %v", r.currentReadingFile)
        dispatcher.status.OnFileProcessingFinished(r.currentReadingFile, startTime)
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

func (r *DirReader) GetPendingFileCount() int {
    r.mutex.Lock()
    c := r.files.Len()
    r.mutex.Unlock()
    return c
}

func (r *DirReader) Wait() int {
    atomic.AddInt32(&r.waiting, 1)
    event := <-r.wakeup
    atomic.AddInt32(&r.waiting, -1)
    return event
}