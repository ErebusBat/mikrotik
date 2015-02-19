package mikrotik

import "github.com/ErebusBat/mikrotik/snmp"

func NewSnmp(host, community string) snmp.SnmpRouterboard {
	rb := &snmp.MikrotikSnmp{
		Host:      host,
		Community: community,
	}
	rb.Initialize()
	return rb
}
