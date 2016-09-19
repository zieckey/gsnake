package gsnake
import (
    "sync"
    "container/list"
    "github.com/golang/glog"
    "time"
    "sync/atomic"
)

type DirReader struct {
    dir                string

    fr                 FileReader

    waiting            int32
    wakeup             chan int
    currentReadingFile string

    Running            bool

    mutex              sync.Mutex
    files              *list.List // The files to be reading
}

func NewPathReader(dir string) (*DirReader, error) {
    r := &DirReader{}
    r.Running = true
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
        glog.Infof("The file <%s> has been modified which we are processing. And the processing goroutine is sleeping, so send kModify signal to it.", file)
        r.wakeup <- kModify
    } else {
        glog.Infof("r.currentReadingFile=%v file=%v r.waiting=%v", r.currentReadingFile, file, r.waiting)
        glog.Infof("do not need to send kModify signal")
    }
    return nil
}

func (r *DirReader) OnFileCreated(file string) (err error) {
    r.add(file)
    if atomic.LoadInt32(&r.waiting) > 0 && r.files.Len() == 1 {
        // r.waiting : we will send a signal only when the goroutine is waiting
        // r.files.Len() == 0 : If we create more than 2 files in the same time, the waiting goroutine may be still waiting when we try to send the second signal
        glog.Infof("send kCreate signal")
        r.wakeup <- kCreate
    } else {
        glog.Infof("do not need to send kCreate signal")
    }
    return nil
}

func (r *DirReader) createReader() FileReader {
    if *readerType == "PTailReader" || *readerType == "GzipReader" {
        return NewTextFileTailReader(r)
    }

    return nil
}

func (r *DirReader) Stop() {
    r.Running = false
}

func (r *DirReader) Read() (err error) {
    glog.Infof("Starting to read files ...")
    for r.Running {
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

    return nil
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