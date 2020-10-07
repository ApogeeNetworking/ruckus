package ruckus

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const loginErr string = "you must first login to perform this action"

// Client struct is used to handle the Connection with the SmartZone Controller
type Client struct {
	BaseURL  string
	host     string
	username string
	password string

	http          *http.Client
	serviceTicket string
}

// New creates a Reference to a Client
func New(host, user, pass string, ignoreSSL bool) *Client {
	return &Client{
		BaseURL:  fmt.Sprintf("https://%s:8443/wsg/api/public/v8_1", host),
		host:     host,
		username: user,
		password: pass,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: ignoreSSL,
				},
			},
			Timeout: 8 * time.Second,
		},
	}
}

// Login est a session with the Ruckus SZ Controller
func (c *Client) Login() error {
	// Create our Auth JSON Object|Convert to Reader for POST REQ
	authObj := struct {
		User string `json:"username"`
		Pass string `json:"password"`
	}{User: c.username, Pass: c.password}
	jdata, _ := json.Marshal(&authObj)
	credentials := strings.NewReader(string(jdata))
	req, err := http.NewRequest("POST", c.BaseURL+"/serviceTicket", credentials)
	if err != nil {
		return fmt.Errorf("failed to create a new request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}
	defer res.Body.Close()
	// Auth RESP returns an JSON Object with serviceTicket Field
	auth := struct {
		Ticket string `json:"serviceTicket"`
	}{}
	json.NewDecoder(res.Body).Decode(&auth)
	c.serviceTicket = auth.Ticket
	return nil
}

// Logout removes a sessions with the Ruckus SZ Controller
func (c *Client) Logout() error {
	req, err := http.NewRequest("DELETE", c.BaseURL+"/serviceTicket", nil)
	if err != nil {
		return fmt.Errorf("failed to create a new request: %v", err)
	}
	q := req.URL.Query()
	q.Add("serviceTicket", c.serviceTicket)
	req.URL.RawQuery = q.Encode()

	_, err = c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to logout: %v", err)
	}
	c.serviceTicket = ""
	return nil
}

// GetZones retrieves a Paginated List of Zones
func (c *Client) GetZones(o RksOptions) (RksCommonRes, error) {
	if c.serviceTicket == "" {
		return RksCommonRes{}, fmt.Errorf(loginErr)
	}
	req, err := c.genGetReq("/rkszones")
	if err != nil {
		return RksCommonRes{}, err
	}
	// Update the Request
	c.addQS(req, o)

	res, err := c.http.Do(req)
	if err != nil {
		return RksCommonRes{}, fmt.Errorf("failed to get resp: %v", err)
	}
	defer res.Body.Close()
	var zones RksCommonRes
	json.NewDecoder(res.Body).Decode(&zones)

	return zones, nil
}

// GetZone retrieve Zone Configuration from Rks Controller
func (c *Client) GetZone(id string) (RksZone, error) {
	if c.serviceTicket == "" {
		return RksZone{}, fmt.Errorf(loginErr)
	}
	req, err := c.genGetReq(fmt.Sprintf("/rkszones/%s", id))
	if err != nil {
		return RksZone{}, err
	}
	c.addQS(req, RksOptions{})
	res, err := c.http.Do(req)
	if err != nil {
		return RksZone{}, fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()
	var zone RksZone
	json.NewDecoder(res.Body).Decode(&zone)
	return zone, nil
}

// GetSysSum retrieves system summary information from the Ruckus Controller
func (c *Client) GetSysSum(o RksOptions) (RksSysSumRes, error) {
	if c.serviceTicket == "" {
		return RksSysSumRes{}, fmt.Errorf(loginErr)
	}
	req, err := c.genGetReq("/controller")
	if err != nil {
		return RksSysSumRes{}, err
	}
	c.addQS(req, o)

	res, err := c.http.Do(req)
	if err != nil {
		return RksSysSumRes{}, fmt.Errorf("failed to get resp: %v", err)
	}
	defer res.Body.Close()
	var sysSum RksSysSumRes
	json.NewDecoder(res.Body).Decode(&sysSum)
	return sysSum, nil
}

// RksAp an Access Point in a SZ Controller
type RksAp struct {
	DeviceName                     string      `json:"deviceName"`
	Description                    string      `json:"description"`
	Status                         string      `json:"status"`
	Alerts                         int         `json:"alerts"`
	IP                             string      `json:"ip"`
	Ipv6Address                    interface{} `json:"ipv6Address"`
	TxRx                           interface{} `json:"txRx"`
	Noise24G                       int         `json:"noise24G"`
	Noise5G                        int         `json:"noise5G"`
	Airtime24G                     int         `json:"airtime24G"`
	Airtime5G                      int         `json:"airtime5G"`
	Latency24G                     int         `json:"latency24G"`
	Latency50G                     int         `json:"latency50G"`
	Capacity                       int         `json:"capacity"`
	Capacity24G                    int         `json:"capacity24G"`
	Capacity50G                    int         `json:"capacity50G"`
	ConnectionFailure              int         `json:"connectionFailure"`
	Model                          string      `json:"model"`
	ApMac                          string      `json:"apMac"`
	Channel24G                     string      `json:"channel24G"`
	Channel5G                      string      `json:"channel5G"`
	Channel24GValue                int         `json:"channel24gValue"`
	Channel50GValue                int         `json:"channel50gValue"`
	MeshRole                       string      `json:"meshRole"`
	MeshMode                       string      `json:"meshMode"`
	ZoneName                       string      `json:"zoneName"`
	ZoneAffinityProfileName        string      `json:"zoneAffinityProfileName"`
	ApGroupName                    string      `json:"apGroupName"`
	ExtIP                          string      `json:"extIp"`
	ExtPort                        string      `json:"extPort"`
	FirmwareVersion                string      `json:"firmwareVersion"`
	Serial                         string      `json:"serial"`
	Retry24G                       int         `json:"retry24G"`
	Retry5G                        int         `json:"retry5G"`
	ConfigurationStatus            string      `json:"configurationStatus"`
	LastSeen                       int64       `json:"lastSeen"`
	NumClients                     int         `json:"numClients"`
	NumClients24G                  int         `json:"numClients24G"`
	NumClients5G                   int         `json:"numClients5G"`
	Tx                             interface{} `json:"tx"`
	Rx                             interface{} `json:"rx"`
	Location                       string      `json:"location"`
	WlanGroup24ID                  interface{} `json:"wlanGroup24Id"`
	WlanGroup50ID                  interface{} `json:"wlanGroup50Id"`
	WlanGroup24Name                interface{} `json:"wlanGroup24Name"`
	WlanGroup50Name                interface{} `json:"wlanGroup50Name"`
	EnabledBonjourGateway          bool        `json:"enabledBonjourGateway"`
	ControlBladeName               string      `json:"controlBladeName"`
	LbsStatus                      string      `json:"lbsStatus"`
	AdministrativeState            string      `json:"administrativeState"`
	RegistrationState              string      `json:"registrationState"`
	ProvisionMethod                string      `json:"provisionMethod"`
	ProvisionStage                 string      `json:"provisionStage"`
	RegistrationTime               int64       `json:"registrationTime"`
	ManagementVlan                 interface{} `json:"managementVlan"`
	ConfigOverride                 bool        `json:"configOverride"`
	IndoorMapID                    interface{} `json:"indoorMapId"`
	ApGroupID                      string      `json:"apGroupId"`
	IndoorMapXy                    interface{} `json:"indoorMapXy"`
	IndoorMapName                  interface{} `json:"indoorMapName"`
	IndoorMapLocation              interface{} `json:"indoorMapLocation"`
	DeviceGps                      string      `json:"deviceGps"`
	ConnectionStatus               string      `json:"connectionStatus"`
	ZoneID                         string      `json:"zoneId"`
	ZoneFirmwareVersion            string      `json:"zoneFirmwareVersion"`
	DomainID                       string      `json:"domainId"`
	DomainName                     interface{} `json:"domainName"`
	DpIP                           string      `json:"dpIp"`
	ControlBladeID                 string      `json:"controlBladeId"`
	IsCriticalAp                   bool        `json:"isCriticalAp"`
	CrashDump                      interface{} `json:"crashDump"`
	CableModemSupported            bool        `json:"cableModemSupported"`
	CableModemResetSupported       bool        `json:"cableModemResetSupported"`
	SwapInMac                      interface{} `json:"swapInMac"`
	SwapOutMac                     interface{} `json:"swapOutMac"`
	IsOverallHealthStatusFlagged   bool        `json:"isOverallHealthStatusFlagged"`
	IsLatency24GFlagged            bool        `json:"isLatency24GFlagged"`
	IsCapacity24GFlagged           bool        `json:"isCapacity24GFlagged"`
	IsConnectionFailure24GFlagged  bool        `json:"isConnectionFailure24GFlagged"`
	IsLatency50GFlagged            bool        `json:"isLatency50GFlagged"`
	IsCapacity50GFlagged           bool        `json:"isCapacity50GFlagged"`
	IsConnectionFailure50GFlagged  bool        `json:"isConnectionFailure50GFlagged"`
	IsConnectionTotalCountFlagged  interface{} `json:"isConnectionTotalCountFlagged"`
	IsConnectionFailureFlagged     bool        `json:"isConnectionFailureFlagged"`
	IsAirtimeUtilization24GFlagged bool        `json:"isAirtimeUtilization24GFlagged"`
	IsAirtimeUtilization50GFlagged bool        `json:"isAirtimeUtilization50GFlagged"`
	Uptime                         interface{} `json:"uptime"`
	Eirp24G                        int         `json:"eirp24G"`
	Eirp50G                        int         `json:"eirp50G"`
	SupportFips                    interface{} `json:"supportFips"`
	IpsecSessionTime               interface{} `json:"ipsecSessionTime"`
	IpsecTxPkts                    interface{} `json:"ipsecTxPkts"`
	IpsecRxPkts                    interface{} `json:"ipsecRxPkts"`
	IpsecTxBytes                   interface{} `json:"ipsecTxBytes"`
	IpsecRxBytes                   interface{} `json:"ipsecRxBytes"`
	IpsecTxDropPkts                interface{} `json:"ipsecTxDropPkts"`
	IpsecRxDropPkts                interface{} `json:"ipsecRxDropPkts"`
	IpsecTxIdleTime                interface{} `json:"ipsecTxIdleTime"`
	IpsecRxIdleTime                interface{} `json:"ipsecRxIdleTime"`
	IPType                         string      `json:"ipType"`
	PacketCaptureState             string      `json:"packetCaptureState"`
	CellularWanInterface           interface{} `json:"cellularWanInterface"`
	CellularConnectionStatus       interface{} `json:"cellularConnectionStatus"`
	CellularSignalStrength         interface{} `json:"cellularSignalStrength"`
	CellularIMSISIM0               interface{} `json:"cellularIMSISIM0"`
	CellularIMSISIM1               interface{} `json:"cellularIMSISIM1"`
	CellularICCIDSIM0              interface{} `json:"cellularICCIDSIM0"`
	CellularICCIDSIM1              interface{} `json:"cellularICCIDSIM1"`
	CellularIsSIM0Present          interface{} `json:"cellularIsSIM0Present"`
	CellularIsSIM1Present          interface{} `json:"cellularIsSIM1Present"`
	CellularTxBytesSIM0            interface{} `json:"cellularTxBytesSIM0"`
	CellularTxBytesSIM1            interface{} `json:"cellularTxBytesSIM1"`
	CellularRxBytesSIM0            interface{} `json:"cellularRxBytesSIM0"`
	CellularRxBytesSIM1            interface{} `json:"cellularRxBytesSIM1"`
	CellularActiveSim              interface{} `json:"cellularActiveSim"`
	CellularIPaddress              interface{} `json:"cellularIPaddress"`
	CellularSubnetMask             interface{} `json:"cellularSubnetMask"`
	CellularDefaultGateway         interface{} `json:"cellularDefaultGateway"`
	CellularOperator               interface{} `json:"cellularOperator"`
	Cellular3G4GChannel            interface{} `json:"cellular3G4GChannel"`
	CellularCountry                interface{} `json:"cellularCountry"`
	CellularRadioUptime            interface{} `json:"cellularRadioUptime"`
	CellularGpsHistory             interface{} `json:"cellularGpsHistory"`
	FipsEnabled                    interface{} `json:"fipsEnabled"`
	MedianTxRadioMCSRate24G        int         `json:"medianTxRadioMCSRate24G"`
	MedianTxRadioMCSRate50G        int         `json:"medianTxRadioMCSRate50G"`
	MedianRxRadioMCSRate24G        int         `json:"medianRxRadioMCSRate24G"`
	MedianRxRadioMCSRate50G        int         `json:"medianRxRadioMCSRate50G"`
}

func (c *Client) genGetReq(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.BaseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	return req, nil
}

func (c *Client) addQS(r *http.Request, o RksOptions) {
	q := r.URL.Query()
	q.Add("serviceTicket", c.serviceTicket)
	if o.Index != "" {
		q.Add("index", o.Index)
	}
	if o.ListSize != "" {
		q.Add("listSize", o.ListSize)
	}
	if o.DomainID != "" {
		q.Add("domainId", o.DomainID)
	}
	r.URL.RawQuery = q.Encode()
}

// RksZone properties and fields of a Ruckus Controller Zone
type RksZone struct {
	ID          string `json:"id"`
	DomainID    string `json:"domainId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CountryCode string `json:"countryCode"`
	Version     string `json:"version"`
	Timezone    struct {
		SystemTimezone     string      `json:"systemTimezone"`
		CustomizedTimezone interface{} `json:"customizedTimezone"`
	} `json:"timezone"`
	IPMode                   string      `json:"ipMode"`
	Ipv6TrafficFilterEnabled interface{} `json:"ipv6TrafficFilterEnabled"`
	Login                    struct {
		ApLoginName     string `json:"apLoginName"`
		ApLoginPassword string `json:"apLoginPassword"`
	} `json:"login"`
	Mesh                       interface{} `json:"mesh"`
	DfsChannelEnabled          bool        `json:"dfsChannelEnabled"`
	CbandChannelEnabled        bool        `json:"cbandChannelEnabled"`
	CbandChannelLicenseEnabled bool        `json:"cbandChannelLicenseEnabled"`
	Channel144Enabled          bool        `json:"channel144Enabled"`
	Wifi24                     struct {
		AutoCellSizing        interface{} `json:"autoCellSizing"`
		TxPower               string      `json:"txPower"`
		ChannelWidth          int         `json:"channelWidth"`
		Channel               int         `json:"channel"`
		ChannelRange          []int       `json:"channelRange"`
		AvailableChannelRange []int       `json:"availableChannelRange"`
	} `json:"wifi24"`
	Wifi50 struct {
		AutoCellSizing               interface{} `json:"autoCellSizing"`
		TxPower                      string      `json:"txPower"`
		ChannelWidth                 int         `json:"channelWidth"`
		IndoorChannel                int         `json:"indoorChannel"`
		OutdoorChannel               int         `json:"outdoorChannel"`
		IndoorSecondaryChannel       interface{} `json:"indoorSecondaryChannel"`
		OutdoorSecondaryChannel      interface{} `json:"outdoorSecondaryChannel"`
		IndoorChannelRange           []int       `json:"indoorChannelRange"`
		OutdoorChannelRange          []int       `json:"outdoorChannelRange"`
		AvailableIndoorChannelRange  []int       `json:"availableIndoorChannelRange"`
		AvailableOutdoorChannelRange []int       `json:"availableOutdoorChannelRange"`
	} `json:"wifi50"`
	ProtectionMode24         string      `json:"protectionMode24"`
	Syslog                   interface{} `json:"syslog"`
	SmartMonitor             interface{} `json:"smartMonitor"`
	ClientAdmissionControl24 interface{} `json:"clientAdmissionControl24"`
	ClientAdmissionControl50 interface{} `json:"clientAdmissionControl50"`
	ChannelModeEnabled       bool        `json:"channelModeEnabled"`
	TunnelType               string      `json:"tunnelType"`
	TunnelProfile            struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tunnelProfile"`
	RuckusGreTunnelProfile struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"ruckusGreTunnelProfile"`
	SoftGreTunnelProflies interface{} `json:"softGreTunnelProflies"`
	IpsecProfiles         interface{} `json:"ipsecProfiles"`
	IpsecTunnelMode       interface{} `json:"ipsecTunnelMode"`
	BackgroundScanning24  struct {
		FrequencyInSec int `json:"frequencyInSec"`
	} `json:"backgroundScanning24"`
	BackgroundScanning50 struct {
		FrequencyInSec int `json:"frequencyInSec"`
	} `json:"backgroundScanning50"`
	ClientLoadBalancing24 interface{} `json:"clientLoadBalancing24"`
	ClientLoadBalancing50 interface{} `json:"clientLoadBalancing50"`
	BandBalancing         struct {
		Mode             string `json:"mode"`
		Wifi24Percentage int    `json:"wifi24Percentage"`
	} `json:"bandBalancing"`
	LoadBalancingMethod string `json:"loadBalancingMethod"`
	Rogue               struct {
		ReportType        string      `json:"reportType"`
		MaliciousTypes    interface{} `json:"maliciousTypes"`
		ProtectionEnabled bool        `json:"protectionEnabled"`
	} `json:"rogue"`
	LocationBasedService interface{} `json:"locationBasedService"`
	ApRebootTimeout      struct {
		GatewayLossTimeoutInSec int `json:"gatewayLossTimeoutInSec"`
		ServerLossTimeoutInSec  int `json:"serverLossTimeoutInSec"`
	} `json:"apRebootTimeout"`
	Location                          string      `json:"location"`
	LocationAdditionalInfo            string      `json:"locationAdditionalInfo"`
	Latitude                          interface{} `json:"latitude"`
	Longitude                         interface{} `json:"longitude"`
	VlanOverlappingEnabled            bool        `json:"vlanOverlappingEnabled"`
	NodeAffinityProfile               interface{} `json:"nodeAffinityProfile"`
	ZoneAffinityProfileID             string      `json:"zoneAffinityProfileId"`
	EnforcePriorityZoneAffinityEnable bool        `json:"enforcePriorityZoneAffinityEnable"`
	AwsVenue                          string      `json:"awsVenue"`
	VenueProfile                      interface{} `json:"venueProfile"`
	IpsecProfile                      interface{} `json:"ipsecProfile"`
	BonjourFencingPolicyEnabled       bool        `json:"bonjourFencingPolicyEnabled"`
	DhcpSiteConfig                    struct {
		SiteEnabled    bool        `json:"siteEnabled"`
		DwpdEnabled    interface{} `json:"dwpdEnabled"`
		ManualSelect   interface{} `json:"manualSelect"`
		SiteMode       interface{} `json:"siteMode"`
		SiteProfileIds interface{} `json:"siteProfileIds"`
		SiteAps        interface{} `json:"siteAps"`
		Eth0ProfileID  interface{} `json:"eth0ProfileId"`
		Eth1ProfileID  interface{} `json:"eth1ProfileId"`
	} `json:"dhcpSiteConfig"`
	BonjourFencingPolicy   interface{} `json:"bonjourFencingPolicy"`
	AutoChannelSelection24 struct {
		ChannelSelectMode string      `json:"channelSelectMode"`
		ChannelFlyMtbc    interface{} `json:"channelFlyMtbc"`
	} `json:"autoChannelSelection24"`
	AutoChannelSelection50 struct {
		ChannelSelectMode string      `json:"channelSelectMode"`
		ChannelFlyMtbc    interface{} `json:"channelFlyMtbc"`
	} `json:"autoChannelSelection50"`
	ChannelEvaluationInterval int `json:"channelEvaluationInterval"`
	ApMgmtVlan                struct {
		ID   int    `json:"id"`
		Mode string `json:"mode"`
	} `json:"apMgmtVlan"`
	ApLatencyInterval struct {
		PingEnabled bool `json:"pingEnabled"`
	} `json:"apLatencyInterval"`
	Altitude struct {
		AltitudeUnit  string      `json:"altitudeUnit"`
		AltitudeValue interface{} `json:"altitudeValue"`
	} `json:"altitude"`
	RecoverySsid struct {
		RecoverySsidEnabled bool `json:"recoverySsidEnabled"`
	} `json:"recoverySsid"`
	DosBarringEnable      int `json:"dosBarringEnable"`
	DosBarringPeriod      int `json:"dosBarringPeriod"`
	DosBarringThreshold   int `json:"dosBarringThreshold"`
	DosBarringCheckPeriod int `json:"dosBarringCheckPeriod"`
	SnmpAgent             struct {
		ApSnmpEnabled bool          `json:"apSnmpEnabled"`
		SnmpV2Agent   []interface{} `json:"snmpV2Agent"`
		SnmpV3Agent   []interface{} `json:"snmpV3Agent"`
	} `json:"snmpAgent"`
	ClusterRedundancyEnabled                   bool        `json:"clusterRedundancyEnabled"`
	AaaAffinityEnabled                         bool        `json:"aaaAffinityEnabled"`
	RogueApReportThreshold                     int         `json:"rogueApReportThreshold"`
	RogueApAggressivenessMode                  int         `json:"rogueApAggressivenessMode"`
	RogueApJammingDetection                    bool        `json:"rogueApJammingDetection"`
	RogueApJammingThreshold                    interface{} `json:"rogueApJammingThreshold"`
	DirectedMulticastFromWiredClientEnabled    bool        `json:"directedMulticastFromWiredClientEnabled"`
	DirectedMulticastFromWirelessClientEnabled bool        `json:"directedMulticastFromWirelessClientEnabled"`
	DirectedMulticastFromNetworkEnabled        bool        `json:"directedMulticastFromNetworkEnabled"`
	HealthCheckSitesEnabled                    bool        `json:"healthCheckSitesEnabled"`
	HealthCheckSites                           []string    `json:"healthCheckSites"`
	SSHTunnelEncryption                        string      `json:"sshTunnelEncryption"`
	LteBandLockChannels                        []struct {
		SimCardID int    `json:"simCardId"`
		Type      string `json:"type"`
		Channel4G string `json:"channel4g"`
		Channel3G string `json:"channel3g"`
	} `json:"lteBandLockChannels"`
	ApHccdEnabled bool `json:"apHccdEnabled"`
	ApHccdPersist bool `json:"apHccdPersist"`
}
