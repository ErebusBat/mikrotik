package mikrotik

import "time"

type Routerboarder interface {
	Routerboard() Routerboard
}

type Routerboard interface {
	Routerboarder

	// Does initial setup in order to communicate to RB
	Initialize()

	// Removes any cached information
	FlushCaches()

	// Returns a time.Duration representing how long the RB has been running
	GetSystemUptime() (uptime time.Duration, err error)

	// Retuns the system description, i.e. RouterOS RB450G
	GetSystemDescription() (sysDesc string, err error)

	// Retuns the system name
	GetSystemName() (sysDesc string, err error)

	// Retuns a string comprised of the system name, description, and uptime
	GetSystemBanner() (banner string, err error)

	// Interface Methods
	GetInterfaces() (ifaces []RbInterface, err error)
	FindInterfaceByName(ifName string) (iface RbInterface, err error)
}
