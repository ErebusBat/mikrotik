package snmp

import (
	"fmt"
	"time"

	"github.com/alouca/gosnmp"

	. "github.com/ErebusBat/mikrotik"
)

type SnmpRouterboard interface {
	Routerboard

	GetOidStringValCached(oid string) (val string, err error)
	SnmpGetPDU(oid string) (gosnmp.SnmpPDU, error)
	SnmpGetPDUList(oid string) ([]gosnmp.SnmpPDU, error)
}

// http://wiki.mikrotik.com/wiki/Munin_Monitoring#mikrotikifrate
type MikrotikSnmp struct {
	SnmpRouterboard

	Host      string
	Community string

	// Caches
	cacheInterfaces []RbInterface
	cacheStrings    map[string]string
}

func (rb *MikrotikSnmp) Routerboard() Routerboard {
	return rb
}

// Setups up object to be ready to works (caches, etc)
func (rb *MikrotikSnmp) Initialize() {
	rb.FlushCaches()
}

// Removes any cached information
func (rb *MikrotikSnmp) FlushCaches() {
	rb.cacheInterfaces = make([]RbInterface, 0, 5)
	rb.cacheStrings = make(map[string]string, 0)
}

// Returns a time.Duration representing how long the RB has been running
func (rb *MikrotikSnmp) GetSystemUptime() (uptime time.Duration, err error) {
	pdu, err := rb.SnmpGetPDU(SnmpOidUptime)
	if err != nil {
		return time.Duration(0), err
	}
	upSeconds := pdu.Value.(int) / 100
	upSecsParseString := fmt.Sprintf("%ds", upSeconds)
	uptime, err = time.ParseDuration(upSecsParseString)
	if err != nil {
		return time.Duration(0), err
	}
	return uptime, nil
}

// Retuns the system description, i.e. RouterOS RB450G
func (rb *MikrotikSnmp) GetSystemDescription() (sysDesc string, err error) {
	return rb.GetOidStringValCached(SnmpOidDescription)
}

// Retuns the system name
func (rb *MikrotikSnmp) GetSystemName() (sysDesc string, err error) {
	return rb.GetOidStringValCached(SnmpOidName)
}

// Retuns a string comprised of the system name, description, and uptime
func (rb *MikrotikSnmp) GetSystemBanner() (banner string, err error) {

	sysName, err := rb.GetSystemName()
	if err != nil {
		return "", err
	}

	sysDesc, err := rb.GetSystemDescription()
	if err != nil {
		return "", err
	}

	sysUptime, err := rb.GetSystemUptime()
	if err != nil {
		return "", err
	}
	sysUpDays := 0
	if sysUptime.Hours() >= 24 {
		sysUpDays = int(sysUptime.Hours()/float64(24) + 0.5)
	}

	return fmt.Sprintf("%s %s (Uptime ~%d days: %s)", sysName, sysDesc, sysUpDays, sysUptime), nil
}
