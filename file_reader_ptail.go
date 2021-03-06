package gsnake
import (
    "os"
    "bufio"
    "bytes"
    "github.com/golang/glog"
)

type PTailFileReader struct {
    r* bufio.Reader // The reader of os.File fp
}

func NewPTailFileReader() *PTailFileReader {
    br := &PTailFileReader{
        r:nil,
    }

    return br
}

func (r *PTailFileReader) LoadFile(filepath string, fp *os.File) (err error) {
    if r.r == nil {
        glog.Infof("LoadFile : it is the first time to come here, we create a new reader: bufio.NewReader(fp)")
        r.r = bufio.NewReader(fp)
    } else {
        glog.Infof("Reset reader")
        r.r.Reset(fp)
    }

    return nil
}

func (r *PTailFileReader) ReadLine() (line []byte, err error) {
    line, err = r.r.ReadBytes('\n')
    glog.Infof("len(line)=%v %v", len(line), string(line))
    line = bytes.TrimRight(line, "\r\n")
    glog.Infof("len(line)=%v %v after trim", len(line), string(line))
    return line, err
}
