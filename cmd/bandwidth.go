package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ErebusBat/mikrotik/core"
	"github.com/ErebusBat/mikrotik/snmp"
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
	Routerboard    core.Routerboard
	SampleInterval time.Duration
	InterfaceName  string

	action         cliAction
	dumpInterfaces bool
}

////////////////////////////////////////////////////////////////////////////////
// ENTRY POINT
////////////////////////////////////////////////////////////////////////////////

func main() {
	app := parseConfig()

	switch app.action {
	case monitorBandwidth:
		app.actionMonitorBandwidth()
	case printInterfaces:
		app.actionPrintKnownInterfaces()
	default:
		log.Fatalf("Unknown action %#v?!?! ", app.action)
	}
}

////////////////////////////////////////////////////////////////////////////////
// ACTIONS
////////////////////////////////////////////////////////////////////////////////

// Prints all known interfaces to the console
func (cfg *AppConfig) actionPrintKnownInterfaces() {
	log.Println("Interfaces")
	ifaces, err := cfg.Routerboard.GetInterfaces()
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
func (cfg *AppConfig) actionMonitorBandwidth() {
	iface := cfg.mustFindInterface()
	var sampleInterval time.Duration = cfg.SampleInterval
	log.Printf("Sampling bandwidth every %s\n", sampleInterval)
	bandwidthSample := iface.MonitorBandwidth(sampleInterval)
	for sample := range bandwidthSample {
		if &sample == nil {
			return
		}
		// log.Printf("%#v\n\n", sample)
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
	var (
		host, community string
	)
	cfg := new(AppConfig)
	cfg.action = monitorBandwidth
	flag.StringVar(&host, "h", "127.0.0.7", "Mikrotik IP")
	flag.StringVar(&community, "c", "public", "SNMP Community Name")

	flag.DurationVar(&cfg.SampleInterval, "s", time.Second, "Sample Interval")
	flag.StringVar(&cfg.InterfaceName, "i", "ether1", "Mikrotik Interface Name")

	// Non operational flags
	flag.BoolVar(&cfg.dumpInterfaces, "list", false, "Lists all known interfaces and exits")
	flag.Parse()
	cfg.Routerboard = snmp.Connect(host, community)

	// Print RB banner (so it is on all output)
	banner, err := cfg.Routerboard.GetSystemBanner()
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
	sysName, err := cfg.Routerboard.GetSystemName()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	return sysName
}

// helper: returns interface or logs fatal error
func (cfg *AppConfig) mustFindInterface() core.RbInterface {
	// Call sysName first so we exit on that
	sysName := cfg.mustGetSystemName()
	iface, err := cfg.Routerboard.FindInterfaceByName(cfg.InterfaceName)
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
