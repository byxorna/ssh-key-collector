package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sort"
)

var (
	vipConfigPath string
	vipConfig     VipConfig
	scanTimeout   int
	outFile       string
)

func init() {
	flag.StringVar(&vipConfigPath, "vips", "", "vips json config file")
	flag.IntVar(&scanTimeout, "timeout", 20, "timeout in seconds to scan a single host")
	flag.StringVar(&outFile, "out", "ssh_known_hosts", "file to write scanned keys to")
	flag.Parse()
}

func main() {
	if vipConfigPath == "" {
		log.Fatal("You need to pass a -vips config")
	}
	vipConfig, err := LoadVipConfig(vipConfigPath)
	if err != nil {
		log.Fatalf("Unable to parse vip config: %s", err)
	}
	log.Printf("Loaded %d vips to scan from config\n", len(vipConfig.Vips))

	var keys SshPublicKeyList
	var failed []string
	for _, host := range vipConfig.Vips {
		log.Printf("Scanning %s\n", host)
		k, err := ScanHost(host, scanTimeout)
		if err != nil {
			log.Printf("Error: %s\n", err)
			failed = append(failed, host)
		} else {
			keys = append(keys, k)
			log.Printf("Success: %s\n", k.String())
		}
	}

	if len(failed) > 0 {
		log.Printf("%d hosts failed to scan\n", len(failed))
		os.Exit(2)
	}

	sort.Sort(keys)

	f, err := os.Create(outFile)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	for _, k := range keys {
		w.WriteString(k.String())
	}
	w.Flush()

}
