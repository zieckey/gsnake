package gsnake

import (
	"github.com/akrennmair/gopcap"
	"github.com/golang/glog"
)

type PcapFileReader struct {
	h *pcap.Pcap
}

func NewPcapFileReader() *PcapFileReader {
	return &PcapFileReader{
		h: nil,
	}
}

func (r *PcapFileReader) ReadFile(file string, pos int) (err error) {
	r.h, err = pcap.Openoffline(file)
	if r.h == nil {
		glog.Errorf("Openoffline(%s) failed: %s", file, err)
		return err
	}
	defer r.h.Close()

	for pkt := r.h.Next(); pkt != nil; pkt = r.h.Next() {
		pkt.Decode()
		if dispatcher.pcapModule != nil {
			dispatcher.pcapModule.OnPcapPacket(pkt)
		}
	}
    return  nil
}
