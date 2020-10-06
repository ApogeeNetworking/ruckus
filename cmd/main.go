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
	// mac := "EC:58:EA:0A:24:D0"
	// mac := "60:d0:2c:2a:52:b0"
	// apIntf, err := sz.GetApLldp(mac)
	// fmt.Println(apIntf)
	rkszones(sz)
	err = sz.Logout()
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

func rkszones(sz *ruckus.Client) {
	zones, err := sz.GetZones(ruckus.RksOptions{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, zone := range zones.List {
		fmt.Println(zone)
	}
	// fmt.Println(zones.List[0].ID)
}

func rkszone(sz *ruckus.Client) {
	zone, _ := sz.GetZone("615d18e9-0cc0-4e3d-b98e-f2a476b5a846")
	fmt.Println(zone)
}
