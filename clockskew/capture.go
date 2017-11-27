package clockskew

import (
	//"log"
	"time"
	"encoding/binary"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func CapturePacket(){

	device := DeviceName

	handle, _ := pcap.OpenLive(
		device,
		int32(65535),
		false,
		-1 * time.Second,
	)

	defer handle.Close()

	bpFilter := BpConfig

	handle.SetBPFFilter(bpFilter)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil{
			continue
		}
		ip, _ := ipLayer.(*layers.IPv4)
		srcIP := ip.SrcIP.String()
		dstIP := ip.DstIP.String()

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort := tcp.SrcPort
		dstPort := tcp.DstPort

		for _, opt := range tcp.Options {
			if opt.OptionType.String() != "Timestamps" {
				continue
			}

			srcTS := binary.BigEndian.Uint32(opt.OptionData[:4])
			dstTS := binary.BigEndian.Uint32(opt.OptionData[4:8])

			cs := ClockSkew{
				Clock : time.Now().UnixNano(),
				SrcIP : srcIP,
				SrcPort : int(srcPort),
				SrcTS : int(srcTS),
				DstIP : dstIP,
				DstPort : int(dstPort),
				DstTS   : int(dstTS)}

			ClockSkewChannel <- cs
		}
	}
}

