package pxecore

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var conf Config

// Config config for pxe
type Config struct {
	HTTP HTTP `yaml:"http"`
	TFTP TFTP `yaml:"tftp"`
	DHCP DHCP `yaml:"dhcp"`
}

// HTTP config
type HTTP struct {
	// which ip address that http server listening
	HTTPIP    string `yaml:"listen_ip,omitempty"`
	HTTPPort  string `yaml:"listen_port,omitempty"` // listening port of http server
	MountPath string `yaml:"mount_path,omitempty"`  // http file server path
}

// TFTP config
type TFTP struct {
	TftpPath string `yaml:"mount_path,omitempty"` // tftp_files server path
	TftpIP   string `yaml:"listen_ip,omitempty"`  // ip address that tftp_files server listening on
}

// DHCP config
type DHCP struct {
	ListenIP   string `yaml:"listen_ip,omitempty"` // which ip address that dhcp server was listening on
	ListenPort string `yaml:"listen_port,omitempty"`
	TftpServer string `yaml:"tftp_server,omitempty"`
	StartIP    string `yaml:"start_ip"`
	Range      int    `yaml:"lease_range"`       // lease ip address count
	NetMask    string `yaml:"netmask,omitempty"` // default /24
	PxeFile    string `yaml:"pxe_file"`          // pxe file name
}

// Refresh runtime configurations
func Refresh() {
	c := new(Config)
	// set default options
	c.HTTP.HTTPIP = "0.0.0.0"
	c.HTTP.HTTPPort = "80"
	c.HTTP.MountPath = "/mnt/dhtp/http"
	c.TFTP.TftpIP = "0.0.0.0"
	c.TFTP.TftpPath = "/mnt/dhtp/tftp"
	c.DHCP.ListenIP = "0.0.0.0"
	c.DHCP.ListenPort = "67"
	c.DHCP.StartIP = "169.169.181.2"
	c.DHCP.Range = 50
	c.DHCP.PxeFile = "pxelinux.0"
	c.DHCP.NetMask = "255.255.255.0"
	f, err := ioutil.ReadFile("/etc/dhtp/dhtp.yml")
	if err != nil {
		panic(fmt.Sprintf("read config file from /etc/dhtp/dhtp.conf failed, %s", err))
	}
	err = yaml.Unmarshal(f, c)
	if err != nil {
		panic(fmt.Sprintf("parse config file failed, %s", err))
	}
	conf = *c
}

// GetConf return runtime configurations
func GetConf() Config {
	return conf
}
