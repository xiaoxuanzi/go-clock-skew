package clockskew

import (
	"time"
	"fmt"
	"encoding/binary"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/gavv/monotime"
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
	//startClock := monotime.Now()

	for packet := range packetSource.Packets() {

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil{
			continue
		}
		ip, _ := ipLayer.(*layers.IPv4)
		srcIP := ip.SrcIP.String()

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort := tcp.SrcPort
		taddr := srcIP + ":" +  fmt.Sprintf("%d",srcPort)

		for _, opt := range tcp.Options {
			if opt.OptionType.String() != "Timestamps" {
				continue
			}

			srcTS := binary.BigEndian.Uint32(opt.OptionData[:4])

			//elapsed := monotime.Since(startClock) / 1000000000
			clock := monotime.Now()
			cs := ClockSkew{
				//Clock : time.Now().UnixNano(),
				Clock : int64(clock),
				Taddr : taddr,
				SrcTS : int(srcTS),
			}

			ClockSkewChannel <- cs
		}
	}
}

