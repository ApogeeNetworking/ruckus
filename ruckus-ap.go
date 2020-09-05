package ruckus

import (
	"encoding/json"
	"fmt"
	"strconv"
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
