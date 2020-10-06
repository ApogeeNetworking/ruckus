package ruckus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//SetApNameAndGroup ...
func (c *Client) SetApNameAndGroup(apMacAddr, apName, zoneID, groupID string) {
	type apChngREQ struct {
		ZoneID  string `json:"zoneId"`
		GroupID string `json:"apGroupId"`
		ApName  string `json:"name"`
	}
	chngObj := apChngREQ{
		ZoneID:  zoneID,
		GroupID: groupID,
		ApName:  apName,
	}
	jdata, _ := json.Marshal(&chngObj)
	chngBody := strings.NewReader(string(jdata))
	ep := fmt.Sprintf("/aps/%s", apMacAddr)
	req, err := http.NewRequest("PATCH", c.BaseURL+ep, chngBody)
	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
		return
	}
	c.addQS(req, RksOptions{})
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := c.http.Do(req)
	if err != nil {
		fmt.Printf("error in response: %v\n", err)
		return
	}
	defer res.Body.Close()
	d, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(d))
}

// RksQuery ...
type RksQuery struct {
	Filters       []Mapper    `json:"filters"`
	FullTxtSearch Mapper      `json:"fullTextSearch"`
	Attrs         []string    `json:"attributes"`
	SortInfo      rksSortInfo `json:"sortInfo"`
	Page          int         `json:"page"`
	Limit         int         `json:"limit"`
}

type rksSortInfo struct {
	SortCol   string `json:"sortColumn"`
	Direction string `json:"dir"`
}

// GetAPs retrieves APs associated with the Controller
func (c *Client) GetAPs(o RksOptions) ([]RksAP, error) {
	var getMore func(o RksOptions, r []RksAP) ([]RksAP, error)
	getMore = func(o RksOptions, rksAps []RksAP) ([]RksAP, error) {
		q := RksQuery{
			Filters:       []Mapper{},
			FullTxtSearch: Mapper{Type: "AND", Value: ""},
			Attrs:         []string{"*"},
			SortInfo:      rksSortInfo{SortCol: "apMac", Direction: "ASC"},
			Page:          1,
			Limit:         10000,
		}
		qjson, _ := json.Marshal(&q)
		body := strings.NewReader(string(qjson))
		req, err := http.NewRequest("POST", c.BaseURL+"/query/ap", body)
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

// GetAp ...
func (c *Client) GetAp(macAddr string) (string, error) {
	req, err := c.genGetReq("/aps/" + macAddr)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	c.addQS(req, RksOptions{})
	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	type getApRESP struct {
		Model string `json:"model"`
	}
	var apSum getApRESP
	json.NewDecoder(res.Body).Decode(&apSum)
	return apSum.Model, nil
}

// GetApGroupName ...
func (c *Client) GetApGroupName(zoneID, groupID string) (string, error) {
	ep := fmt.Sprintf("/rkszones/%s/apgroups/%s", zoneID, groupID)
	req, err := c.genGetReq(ep)
	if err != nil {
		return "", err
	}
	c.addQS(req, RksOptions{})
	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	type grpRES struct {
		GroupName string `json:"name"`
	}
	var grpNameRes grpRES
	json.NewDecoder(res.Body).Decode(&grpNameRes)
	return grpNameRes.GroupName, nil
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
	return grps.List, nil
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
		RksCommonReq
		List []ApLldp `json:"list"`
	}
	var apResp getApRESP
	var apLldp ApLldp
	if err := json.NewDecoder(res.Body).Decode(&apResp); err != nil {
		return apLldp, err
	}
	if apResp.TotalCount == 0 {
		return apLldp, nil
	}
	for _, lldp := range apResp.List {
		apLldp = lldp
		break
	}
	return apLldp, nil
}

// RebootAp ...
func (c *Client) RebootAp(macAddr string) (bool, error) {
	uri := fmt.Sprintf("https://%s:8443/wsg/api/scg/aps/%s/reboot", c.host, macAddr)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return false, err
	}
	c.addQS(req, RksOptions{})
	res, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	type apRebootRES struct {
		Success bool `json:"success"`
	}
	var result apRebootRES
	json.NewDecoder(res.Body).Decode(&result)
	return result.Success, nil
}
