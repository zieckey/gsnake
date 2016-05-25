package gsnake

import (
    "github.com/akrennmair/gopcap"
)

type PcapModule interface {
    OnPcapPacket(pkt *pcap.Packet)
}


//type DefaultTextModule struct {}
//type DefaultPcapModule struct {}
//
//func (m *DefaultTextModule) OnRecord(line []byte) {
//    glog.Infof("DefaultTextModule : Read a new line, len=%v <%s> ", len(line), string(line))
//}
//
//func (m *DefaultPcapModule) OnPcapPacket(pkt *pcap.Packet) {
//    glog.Infof("time: %d.%06d (%s) caplen: %d len: %d",
//        int64(pkt.Time.Second()), int64(pkt.Time.Nanosecond()),
//        time.Unix(int64(pkt.Time.Second()), 0).String(), int64(pkt.Caplen), int64(pkt.Len))
//}