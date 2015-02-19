package mikrotik_util

import (
	"fmt"
	"strconv"
	"strings"
)

type MtInterface struct {
	Index   int
	Name    string
	rb      *MikrotikSnmp
	snmpOid string
}

func (iff MtInterface) String() string {
	idx := fmt.Sprintf(".%d", iff.Index)
	return fmt.Sprintf("%3s %-23s %s", idx, iff.snmpOid, iff.Name)
}

// Returns the System + Interface Name
func (iff MtInterface) FullName() string {
	sysName, _ := iff.rb.GetSystemName()
	return fmt.Sprintf("%s-%s",
		sysName,
		iff.Name,
	)
}

// Takes a base OID (like bytes TX/RX) and returns that contatinated with
// the interfaces index
func (iff MtInterface) GetOidForInterface(oidBase string) string {
	oidIdx := strconv.Itoa(iff.Index)
	if !strings.HasSuffix(oidBase, ".") {
		oidBase += "."
	}
	return oidBase + oidIdx
}

// Returns the total bytes received on the interface
func (iff MtInterface) GetBytesRX() (ifBytes int64, err error) {
	pdu, err := iff.rb.SnmpGetPDU(iff.GetOidForInterface(SnmpOidIfBytesRX))
	if err != nil {
		return -1, err
	}
	return pdu.Value.(int64), nil
}

// Returns the total bytes transmitted on the interface
func (iff MtInterface) GetBytesTX() (ifBytes int64, err error) {
	pdu, err := iff.rb.SnmpGetPDU(iff.GetOidForInterface(SnmpOidIfBytesTX))
	if err != nil {
		return -1, err
	}
	return pdu.Value.(int64), nil
}

// Returns an array of the known interfaces
func (rb *MikrotikSnmp) GetInterfaces() (ifaces []MtInterface, err error) {
	// Check caches first
	if len(rb.cacheInterfaces) > 0 {
		// log.Println("Returning cached interfaces")
		return rb.cacheInterfaces, nil
	}

	snmpList, err := rb.SnmpGetPDUList(SnmpOidInterfaces)
	if err != nil {
		return nil, err
	}

	ifaces = make([]MtInterface, len(snmpList))
	for x, snmpIface := range snmpList {
		oidParts := strings.Split(snmpIface.Name, ".")
		// oidIndex, err := 0, error(nil)
		oidIndex, err := strconv.Atoi(oidParts[len(oidParts)-1])
		if err != nil {
			return nil, err
		}
		iface := MtInterface{
			Index:   oidIndex,
			Name:    snmpIface.Value.(string),
			rb:      rb,
			snmpOid: snmpIface.Name,
		}
		// log.Printf("Found interface: %s\n", iface)
		ifaces[x] = iface
	}
	rb.cacheInterfaces = ifaces
	return ifaces, nil
}

// Returns the given MtInterface, searching by name (exact match)
func (rb *MikrotikSnmp) FindInterfaceByName(ifName string) (iface MtInterface, err error) {
	if ifaces, err := rb.GetInterfaces(); err == nil {
		for _, iface := range ifaces {
			if iface.Name == ifName {
				return iface, nil
			}
		}
		return MtInterface{}, fmt.Errorf("Interface %s not found", ifName)
	}
	return MtInterface{}, err
}
