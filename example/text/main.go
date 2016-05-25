package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/zieckey/gsnake"
	"log"
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
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	dispatcher, err := gsnake.NewDispatcher()
	if err != nil {
		log.Panic(err.Error())
		return
	}
	dispatcher.RegisterTextModule(&MyTextModule{})
	dispatcher.Run()
}
