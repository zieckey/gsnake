package freader
import (
    "os"
    "github.com/golang/glog"
    "bufio"
    "bytes"
)

type PTailFileReader struct {
    path string
    pos int
    fp *os.File
    r* bufio.Reader // The reader of os.File fp
}

func NewPTailFileReader() *PTailFileReader {
    br := &PTailFileReader{
        path : "",
        pos:0,
        fp:nil,
    }

    return br
}

func (r *PTailFileReader) LoadFile(file string, pos int) (err error) {
    if r.fp != nil {
        r.fp.Close()
        r.fp = nil
    }

    r.path = file
    r.fp, err = os.OpenFile(file, os.O_RDONLY, 0644)
    if err != nil {
        glog.Errorf("OpenFile <%s> failed : %v\n", file, err.Error())
        return err
    }
    glog.Infof("OpenFile %v OK", file)

    if pos > 0 {
        r.fp.Seek(int64(pos), os.SEEK_SET)
    }

    if r.r == nil {
        r.r = bufio.NewReader(r.fp)
    } else {
        r.r.Reset(r.fp)
    }

    return nil
}

func (r *PTailFileReader) ReadLine() (line []byte, err error) {
    line, err = r.r.ReadBytes('\n')
    //glog.Infof("len(line)=%v %v", len(line), base64.StdEncoding.EncodeToString(line))
    line = bytes.TrimRight(line, "\r\n")
    //glog.Infof("len(line)=%v %v after trim", len(line), base64.StdEncoding.EncodeToString(line))
    return line, err
}

func (r *PTailFileReader) GetPos() int {
    return r.pos
}