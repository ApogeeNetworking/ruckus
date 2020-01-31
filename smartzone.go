package rkssz

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client struct is used to handle the Connection with the SmartZone Controller
type Client struct {
	BaseURL  string
	Username string
	Password string

	http          *http.Client
	serviceTicket string
}

// New creates a Reference to a Client
func New(host, user, pass string, ignoreSSL bool) *Client {
	return &Client{
		BaseURL:  fmt.Sprintf("https://%s:8443/wsg/api/public/v8_1", host),
		Username: user,
		Password: pass,
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

// SZAuthObj smart zone authorization object
type SZAuthObj struct {
	ControllerVersion string `json:"controllerVersion"`
	ServiceTicket     string `json:"serviceTicket"`
}

// Login est a session with the Ruckus SZ Controller
func (c *Client) Login() error {
	type creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	authObj := creds{
		Username: c.Username,
		Password: c.Password,
	}
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

	var auth SZAuthObj
	json.NewDecoder(res.Body).Decode(&auth)
	c.serviceTicket = auth.ServiceTicket
	fmt.Println(auth)
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

// RksZone the object returned when retrieving Rks Zones
type RksZone struct {
	TotalCount int  `json:"totalCount"`
	HasMore    bool `json:"hasMore"`
	FirstIndex int  `json:"firstIndex"`
	List       []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"list"`
}

// ZoneOptions query options for GetZones
type ZoneOptions struct {
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

// GetZones retrieves a Paginated List of Zones
func (c *Client) GetZones(o ZoneOptions) (RksZone, error) {
	if c.serviceTicket == "" {
		e := "you must first login to perform this action"
		return RksZone{}, fmt.Errorf(e)
	}
	req, err := http.NewRequest("GET", c.BaseURL+"/rkszones", nil)
	if err != nil {
		return RksZone{}, fmt.Errorf("failed to create request: %v", err)
	}
	q := req.URL.Query()
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
	req.URL.RawQuery = q.Encode()

	res, err := c.http.Do(req)
	if err != nil {
		return RksZone{}, fmt.Errorf("failed to get resp: %v", err)
	}
	defer res.Body.Close()
	var zones RksZone
	json.NewDecoder(res.Body).Decode(&zones)

	return zones, nil
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
