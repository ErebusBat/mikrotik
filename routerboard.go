package mikrotik_util

import "github.com/ErebusBat/mikrotik_util/snmp"

func NewSnmp(host, community string) snmp.SnmpRouterboard {
	rb := &snmp.MikrotikSnmp{
		Host:      host,
		Community: community,
	}
	rb.Initialize()
	return rb
}
