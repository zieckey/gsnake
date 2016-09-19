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
 go build  && ./text -file_path="." -file_pattern="*.log" -status=status.txt -reader_type="PTailReader" -stderrthreshold=0 -logtostderr=true
 */
func main() {
	flag.Parse()
	dispatcher, err := gsnake.NewDispatcher()
	if err != nil {
		panic(err.Error())
		return
	}
	dispatcher.Register(&MyTextModule{})
	dispatcher.Run()
}
