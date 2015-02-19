package snmp

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/ErebusBat/mikrotik"
)

// Returns an array of the known interfaces
func (rb *MikrotikSnmp) GetInterfaces() (ifaces []RbInterface, err error) {
	// Check caches first
	if len(rb.cacheInterfaces) > 0 {
		// log.Println("Returning cached interfaces")
		return rb.cacheInterfaces, nil
	}

	snmpList, err := rb.SnmpGetPDUList(SnmpOidInterfaces)
	if err != nil {
		return nil, err
	}

	ifaces = make([]RbInterface, len(snmpList))
	for x, snmpIface := range snmpList {
		oidParts := strings.Split(snmpIface.Name, ".")
		// oidIndex, err := 0, error(nil)
		oidIndex, err := strconv.Atoi(oidParts[len(oidParts)-1])
		if err != nil {
			return nil, err
		}
		iface := MtInterface{
			_index:   oidIndex,
			_name:    snmpIface.Value.(string),
			_rb:      rb,
			_snmpOid: snmpIface.Name,
		}
		// log.Printf("Found interface: %s\n", iface)
		ifaces[x] = iface
	}
	rb.cacheInterfaces = ifaces
	return ifaces, nil
}

// Returns the given RbInterface, searching by name (exact match)
func (rb *MikrotikSnmp) FindInterfaceByName(ifName string) (iface RbInterface, err error) {
	if ifaces, err := rb.GetInterfaces(); err == nil {
		for _, iface := range ifaces {
			if iface.Name() == ifName {
				return iface, nil
			}
		}
		return nil, fmt.Errorf("Interface %s not found", ifName)
	}
	return nil, err
}
