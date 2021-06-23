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
func New(apiVersion, host, user, pass string, ignoreSSL bool) *Client {
	return &Client{
		BaseURL:  fmt.Sprintf("https://%s:8443/wsg/api/public/v%s", host, apiVersion),
		host:     host,
		username: user,
		password: pass,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: ignoreSSL,
				},
			},
			Timeout: 120 * time.Second,
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
