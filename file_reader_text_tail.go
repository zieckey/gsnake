package gsnake

import (
	"github.com/golang/glog"
	"io"
	"os"
)

type TextFileReader interface {
	LoadFile(filepath string, fp *os.File) error
	ReadLine() ([]byte, error)
}

type TextFileTailReader struct {
	path   string
	offset int
	fp     *os.File
	r      TextFileReader
	dr     *DirReader
}

func NewTextFileTailReader(dr *DirReader) *TextFileTailReader {
	r := &TextFileTailReader{
		path:   "",
		offset: 0,
		fp:     nil,
		dr:     dr,
	}

	if *reader_type == "PTailReader" {
		r.r = NewPTailFileReader()
	} else if *reader_type == "GzipReader" {
		r.r = NewGzipFileReader()
	} else {
		glog.Fatal("Error reader_type %v", *reader_type)
	}

	return r
}

func (r *TextFileTailReader) ReadFile(file string, offset int) (err error) {
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

	if offset > 0 {
		r.fp.Seek(int64(offset), os.SEEK_SET)
	}

	r.r.LoadFile(file, r.fp)

	r.readTextFile()

	return nil
}

func (r *TextFileTailReader) readTextFile() {
	var lastLine []byte
	for {
		line, err := r.r.ReadLine()
		//glog.Infof("ReadLine: lastLine=<%s> current-read=<%s> <%v>", string(lastLine), string(line), err)
		if len(lastLine) > 0 {
			line = append(lastLine, line...)
		}

		if err == io.EOF {
			if len(line) > 0 {
				lastLine = line
			}

			// there are still files which are ready to be processed
			if r.dr.GetPendingFileCount() > 0 {
				break
			}

			// no more files. we wait this file to be updated or wait new file created
			glog.Infof("no more files, we wait this file <%v> to be updated. Waiting ...", r.path)
			r.dr.Wait()
			continue
		} else if err != nil {
			glog.Errorf("Read data from <%s> failed : %v", r.path, err.Error())
			break
		} else {
			lastLine = []byte{}
		}

		dispatcher.textModule.OnRecord(line)
	}
}
