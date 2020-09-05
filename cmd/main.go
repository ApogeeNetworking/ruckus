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
	host = os.Getenv("HOST")
	user = os.Getenv("USER")
	pass = os.Getenv("PASS")
}

func main() {
	sz := ruckus.New(
		host,
		user,
		pass,
		true,
	)
	err := sz.Login()
	if err != nil {
		log.Fatalf("login failed: %v", err)
	}
	aps(sz)
	// rkszone(sz)
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
	fmt.Println(zones.List[0].ID)
}

func rkszone(sz *ruckus.Client) {
	zone, _ := sz.GetZone("615d18e9-0cc0-4e3d-b98e-f2a476b5a846")
	fmt.Println(zone)
}
