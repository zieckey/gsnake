package gsnake
import (
    "os"
    "github.com/golang/glog"
    "bufio"
    "bytes"
    "compress/gzip"
)

type GzipFileReader struct {
    r* bufio.Reader // The reader of os.File fp
    gr* gzip.Reader
}

func NewGzipFileReader() *GzipFileReader {
    br := &GzipFileReader {
        r : nil,
        gr : nil,
    }

    return br
}

func (r *GzipFileReader) LoadFile(filepath string, fp *os.File) (err error) {
    if r.r == nil {
        r.gr, err = gzip.NewReader(fp)
        if err != nil {
            glog.Errorf("Create gzip Reader failed : %v", err.Error())
            return err
        }
        r.r = bufio.NewReader(r.gr)
    } else {
        r.gr.Reset(fp)
        r.r.Reset(r.gr)
    }

    return nil
}

func (r *GzipFileReader) ReadLine() (line []byte, err error) {
    line, err = r.r.ReadBytes('\n')
    //glog.Infof("len(line)=%v %v", len(line), base64.StdEncoding.EncodeToString(line))
    line = bytes.TrimRight(line, "\r\n")
    //glog.Infof("len(line)=%v %v after trim", len(line), base64.StdEncoding.EncodeToString(line))
    return line, err
}
