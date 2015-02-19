package mikrotik_util

import "github.com/alouca/gosnmp"

// http://wiki.mikrotik.com/wiki/Munin_Monitoring#mikrotikifrate
type MikrotikSnmp struct {
	Host      string
	Community string

	// Caches
	cacheInterfaces []MtInterface
	cacheStrings    map[string]string
}

const (
	SnmpOidDescription = ".1.3.6.1.2.1.1.1.0"
	SnmpOidUptime      = ".1.3.6.1.2.1.1.3.0"
	SnmpOidName        = ".1.3.6.1.2.1.1.5.0"
	SnmpOidInterfaces  = ".1.3.6.1.2.1.2.2.1.2"
	SnmpOidIfBytesRX   = ".1.3.6.1.2.1.31.1.1.1.6."
	SnmpOidIfBytesTX   = ".1.3.6.1.2.1.31.1.1.1.10."
)

var (
	zeroValPdu = gosnmp.SnmpPDU{}
)

// Setups up object to be ready to works (caches, etc)
func (rb *MikrotikSnmp) Initialize() {
	rb.FlushCaches()
}

// Removes any cached information
func (rb *MikrotikSnmp) FlushCaches() {
	rb.cacheInterfaces = make([]MtInterface, 0, 5)
	rb.cacheStrings = make(map[string]string, 0)
}

// Looks up the given OID value, casts it to a string and caches it.  Returns
// cached value early if it exists
func (rb *MikrotikSnmp) GetOidStringValCached(oid string) (val string, err error) {
	if cacheVal, foundInCache := rb.cacheStrings[oid]; foundInCache {
		return cacheVal, nil
	}

	pdu, err := rb.SnmpGetPDU(oid)
	if err != nil {
		return "", err
	}

	// Cast val, cache it and return it
	val = pdu.Value.(string)
	rb.cacheStrings[oid] = val
	return val, nil
}

func (rb *MikrotikSnmp) SnmpGetPDU(oid string) (gosnmp.SnmpPDU, error) {
	s, err := gosnmp.NewGoSNMP(rb.Host, rb.Community, gosnmp.Version2c, 5)
	if err != nil {
		return zeroValPdu, err
	}
	resp, err := s.Get(oid)
	if err != nil {
		return zeroValPdu, err
	}

	if len(resp.Variables) > 0 {
		return resp.Variables[0], nil
	}
	return zeroValPdu, nil
}

func (rb *MikrotikSnmp) SnmpGetPDUList(oid string) ([]gosnmp.SnmpPDU, error) {
	s, err := gosnmp.NewGoSNMP(rb.Host, rb.Community, gosnmp.Version2c, 5)
	if err != nil {
		return nil, err
	}
	resp, err := s.Walk(oid)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
