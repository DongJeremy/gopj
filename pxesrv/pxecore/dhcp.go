package pxecore

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"honnef.co/go/tools/config"
)

var log = GetLogger("coredhcp")

// Handler4 behaves like Handler6, but for DHCPv4 packets.
type Handler4 func(req, resp *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, bool)

// Server is a CoreDHCP server structure that holds information about
// DHCPv4 servers, and their respective handlers.
type Server struct {
	Handlers4 []Handler4
	Config    *config.Config
	Server4   *server4.Server
	errors    chan error
}

// MainHandler4 is like MainHandler6, but for DHCPv4 packets.
func (s *Server) MainHandler4(conn net.PacketConn, _peer net.Addr, req *dhcpv4.DHCPv4) {
	var (
		resp, tmp *dhcpv4.DHCPv4
		err       error
		stop      bool
	)
	if req.OpCode != dhcpv4.OpcodeBootRequest {
		log.Printf("MainHandler4: unsupported opcode %d. Only BootRequest (%d) is supported", req.OpCode, dhcpv4.OpcodeBootRequest)
		return
	}
	tmp, err = dhcpv4.NewReplyFromRequest(req)
	if err != nil {
		log.Printf("MainHandler4: failed to build reply: %v", err)
		return
	}
	switch mt := req.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		tmp.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	case dhcpv4.MessageTypeRequest:
		tmp.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	default:
		log.Printf("plugins/server: Unhandled message type: %v", mt)
		return
	}

	resp = tmp
	for _, handler := range s.Handlers4 {
		resp, stop = handler(req, resp)
		if stop {
			break
		}
	}

	if resp != nil {
		var peer net.Addr
		if !req.GatewayIPAddr.IsUnspecified() {
			// TODO: make RFC8357 compliant
			peer = &net.UDPAddr{IP: req.GatewayIPAddr, Port: dhcpv4.ServerPort}
		} else if resp.MessageType() == dhcpv4.MessageTypeNak {
			peer = &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ClientPort}
		} else if !req.ClientIPAddr.IsUnspecified() {
			peer = &net.UDPAddr{IP: req.ClientIPAddr, Port: dhcpv4.ClientPort}
		} else if req.IsBroadcast() {
			peer = &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ClientPort}
		} else {
			// FIXME: we're supposed to unicast to a specific *L2* address, and an L3
			// address that's not yet assigned.
			// I don't know how to do that with this API...
			//peer = &net.UDPAddr{IP: resp.YourIPAddr, Port: dhcpv4.ClientPort}
			log.Warn("Cannot handle non-broadcast-capable unspecified peers in an RFC-compliant way. " +
				"Response will be broadcast")

			peer = &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ClientPort}
		}

		if _, err := conn.WriteTo(resp.ToBytes(), peer); err != nil {
			log.Printf("MainHandler4: conn.Write to %v failed: %v", peer, err)
		}

	} else {
		log.Print("MainHandler4: dropping request because response is nil")
	}
}

// LoadHandlers load handlers
func (s *Server) LoadHandlers(conf *config.Config) {
	s.Handlers4 = append(s.Handlers4, h6)
}

// Start will start the server asynchronously. See `Wait` to wait until
// the execution ends.
func (s *Server) Start() error {
	_, _, err := s.LoadPlugins(s.Config)
	if err != nil {
		return err
	}

	// listen
	log.Printf("Starting DHCPv4 listener on %v", s.Config.Server4.Listener)
	s.Server4, err = server4.NewServer(s.Config.Server4.Interface, s.Config.Server4.Listener, s.MainHandler4)
	if err != nil {
		return err
	}
	go func() {
		s.errors <- s.Server4.Serve()
	}()

	return nil
}

// Wait waits until the end of the execution of the server.
func (s *Server) Wait() error {
	log.Print("Waiting")
	err := <-s.errors
	if s.Server4 != nil {
		s.Server4.Close()
	}
	return err
}

// NewServer creates a Server instance with the provided configuration.
func NewServer(config *Config) *Server {
	return &Server{Config: config, errors: make(chan error, 1)}
}
