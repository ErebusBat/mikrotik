package main

import (
	"log"
	"time"

	"github.com/ErebusBat/mikrotik/snmp"
)

func main() {
	rb := snmp.Connect("10.0.1.250", "public")
	iface, err := rb.FindInterfaceByName("UPSTREAM")
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	sampleChannel := iface.MonitorBandwidth(time.Second * 5)
	for sample := range sampleChannel {
		log.Printf("%s tx/rx %s/%s",
			iface.FullName(),
			sample.TX().BitsString(),
			sample.RX().BitsString(),
		)
	}
}
