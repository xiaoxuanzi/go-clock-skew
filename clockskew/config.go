package clockskew

import (
	"os"
)

type ClockSkew struct {

	Clock       int64    `json:clock`
	Taddr       string   `json:taddr`
	//SrcPort     int      `json:srcPort`
	SrcTS       int      `json:srcTimeStamp`
	//DstIP       string   `json:dstIp`
	//DstPort     int      `json:dstPort`
	//DstTS       int      `json:dstTimeStamp`

}

var DeviceName   string
var BpConfig     string
var ClockSkewChannel chan ClockSkew
var StorageFile *os.File

