package sip

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// UDPConnInterface abstracts the UDP connection methods needed by the server
type UDPConnInterface interface {
	ReadFromUDP(b []byte) (int, *net.UDPAddr, error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (int, error)
	Close() error
}

// Server represents a SIP server
type Server struct {
	Port      string
	BindAddr  string
	conn      UDPConnInterface
	registrar map[string]string // user -> address mapping
	calls     map[string]string // callID -> status mapping
}

// NewServer creates a new SIP server instance
func NewServer(port string) *Server {
	return &Server{
		Port:      port,
		BindAddr:  "0.0.0.0",
		registrar: make(map[string]string),
		calls:     make(map[string]string),
	}
}

// SetBindAddr sets the bind address for the server
func (s *Server) SetBindAddr(addr string) {
	s.BindAddr = addr
}

// Start begins listening for SIP messages
func (s *Server) Start() error {
	// Combine bind address and port
	listenAddr := fmt.Sprintf("%s:%s", s.BindAddr, s.Port)
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return fmt.Errorf("address resolution error: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("UDP listening error: %v", err)
	}
	s.conn = conn

	log.Printf("SIP server started on %s", listenAddr)

	buffer := make([]byte, 65535)
	for {
		n, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("packet reading error: %v", err)
			continue
		}

		// Process received message in a separate goroutine
		go s.handleMessage(addr, buffer[:n])
	}
}

// handleMessage processes incoming SIP messages
func (s *Server) handleMessage(addr *net.UDPAddr, data []byte) {
	msgStr := string(data)
	log.Printf("received message: %s", msgStr)

	msg, err := ParseMessage(msgStr)
	if err != nil {
		log.Printf("message parsing error: %v", err)
		return
	}

	// Process based on message type
	if strings.HasPrefix(msg.StartLine, "REGISTER") {
		s.handleRegister(addr, msg)
	} else if strings.HasPrefix(msg.StartLine, "INVITE") {
		s.handleInvite(addr, msg)
	} else if strings.HasPrefix(msg.StartLine, "BYE") {
		s.handleBye(addr, msg)
	} else if strings.HasPrefix(msg.StartLine, "ACK") {
		// ACK typically doesn't require a response
		log.Printf("ACK received: %s", msg.Headers["Call-ID"])
	} else {
		log.Printf("unhandled message type: %s", msg.StartLine)
	}
}

// handleRegister processes REGISTER requests
func (s *Server) handleRegister(addr *net.UDPAddr, msg *Message) {
	// Extract user information from From header
	fromHeader := msg.Headers["From"]
	uri := extractSIPURI(fromHeader)

	// Register the user
	s.registrar[uri] = addr.String()
	log.Printf("user registered: %s -> %s", uri, addr.String())

	// Send 200 OK response
	resp := NewResponse("200", "OK", msg)
	s.sendResponse(addr, resp)
}

// handleInvite processes INVITE requests
func (s *Server) handleInvite(addr *net.UDPAddr, msg *Message) {
	callID := msg.Headers["Call-ID"]

	// Send 100 Trying response
	tryingResp := NewResponse("100", "Trying", msg)
	s.sendResponse(addr, tryingResp)

	// Send 180 Ringing response
	ringingResp := NewResponse("180", "Ringing", msg)
	s.sendResponse(addr, ringingResp)

	// Save call state
	s.calls[callID] = "ringing"

	// Send 200 OK response (normally sent after user accepts call)
	okResp := NewResponse("200", "OK", msg)
	s.sendResponse(addr, okResp)

	// Update call state
	s.calls[callID] = "connected"
	log.Printf("call established: %s", callID)
}

// handleBye processes BYE requests
func (s *Server) handleBye(addr *net.UDPAddr, msg *Message) {
	callID := msg.Headers["Call-ID"]

	// Check call state
	if _, exists := s.calls[callID]; exists {
		// Terminate the call
		delete(s.calls, callID)
		log.Printf("call terminated: %s", callID)
	}

	// Send 200 OK response
	resp := NewResponse("200", "OK", msg)
	s.sendResponse(addr, resp)
}

// sendResponse sends a SIP response message
func (s *Server) sendResponse(addr *net.UDPAddr, msg *Message) {
	msgStr := msg.String()
	_, err := s.conn.WriteToUDP([]byte(msgStr), addr)
	if err != nil {
		log.Printf("response sending error: %v", err)
	}
}

// extractSIPURI extracts SIP URI from header
func extractSIPURI(header string) string {
	start := strings.Index(header, "sip:")
	if start == -1 {
		return ""
	}

	end := strings.Index(header[start:], ">")
	if end == -1 {
		return header[start:]
	}

	return header[start : start+end]
}
