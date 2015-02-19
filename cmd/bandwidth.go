package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	. "github.com/ErebusBat/mikrotik_util"
)

// Custom type/consts to make our action routing code easier to read
type cliAction int

const (
	unknown cliAction = iota
	monitorBandwidth
	printInterfaces
)

// Configuration Struct
type AppConfig struct {
	RouterBoard    MikrotikSnmp
	SampleInterval time.Duration
	InterfaceName  string

	action         cliAction
	dumpInterfaces bool
}

////////////////////////////////////////////////////////////////////////////////
// ENTRY POINT
////////////////////////////////////////////////////////////////////////////////

func main() {
	cfg := parseConfig()

	switch cfg.action {
	case monitorBandwidth:
		iface := cfg.mustFindInterface()
		actionMonitorBandwidth(iface, cfg.SampleInterval)
	case printInterfaces:
		actionPrintKnownInterfaces(cfg.RouterBoard)
	default:
		log.Fatalf("Unknown action %#v?!?! ", cfg.action)
	}
}

////////////////////////////////////////////////////////////////////////////////
// ACTIONS
////////////////////////////////////////////////////////////////////////////////

// Prints all known interfaces to the console
func actionPrintKnownInterfaces(rb MikrotikSnmp) {
	log.Println("Interfaces")
	ifaces, err := rb.GetInterfaces()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	header := fmt.Sprintf("%5s %-23s %s",
		"idx",
		"       Full Oid",
		"Name")
	log.Println(header)
	log.Println("--------------------------------------------------------------------------------")
	for _, iface := range ifaces {
		log.Println(" ", iface)
	}
}

// action: monitors the bandwidth of given interface and prints to console
func actionMonitorBandwidth(iface MtInterface, sampleInterval time.Duration) {
	log.Printf("Sampling bandwidth every %s\n", sampleInterval)
	bandwidthSample := iface.MonitorBandwidth(sampleInterval)
	for sample := range bandwidthSample {
		log.Printf("%s tx/rx %s/%s",
			iface.FullName(),
			sample.TX().BitsString(),
			sample.RX().BitsString(),
		)
	}
}

////////////////////////////////////////////////////////////////////////////////
// HELPERS
////////////////////////////////////////////////////////////////////////////////

// Reads the configuration, inits objects, and returns the AppConfig
func parseConfig() *AppConfig {
	cfg := new(AppConfig)
	cfg.action = monitorBandwidth
	flag.StringVar(&cfg.RouterBoard.Community, "c", "public", "SNMP Community Name")
	flag.StringVar(&cfg.RouterBoard.Host, "h", "127.0.0.7", "Mikrotik IP")

	flag.DurationVar(&cfg.SampleInterval, "s", time.Second, "Sample Interval")
	flag.StringVar(&cfg.InterfaceName, "i", "ether1", "Mikrotik Interface Name")

	// Non operational flags
	flag.BoolVar(&cfg.dumpInterfaces, "list", false, "Lists all known interfaces and exits")
	flag.Parse()
	cfg.RouterBoard.Initialize()

	// Print RB banner (so it is on all output)
	banner, err := cfg.RouterBoard.GetSystemBanner()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	log.Println("Connected to", banner)

	// Now check non-operational status

	if cfg.dumpInterfaces {
		cfg.action = printInterfaces
	}

	return cfg
}

// helper: returns system name or logs fatal error
func (cfg *AppConfig) mustGetSystemName() string {
	sysName, err := cfg.RouterBoard.GetSystemName()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	return sysName
}

// helper: returns interface or logs fatal error
func (cfg *AppConfig) mustFindInterface() MtInterface {
	// Call sysName first so we exit on that
	sysName := cfg.mustGetSystemName()
	iface, err := cfg.RouterBoard.FindInterfaceByName(cfg.InterfaceName)
	if err != nil {
		log.Fatalf("ERROR: %v (maybe try --list)", err)
	}
	log.Printf("%s Found %s Interface at %d\n",
		sysName,
		iface.Name,
		iface.Index,
	)
	return iface
}
