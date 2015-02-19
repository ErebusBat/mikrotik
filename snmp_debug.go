package mikrotik_util

import (
	"log"

	"github.com/alouca/gosnmp"
)

func (rb *MikrotikSnmp) DumpOID(walk bool, oid string) {
	s, err := gosnmp.NewGoSNMP(rb.Host, rb.Community, gosnmp.Version2c, 5)
	if err != nil {
		log.Fatal(err)
	}

	if walk {
		resp, err := s.Walk(oid)
		if err != nil {
			log.Fatal(err)
		}
		dumpSnmpPDU(resp)
	} else {
		resp, err := s.Get(oid)
		if err != nil {
			log.Fatal(err)
		}
		dumpSnmpPacket(resp)
	}
}

func dumpSnmpPacket(packet *gosnmp.SnmpPacket) {
	for _, v := range packet.Variables {
		log.Printf("Response: n=%s : v=%v : t=%s \n",
			v.Name,
			v.Value,
			v.Type.String(),
		)
	}
}
func dumpSnmpPDU(slice []gosnmp.SnmpPDU) {
	// log.Printf("Response: %#v\n\n", slice)
	for idx, val := range slice {
		log.Printf("[%d] %#v\n", idx, val)
	}
}
