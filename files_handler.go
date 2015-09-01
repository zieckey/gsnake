package freader
import (
    "path/filepath"
    "strconv"
    "log"
    "github.com/golang/glog"
)


type FilesHandler struct {
    dir string
    priorityLevel int

    readers map[string/*path dir*/]*DirReader
    paths []string
}

func NewFilesHandler(dir string) (h *FilesHandler, err error) {
    h = &FilesHandler{}
    h.dir = dir
    h.priorityLevel = *priorityLevel
    h.readers = make(map[string/*path*/]*DirReader)

    if h.priorityLevel <= 0 {
        h.readers[dir], err = NewPathReader(dir)
        if err != nil {
            return nil, err
        }
        h.paths = append(h.paths, dir)
    } else {
        for i := 0; i < h.priorityLevel; i++ {
            p := filepath.Join(dir, strconv.Itoa(i))
            h.readers[p], err = NewPathReader(p)
            if err != nil {
                return nil, err
            }
            h.paths = append(h.paths, p)
        }
    }

    return h, nil
}

func (h *FilesHandler) Run() {
    glog.Infof("Running ...")
    ff, err := LookupFiles(h.dir, *filePattern)
    if err != nil {
        log.Fatal("LoopupFiles <%s> with pathern <%s> failed : %v\n", dir, *filePattern, err.Error())
    }

    glog.Infof("existing files: %v", ff)
    for _, f := range ff {
        if !dispatcher.status.IsProcessed(f) {
            h.OnFileCreated(f)
        } else {
            glog.Infof("Skip processed file: %v", f)
        }
    }

    if h.priorityLevel <= 0 { // no priority
        r, _ := h.readers[h.dir]
        r.StartToRead()
    } else {
        for priority := 0; priority < len(h.readers); priority++ {
            r, _ := h.readers[h.paths[priority]]
            r.StartToRead() //TODO add priority and we can process next dir
        }
    }
}

func (h *FilesHandler) OnFileModified(file string) {
    dir := filepath.Dir(file)
    if r, ok := h.readers[dir]; ok {
        r.OnFileModified(file)
    } else {
        glog.Errorf("Append file failed, cannot found reader for this file <%v>", file)
    }
}

func (h *FilesHandler) OnFileCreated(file string) {
    dir := filepath.Dir(file)
    if r, ok := h.readers[dir]; ok {
        r.OnFileCreated(file)
    } else {
        glog.Errorf("Append file failed, cannot found reader for this file <%v>", file)
    }
}

func (h* FilesHandler) RecordPos() (err error) {
    //TODO
    return nil
}