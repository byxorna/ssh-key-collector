package main

import (
  "encoding/json"
  "os"
)

// list of fqdn
type HostList []string
type VipConfig struct {
  Vips HostList `json:"vips"`
}

func LoadVipConfig(path string) (VipConfig, error) {
  c := VipConfig{}
  f, err := os.Open(path)
  if err != nil {
    return c, err
  }
  defer f.Close()
  err = json.NewDecoder(f).Decode(&c)
  if err != nil {
    return c, err
  }
  return c, nil
}
