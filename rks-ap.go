package rkssz

import (
	"encoding/json"
	"fmt"
)

// RksAP ruckus ap properties
type RksAP struct {
	MacAddr   string `json:"mac"`
	ZoneID    string `json:"zoneId"`
	ApGroupID string `json:"apGroupId"`
	Serial    string `json:"serial"`
	Name      string `json:"name"`
}

// RksApRes ... GetAP Result
type RksApRes struct {
	RksCommonReq
	List []RksAP `json:"list"`
}

// GetAPs retrieves APs associated with the Controller
func (c *Client) GetAPs(o RksOptions) (RksApRes, error) {
	req, err := c.genGetReq("/aps")
	if err != nil {
		return RksApRes{}, err
	}
	c.addQS(req, o)

	res, err := c.http.Do(req)
	if err != nil {
		return RksApRes{}, fmt.Errorf("failed to get resp: %v", err)
	}
	defer res.Body.Close()
	var rksAps RksApRes
	json.NewDecoder(res.Body).Decode(&rksAps)
	return rksAps, nil
}
