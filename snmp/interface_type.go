package snmp

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/ErebusBat/mikrotik/core"
)

type MtInterface struct {
	_index   int
	_name    string
	_rb      SnmpRouterboard
	_snmpOid string
}

func (iff MtInterface) Index() int {
	return iff._index
}
func (iff MtInterface) Name() string {
	return iff._name
}
func (iff MtInterface) Routerboard() Routerboard {
	return iff._rb
}
func (iff MtInterface) SnmpOid() string {
	return iff._snmpOid
}

func (first MtInterface) Equals(second RbInterface) bool {
	return first.Name() == second.Name()
}
func (iff MtInterface) String() string {
	idx := fmt.Sprintf(".%d", iff.Index)
	return fmt.Sprintf("%3s %-23s %s", idx, iff._snmpOid, iff._name)
}

// Returns the System + Interface Name
func (iff MtInterface) FullName() string {
	sysName, _ := iff.Routerboard().GetSystemName()
	return fmt.Sprintf("%s-%s",
		sysName,
		iff._name,
	)
}

// Takes a base OID (like bytes TX/RX) and returns that contatinated with
// the interfaces index
func (iff MtInterface) GetOidForInterface(oidBase string) string {
	oidIdx := strconv.Itoa(iff._index)
	if !strings.HasSuffix(oidBase, ".") {
		oidBase += "."
	}
	return oidBase + oidIdx
}

// Returns the total bytes received on the interface
func (iff MtInterface) GetBytesRX() (ifBytes int64, err error) {
	pdu, err := iff._rb.SnmpGetPDU(iff.GetOidForInterface(SnmpOidIfBytesRX))
	if err != nil {
		return -1, err
	}
	return pdu.Value.(int64), nil
}

// Returns the total bytes transmitted on the interface
func (iff MtInterface) GetBytesTX() (ifBytes int64, err error) {
	pdu, err := iff._rb.SnmpGetPDU(iff.GetOidForInterface(SnmpOidIfBytesTX))
	if err != nil {
		return -1, err
	}
	return pdu.Value.(int64), nil
}
