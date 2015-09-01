package gsnake
import (
    "os"
    "bufio"
    "bytes"
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

func (r *PTailFileReader) LoadFile(fp *os.File) (err error) {

    if r.r == nil {
        r.r = bufio.NewReader(fp)
    } else {
        r.r.Reset(fp)
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
