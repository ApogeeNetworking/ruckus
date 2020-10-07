package ruckus

// Mapper ...
type Mapper struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SZAuthObj SmartZone Auth Object (*ServiceTicket)
type SZAuthObj struct {
	ControllerVersion string `json:"controllerVersion"`
	ServiceTicket     string `json:"serviceTicket"`
}

// RksCommonReq contains fields used in ALL Get Reqs
type RksCommonReq struct {
	TotalCount int  `json:"totalCount"`
	HasMore    bool `json:"hasMore"`
	FirstIndex int  `json:"firstIndex"`
}

// RksOptions common ruckus query options for data Retrieval
type RksOptions struct {
	// optional: the index of the 1st Entry to be retrieved.
	// Default 0
	Index string
	// optional: the max number of entries to be retrieved.
	// Default 100
	ListSize string
	// optional: The Domain ID.
	// Default: Current Domain ID
	DomainID string
}

// RksObject properties
type RksObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RksCommonRes the object returned when retrieving Rks Objects with Name/ID
type RksCommonRes struct {
	RksCommonReq
	List []RksObject `json:"list"`
}

// RksController Properties
type RksController struct {
	ID           string      `json:"id"`
	Model        string      `json:"model"`
	Description  string      `json:"description"`
	Hostname     string      `json:"hostName"`
	Mac          string      `json:"mac"`
	SerialNumber string      `json:"serialNumber"`
	ClusterRole  string      `json:"clusterRole"`
	ControlNatIP string      `json:"controlNatIp"`
	UptimeInSec  int         `json:"uptimeInSec"`
	Name         string      `json:"name"`
	Version      string      `json:"version"`
	ApVersion    string      `json:"apVersion"`
	ControlIP    string      `json:"controlIp"`
	ClusterIP    string      `json:"clusterIp"`
	MgmtIP       string      `json:"managementIp"`
	ControlIpv6  interface{} `json:"controlIpv6"`
	ClusterIpv6  interface{} `json:"clusterIpv6"`
	MgmtIpv6     interface{} `json:"managementIpv6"`
}

// RksSysSumRes ruckus controller result
type RksSysSumRes struct {
	RksCommonReq
	List []RksController `json:"list"`
}

// RksAP ruckus ap properties
type RksAP struct {
	MacAddr   string `json:"apMac"`
	ZoneID    string `json:"zoneId"`
	GroupID   string `json:"apGroupId"`
	GroupName string `json:"apGroupName"`
	ZoneName  string `json:"zoneName"`
	Serial    string `json:"serial"`
	ApName    string `json:"deviceName"`
	Model     string `json:"model"`
	Status    string `json:"status"`
	IPAddr    string `json:"ip"`
	ExtIPAddr string `json:"extIp"`
	Firmware  string `json:"firmwareVersion"`
}

// ApIntf ...
type ApIntf struct {
	MacAddr string `json:"apMac"`
	Speed   string `json:"phyLink"`
	Status  string `json:"logicLink"`
	Duplex  string
}

// ApLldp ...
type ApLldp struct {
	RemoteHostname string `json:"lldpSysName"`
	RemoteIntf     string `json:"lldpPortID"`
	RemoteIP       string `json:"lldpMgmtIP"`
}
