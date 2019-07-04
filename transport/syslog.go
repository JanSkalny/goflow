package transport

import (
	log "github.com/Sirupsen/logrus"
	syslog "log/syslog"
	flowmessage "github.com/cloudflare/goflow/pb"
	"fmt"
	"net"
	//proto "github.com/golang/protobuf/proto"
	//"os"
)

type SyslogState struct {
	writer *syslog.Writer

	addrs []string
}

func StartSyslogProducer(addrs []string) *SyslogState {
	writer, err := syslog.Dial("tcp", "127.0.0.1:514", syslog.LOG_NOTICE | syslog.LOG_KERN, "netflow")
	if err != nil {
	    log.Fatal("failed to dial remote syslog")
	}

	state := SyslogState{
		addrs: addrs,
		writer: writer,
	}

	return &state
}

func SyslogFormatProto(proto uint32) string {
	switch (proto) {
	case 1:
		return "ICMP"
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	}

	return fmt.Sprintf("%d", proto)
}
func SyslogFormatIP(addr []byte) string {
	ip := net.IP(addr)
	return ip.String()
}

func (s SyslogState) SendSyslogFlowMessage(f *flowmessage.FlowMessage) {
	log.Infof("new syslog flow message!")
	var msg string
	if (f.IPversion == 4) {
		msg = fmt.Sprintf("FW4-NETFLOW IN=%d OUT=%d MAC=%s:%s:%s SRC=%s DST=%s LEN=%d TOS=0x%02x PREC=0x%02x TTL=%d ID=0 PROTO=%s SPT=%d DPT=%d LEN=%d", f.SrcIf, f.DstIf, "", "", "", SyslogFormatIP(f.SrcIP), SyslogFormatIP(f.DstIP), f.Bytes, f.IPTos, 0, f.IPTTL, SyslogFormatProto(f.Proto), f.SrcPort, f.DstPort, f.Bytes)
	} else {
		msg = fmt.Sprintf("FW6-NETFLOW IN=%d OUT=%d MAC=%s:%s:%s SRC=%s DST=%s LEN=%d TOS=0x%02x PREC=0x%02x TTL=%d ID=0 PROTO=%s SPT=%d DPT=%d LEN=%d", f.SrcIf, f.DstIf, "", "", "", SyslogFormatIP(f.SrcIP), SyslogFormatIP(f.DstIP), f.Bytes, f.IPTos, 0, f.IPTTL, SyslogFormatProto(f.Proto), f.SrcPort, f.DstPort, f.Bytes)
	}
	log.Infof(msg)
	s.writer.Info("hello world")
	s.writer.Info(msg)
	/*
	b, _ := proto.Marshal(flowMessage)
	s.producer.Input() <- &sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.ByteEncoder(b),
	}
	*/
}
