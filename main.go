package main

import (
  "log"
  "flag"
)

var (
  vipConfigPath string
  vipConfig VipConfig
  scanTimeout int
)

func init(){
  flag.StringVar(&vipConfigPath, "vips", "", "vips json config file")
  flag.IntVar(&scanTimeout, "timeout", 20, "timeout in seconds to scan a single host")
  flag.Parse()
}

func main(){
  if vipConfigPath == "" {
    log.Fatal("You need to pass a -vips config")
  }
  vipConfig, err := LoadVipConfig(vipConfigPath)
  if err != nil {
    log.Fatalf("Unable to parse vip config: %s", err)
  }
  log.Printf("Loaded %d vips to scan from config\n", len(vipConfig.Vips))

  for _, host := range vipConfig.Vips {
    log.Printf("Scanning %s\n", host)
    k, err := ScanHost(host, scanTimeout)
    if err != nil {
      log.Printf("Error: %s\n", err)
    } else {
      log.Printf("Success: %s: %s - %s\n", k.Hostname, k.Ips, k.PublicKey)
    }
  }

}

