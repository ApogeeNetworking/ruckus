# Ruckus SmartZone Controller API Wrapper

This is not a comprehensive|exhaustive library

## Usage

```go
package main

import (
    "fmt"

    "github.com/ApogeeNetworking/ruckus"
)

func main() {
    var ignoreSSL bool = true
    // Instantiate SmartZone Client Struct
    smartZone := ruckus.New("ip_address", "username", "password", ignoreSSL)
    // Perform Login to SmartZone (issues a ServiceTicket hidden by API)
    err := smartZone.Login()
    if err != nil {
        log.Fatalf("login failed: %v", err)
    }
    // Ensure you logout when done
    defer smartZone.Logout()

    rksAps, err := smartZone.GetAPs(ruckus.RksOptions{})
    if err != nil {
        // Handle Error ...
    }
    /*
    The Definition of RksAP can be found in types.go but consists of:
    MacAddr
    ZoneID
    ZoneName
    GroupID
    GroupName
    Serial
    ApName
    Model
    Status (Online|Offline|Flag)
    IPAddr
    ExtIPAddr
    Firmware
    */
}
```