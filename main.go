package main

import (
	"log"
	"flag"
	"os"
	"github.com/google/gopacket/pcap"
	"github.com/go-clock-skew/clockskew"
)

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

func handleBPFFilter(device, bpFilter, storageFile  string) {

	clockskew.BpConfig   = bpFilter
	clockskew.DeviceName = device

	newfile, err := os.OpenFile(storageFile,
	os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Println("handleBPFFilter error: ", err)
		os.Exit(0)
	}

	clockskew.StorageFile = newfile
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func init(){

	clockskew.ClockSkewChannel = make(chan clockskew.ClockSkew, 1000)
	device  := flag.String("e", "eth0", "device name")
	help := flag.Bool("h", false, "help")

	bpFilter := flag.String("filter", "tcp", "bpFilter")

	storageFile := flag.String("f", "storage.csv", "storage file")

	flag.Parse()

	handleHelp(*help)
	handleDevice(*device)
	handleBPFFilter(*device, *bpFilter, *storageFile)

}

func main() {

	go clockskew.Storage()
	//go clockskew.SendClockSkew()
	clockskew.CapturePacket()
}

