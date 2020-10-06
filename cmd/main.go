package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ApogeeNetworking/ruckus"
	"github.com/subosito/gotenv"
)

var host, user, pass string

func init() {
	gotenv.Load()
	host = os.Getenv("RKS_HOST")
	user = os.Getenv("RKS_USER")
	pass = os.Getenv("RKS_PASS")
}

func main() {
	sz := ruckus.New(host, user, pass, true)
	err := sz.Login()
	if err != nil {
		log.Fatalf("login failed: %v", err)
	}
	defer func() {
		err := sz.Logout()
		if err != nil {
			// Why did we err Logging Out
			fmt.Println(err)
		}
	}()
	type ApChangeREQ struct {
		MacAddr   string
		GroupName string
		ApName    string
		GroupID   string
		ZoneID    string
	}
	// apChng := ApChangeREQ{
	// 	MacAddr:   "60:D0:2C:2A:52:B0",
	// 	ApName:    "ap01.austin.introom.saml.tx",
	// 	GroupName: "Dev-Lab-Zone1-Group3",
	// }
	// zoneIds := rkszones(sz)
	// for _, zoneID := range zoneIds {
	// 	apGroups, _ := sz.GetApGroups(ruckus.RksOptions{}, zoneID)
	// 	for _, apGroup := range apGroups {
	// 		if apChng.GroupName == apGroup.Name {
	// 			apChng.GroupID = apGroup.ID
	// 			apChng.ZoneID = zoneID
	// 		}
	// 	}
	// }
	// sz.SetApNameAndGroup(
	// 	apChng.MacAddr,
	// 	apChng.ApName,
	// 	apChng.ZoneID,
	// 	apChng.GroupID,
	// )

	// mac := "EC:58:EA:0A:24:D0"
	mac := "60:d0:2c:2a:52:b0"
	sz.RebootAp(mac)
	// sz.GetAp(mac)
	// apIntf, err := sz.GetApLldp(mac)
	// fmt.Println(apIntf)
	// aps(sz)
}

func aps(sz *ruckus.Client) {
	rkAps, err := sz.GetAPs(ruckus.RksOptions{})
	if err != nil {
		sz.Logout()
		log.Fatalf("%v", err)
	}
	for _, ap := range rkAps {
		fmt.Println(ap)
	}
}

func rkszones(sz *ruckus.Client) []string {
	zones, err := sz.GetZones(ruckus.RksOptions{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	var zoneIds []string
	for _, zone := range zones.List {
		zoneIds = append(zoneIds, zone.ID)
	}
	return zoneIds
	// fmt.Println(zones.List[0].ID)
}

func rkszone(sz *ruckus.Client) {
	zone, _ := sz.GetZone("615d18e9-0cc0-4e3d-b98e-f2a476b5a846")
	fmt.Println(zone)
}
