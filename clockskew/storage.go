package clockskew

import (
	//"log"
	"fmt"
)

func Storage(){

	defer StorageFile.Close()
	defer close(ClockSkewChannel)

	for{
		cs := <- ClockSkewChannel

		item := fmt.Sprintf("%d %s %d %s %d %d %d\n",cs.Clock, cs.SrcIP, cs.SrcPort, cs.SrcTS, cs.DstIP, cs.DstPort, cs.DstTS)
		//log.Println(cs.Clock, cs.SrcIP, cs.SrcPort, cs.DstIP, cs.DstPort, cs.SrcTS, cs.DstTS)
		StorageFile.WriteString(item)
	}
}
