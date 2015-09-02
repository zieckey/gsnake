package main

import (
	"flag"
	"fmt"
	"github.com/akrennmair/gopcap"
	"github.com/golang/glog"
	"github.com/zieckey/gsnake"
	"log"
	"time"
)

type MyTextModule struct{}

func (m *MyTextModule) OnRecord(line []byte) {
	glog.Infof("DefaultTextModule : Read a new line, len=%v <%s> ", len(line), string(line))
}

type MyPcapModule struct{}

func (m *MyPcapModule) OnPcapPacket(pkt *pcap.Packet) {
	glog.Infof("time: %d.%06d (%s) caplen: %d len: %d",
		int64(pkt.Time.Second()), int64(pkt.Time.Nanosecond()),
		time.Unix(int64(pkt.Time.Second()), 0).String(), int64(pkt.Caplen), int64(pkt.Len))

	for i := uint32(0); i < pkt.Caplen; i++ {
		if i%32 == 0 {
			fmt.Printf("\n")
		}
		if 32 <= pkt.Data[i] && pkt.Data[i] <= 126 {
			fmt.Printf("%c", pkt.Data[i])
		} else {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n\n")
}

/*
 go build && ./freader.exe -file_path="e:\1\1" -file_pattern="ddd*"       -reader_type="PTailReader" -stderrthreshold=0 -logtostderr=true
 go build && ./freader.exe -file_path="e:\1\gzip" -file_pattern="ddd*.gz" -reader_type="GzipReader"  -stderrthreshold=0 -logtostderr=true
 go build && ./freader.exe -file_path="e:\1\pcap" -file_pattern="*.pcap"  -reader_type="PcapReader"  -stderrthreshold=0 -logtostderr=true
 */
func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	dispatcher, err := gsnake.NewDispatcher()
	if err != nil {
		log.Panic(err.Error())
		return
	}
	dispatcher.RegisterPcapModule(&MyPcapModule{})
	dispatcher.RegisterTextModule(&MyTextModule{})
	dispatcher.Run()
}
