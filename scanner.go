package main

import (
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
	Type      string
	PublicKey string
}

func ScanHost(host string, timeout int) (SshPublicKey, error) {
	scan_type := "rsa"
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
	scan_output, err := RunCommand(fmt.Sprintf("ssh-keyscan -t %s %s", scan_type, aliaslist))
	if err != nil {
		return s, fmt.Errorf("Unable to scan %s: %s", s.Hostname, err)
	}
	lines := strings.Split(scan_output, "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, s.Hostname) {
			//github.com,192.30.252.131 ssh-rsa <key>
			fields := strings.Fields(l)
			if len(fields) < 3 {
				return s, fmt.Errorf("Unexpected number of fields received from ssh-keyscan")
			}
			s.Type = fields[1]
			s.PublicKey = fields[2]
			break
		}
	}
	if s.PublicKey == "" {
		return s, fmt.Errorf("Unable to extract public key from ssh-keyscan")
	}
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
	return fmt.Sprintf("%s %s %s", strings.Join(s.AliasList(), ","), s.Type, s.PublicKey)
}

func RunCommand(cmdstr string) (string, error) {
	output := ""
	cmd := exec.Command("sh", "-c", cmdstr)
	log.Printf("Running: %s\n", cmdstr)
	combined, err := cmd.CombinedOutput()
	// []bytes are annoying as shit. just use string
	output = string(combined)
	return output, err
}
