package main

import (
  "fmt"
  "net"
  "os/exec"
  "log"
  "bufio"
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
  // append comment to end
  return s, nil
}

func (s *SshPublicKey) String() string {
  return fmt.Sprintf("FIXME")
}

func RunCommand(cmdstr string) error {
  cmd := exec.Command("sh", "-c", cmdstr)
  log.Printf("Running: %s\n", cmdstr)
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return err
  }
  stderr, err := cmd.StderrPipe()
  if err != nil {
    return err
  }
  if err := cmd.Start(); err != nil {
    return err
  }
  se := bufio.NewScanner(stderr)
  so := bufio.NewScanner(stdout)
  logIO := func(s *bufio.Scanner, prefix string) {
    for s.Scan() {
      log.Printf("%s%s\n", prefix, s.Bytes())
    }
  }
  go logIO(se, "! ")
  go logIO(so, "> ")
  return cmd.Wait()
}
