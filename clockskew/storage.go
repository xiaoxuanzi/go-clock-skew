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

		item := fmt.Sprintf("%d %s %d\n",cs.Clock, cs.Taddr, cs.SrcTS)
		StorageFile.WriteString(item)
	}
}
