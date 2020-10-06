package ruckus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// RksAP ruckus ap properties
type RksAP struct {
	MacAddr   string `json:"mac"`
	ZoneID    string `json:"zoneId"`
	ApGroupID string `json:"apGroupId"`
	Serial    string `json:"serial"`
	Name      string `json:"name"`
	LanPorts  int    `json:"lanPortSize"`
}

// GetAPs retrieves APs associated with the Controller
func (c *Client) GetAPs(o RksOptions) ([]RksAP, error) {
	var getMore func(o RksOptions, r []RksAP) ([]RksAP, error)
	getMore = func(o RksOptions, rksAps []RksAP) ([]RksAP, error) {
		req, err := c.genGetReq("/aps")
		if err != nil {
			return rksAps, err
		}
		c.addQS(req, o)

		res, err := c.http.Do(req)
		if err != nil {
			return rksAps, fmt.Errorf("failed to get resp: %v", err)
		}
		defer res.Body.Close()
		type rksApResult struct {
			RksCommonReq
			List []RksAP `json:"list"`
		}
		var aps rksApResult
		json.NewDecoder(res.Body).Decode(&aps)
		for _, ap := range aps.List {
			rksAps = append(rksAps, ap)
		}
		if aps.HasMore {
			i := aps.FirstIndex + 100
			o.Index = strconv.Itoa(i)
			return getMore(o, rksAps)
		}
		return rksAps, nil
	}
	return getMore(o, []RksAP{})
}

// GetApGroups retrieves list of AP Group Names with IDs
func (c *Client) GetApGroups(o RksOptions, zoneID string) ([]RksObject, error) {
	if c.serviceTicket == "" {
		return nil, fmt.Errorf(loginErr)
	}
	ep := fmt.Sprintf("/rkszones/%s/apgroups", zoneID)
	req, err := c.genGetReq(ep)
	if err != nil {
		return nil, err
	}
	c.addQS(req, o)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get resp: %v", err)
	}
	defer res.Body.Close()
	var grps RksCommonRes
	json.NewDecoder(res.Body).Decode(&grps)
	fmt.Println(grps.TotalCount)
	return grps.List, nil
}

// ApIntf ...
type ApIntf struct {
	MacAddr string `json:"apMac"`
	Speed   string `json:"phyLink"`
	Status  string `json:"logicLink"`
	Duplex  string
}

// GetApIntf ...
func (c *Client) GetApIntf(macAddr string) (ApIntf, error) {
	uri := fmt.Sprintf("https://%s:8443/wsg/api/scg/aps/%s", c.host, macAddr)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return ApIntf{}, err
	}
	c.addQS(req, RksOptions{})
	res, err := c.http.Do(req)
	if err != nil {
		return ApIntf{}, err
	}
	defer res.Body.Close()
	type getApRESP struct {
		Success bool `json:"success"`
		Data    struct {
			LanPorts []ApIntf `json:"lanPortStatus"`
		} `json:"data"`
	}
	var apResp getApRESP
	json.NewDecoder(res.Body).Decode(&apResp)
	var apIntf ApIntf
	if apResp.Success {
		for _, port := range apResp.Data.LanPorts {
			status := strings.ToLower(port.Status)
			if status == "up" {
				speedClean := strings.Replace(
					port.Speed, "Up ", "", -1,
				)
				speedSplit := strings.Split(speedClean, " ")
				apIntf = ApIntf{
					MacAddr: port.MacAddr,
					Speed:   speedSplit[0],
					Duplex:  strings.ToUpper(speedSplit[1]),
					Status:  status,
				}
				break
			}
			apIntf = ApIntf{
				MacAddr: port.MacAddr,
				Speed:   status,
				Status:  status,
			}
		}
	}
	return apIntf, nil
}

// ApLldp ...
type ApLldp struct {
	RemoteHostname string `json:"lldpSysName"`
	RemoteIntf     string `json:"lldpPortDesc"`
	RemoteIP       string `json:"lldpMgmtIP"`
}

// GetApLldp ...
func (c *Client) GetApLldp(macAddr string) (ApLldp, error) {
	uri := fmt.Sprintf("/aps/%s/apLldpNeighbors", macAddr)
	req, err := c.genGetReq(uri)
	if err != nil {
		return ApLldp{}, err
	}
	c.addQS(req, RksOptions{})
	res, err := c.http.Do(req)
	if err != nil {
		return ApLldp{}, err
	}
	defer res.Body.Close()
	type getApRESP struct {
		Count int      `json:"totalCount"`
		List  []ApLldp `json:"list"`
	}
	var apResp getApRESP
	var apLldp ApLldp
	if err := json.NewDecoder(res.Body).Decode(&apResp); err != nil {
		return apLldp, err
	}
	if apResp.Count == 0 {
		return apLldp, nil
	}
	for _, lldp := range apResp.List {
		apLldp = lldp
		break
	}
	return apLldp, nil
}
