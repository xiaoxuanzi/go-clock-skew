package clockskew

import (
	//"fmt"
	//"log"
	//"strconv"
	"os"
	//"net"
)

type ClockSkew struct {

	Clock       int64    `json:clock`
	SrcIP       string   `json:srcIp`
	SrcPort     int      `json:srcPort`
	SrcTS       int      `json:srcTimeStamp`
	DstIP       string   `json:dstIp`
	DstPort     int      `json:dstPort`
	DstTS       int      `json:dstTimeStamp`

}

var DeviceName   string
var BpConfig     string
var ClockSkewChannel chan ClockSkew
var StorageFile *os.File

