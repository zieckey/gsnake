package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/zieckey/gsnake"
)

type MyTextModule struct{}

func (m *MyTextModule) OnRecord(line []byte) {
	glog.Infof("DefaultTextModule : Read a new line, len=%v <%s> ", len(line), string(line))
}

/*
 go build  && ./freader.exe -file_path="d:\1\1" -file_pattern="ddd*"       -reader_type="PTailReader" -stderrthreshold=0 -logtostderr=true
 go build  && ./freader.exe -file_path="d:\1\gzip" -file_pattern="ddd*.gz" -reader_type="GzipReader"  -stderrthreshold=0 -logtostderr=true
 */
func main() {
	flag.Parse()
	dispatcher, err := gsnake.New()
	if err != nil {
		panic(err.Error())
		return
	}
	dispatcher.Register(&MyTextModule{})
	dispatcher.Run()
}
