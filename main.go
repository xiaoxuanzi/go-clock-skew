package main

import (
	//"fmt"
	"log"
	"flag"
	"os"
	"strconv"
	"time"
	//"net"
	"encoding/binary"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type ClockSkew struct {

	TimeStamp   int64    `json:timeStamp`
	SrcIP       string   `json:srcIp`
	SrcPort     int      `json:srcPort`
	SrcTS       int      `json:srcTimeStamp`
	DstIP       string   `json:dstIp`
	DstPort     int      `json:dstPort`
	DstTS       int      `json:dstTimeStamp`

}

type BPFFilterConfig struct {

	DeviceName  string   `json:deviceName`
	Protocol    string   `json:protocol`
	SrcIP       string   `json:srcIp`
	SrcPort     int      `json:srcPort`
	DstIP       string   `json:dstIp`
	DstPort     int      `json:dstPort`

}

var bpConfig *BPFFilterConfig

func handleDevice(device string) {

	flag := false

	devices, err := pcap.FindAllDevs()
		if err != nil {
		log.Fatal(err)
	}

	for _, dev := range devices{
		if device == dev.Name{
			flag = true
		}
	}

	if flag == false{
		log.Fatalln("[ERROR] not found device ", device)
	}
}

func handleBPFFilter(device, protocol, srcIP, dstIP string, srcPort, dstPort int) {

	var cfg BPFFilterConfig

	cfg.DeviceName = device
	cfg.Protocol   = protocol
	cfg.SrcIP      = srcIP
	cfg.SrcPort    = srcPort
	cfg.DstIP      = dstIP
	cfg.DstPort    = dstPort

	bpConfig = &cfg

}
func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func init(){

	device  := flag.String("e", "eth0", "device name")
	help := flag.Bool("h", false, "help")

	protocol := flag.String("proto", "tcp", "protocol")
	srcIP := flag.String("h1", "None", "src host")
	srcPort := flag.Int("d1", 0, "src port")

	dstIP := flag.String("h2", "None", "dst host")
	dstPort := flag.Int("d2", 0, "dst host")

	flag.Parse()

	handleHelp(*help)
	handleDevice(*device)
	handleBPFFilter(*device, *protocol, *srcIP, *dstIP, *srcPort, *dstPort)

}

func main() {

	bpFilter := getBPFFilter()

	log.Println("bpConfig: ", bpConfig, " bpFilter: ", bpFilter)
	CapturePacket()
}

func CapturePacket(){

	device := bpConfig.DeviceName

	handle, _ := pcap.OpenLive(
		device,
		int32(65535),
		false,
		-1 * time.Second,
	)

	defer handle.Close()

	bpFilter := getBPFFilter()

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
				TimeStamp : time.Now().UnixNano(),
				SrcIP : srcIP,
				SrcPort : int(srcPort),
				SrcTS : int(srcTS),
				DstIP : dstIP,
				DstPort : int(dstPort),
				DstTS   : int(dstTS)}

			log.Println(cs.TimeStamp, cs.SrcIP, cs.SrcPort, cs.DstIP, cs.DstPort, cs.SrcTS, cs.DstTS)
		}
	}
}

func getBPFFilter() string{

	bpFilter := bpConfig.Protocol

	if bpConfig.SrcIP != "None"{
		bpFilter += (" and src host " + bpConfig.SrcIP)
	}

	if bpConfig.SrcPort != 0{
		bpFilter += (" and src port " + strconv.Itoa(bpConfig.SrcPort))
	}

	if bpConfig.DstIP != "None"{
		bpFilter += (" and dst host " + bpConfig.DstIP)
	}

	if bpConfig.DstPort != 0{
		bpFilter += (" and dst port " + strconv.Itoa(bpConfig.DstPort))
	}

	return bpFilter
}

