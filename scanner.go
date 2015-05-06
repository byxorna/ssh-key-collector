package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

type SshPublicKey struct {
	Hostname  string
	Ips       []net.IP
	Aliases   []string
	PublicKey string
}

func ScanHost(host string, timeout int) (SshPublicKey, error) {
	s := SshPublicKey{Hostname: host}
	// first, resolve the host to IP[s]. Some VIPs may map to multiple backends, so track them all
	ips, err := net.LookupIP(host)
	if err != nil {
		return s, fmt.Errorf("Unable to resolve %s to IP: %s", host, err)
	}
	s.Ips = ips
	// ssh-keyscan -t rsa,dsa,ecdsa $(hostname -f),172.18.110.119,fuckface
	aliases := s.AliasList()
	aliaslist := strings.Join(aliases, ",")
	scan_output, err := RunCommand(fmt.Sprintf("ssh-keyscan -t rsa %s", aliaslist))
	if err != nil {
		return s, fmt.Errorf("Unable to scan %s: %s", s.Hostname, err)
	}
	//TODO format the scan_output properly, handle failures, etc
	s.PublicKey = scan_output
	return s, nil
}

func (s *SshPublicKey) AliasList() []string {
	/// first we want the canonical hostname, followed by aliases, followed by any IPs
	aliases := []string{}
	aliases = append(aliases, s.Hostname)
	aliases = append(aliases, s.Aliases...)
	for _, val := range s.Ips {
		aliases = append(aliases, val.String())
	}
	return aliases
}

func (s *SshPublicKey) String() string {
	return fmt.Sprintf("FIXME")
}

func RunCommand(cmdstr string) (string, error) {
	output := ""
	cmd := exec.Command("sh", "-c", cmdstr)
	log.Printf("Running: %s\n", cmdstr)
	combined, err := cmd.CombinedOutput()
	// figure out the index of the null terminator
	combined_len := bytes.Index(combined, []byte{0})
	output = string(combined[:combined_len])
	return output, err
}
