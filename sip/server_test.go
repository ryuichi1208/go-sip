package sip

import (
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

// MockConn implements a mock net.UDPConn for testing
type MockConn struct {
	receivedData []byte
	sentData     []byte
	addr         *net.UDPAddr
	mutex        sync.Mutex
}

// ReadFromUDP is a mock implementation that returns predefined data
func (m *MockConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	copy(b, m.receivedData)
	return len(m.receivedData), m.addr, nil
}

// WriteToUDP records data that would be sent
func (m *MockConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sentData = make([]byte, len(b))
	copy(m.sentData, b)
	return len(b), nil
}

// Close is a mock implementation
func (m *MockConn) Close() error {
	return nil
}

// SetReceivedData allows setting test data to be "received"
func (m *MockConn) SetReceivedData(data []byte, addr *net.UDPAddr) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.receivedData = make([]byte, len(data))
	copy(m.receivedData, data)
	m.addr = addr
}

// GetSentData returns the data that was "sent"
func (m *MockConn) GetSentData() []byte {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.sentData
}

func setupTestServer(t *testing.T) *Server {
	// Create a server
	server := NewServer("5060")

	// Create a mock connection
	mockConn := &MockConn{}
	server.conn = mockConn

	return server
}

func TestHandleRegister(t *testing.T) {
	server := setupTestServer(t)
	mockConn, ok := server.conn.(*MockConn)
	if !ok {
		t.Fatal("Failed to cast server.conn to *MockConn")
	}

	// Create a test client address
	clientAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 12345,
	}

	// Create a REGISTER message
	registerMsg := NewMessage()
	registerMsg.StartLine = "REGISTER sip:example.com SIP/2.0"
	registerMsg.Headers["Via"] = "SIP/2.0/UDP 127.0.0.1:12345;branch=z9hG4bK123"
	registerMsg.Headers["From"] = "<sip:alice@example.com>;tag=123"
	registerMsg.Headers["To"] = "<sip:alice@example.com>"
	registerMsg.Headers["Call-ID"] = "register-test-123"
	registerMsg.Headers["CSeq"] = "1 REGISTER"
	registerMsg.Headers["Contact"] = "<sip:alice@127.0.0.1:12345>"
	registerMsg.Headers["Content-Length"] = "0"

	// Handle the REGISTER message
	server.handleRegister(clientAddr, registerMsg)

	// Check that user was registered
	uri := "sip:alice@example.com"
	addr, exists := server.registrar[uri]
	if !exists {
		t.Fatalf("User %s not registered", uri)
	}
	if addr != clientAddr.String() {
		t.Errorf("Wrong address registered: got %s, want %s", addr, clientAddr.String())
	}

	// Check the response
	sentData := mockConn.GetSentData()
	if sentData == nil {
		t.Fatal("No response was sent")
	}

	responseStr := string(sentData)
	if !strings.Contains(responseStr, "SIP/2.0 200 OK") {
		t.Error("Response is not 200 OK")
	}
	if !strings.Contains(responseStr, "Call-ID: register-test-123") {
		t.Error("Response has wrong Call-ID")
	}
}

func TestHandleInvite(t *testing.T) {
	server := setupTestServer(t)
	mockConn, ok := server.conn.(*MockConn)
	if !ok {
		t.Fatal("Failed to cast server.conn to *MockConn")
	}

	// Create a test client address
	clientAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 12345,
	}

	// Create an INVITE message
	inviteMsg := NewMessage()
	inviteMsg.StartLine = "INVITE sip:bob@example.com SIP/2.0"
	inviteMsg.Headers["Via"] = "SIP/2.0/UDP 127.0.0.1:12345;branch=z9hG4bK123"
	inviteMsg.Headers["From"] = "<sip:alice@example.com>;tag=123"
	inviteMsg.Headers["To"] = "<sip:bob@example.com>"
	inviteMsg.Headers["Call-ID"] = "invite-test-123"
	inviteMsg.Headers["CSeq"] = "1 INVITE"
	inviteMsg.Headers["Contact"] = "<sip:alice@127.0.0.1:12345>"
	inviteMsg.Headers["Content-Length"] = "0"

	// Handle the INVITE message
	server.handleInvite(clientAddr, inviteMsg)

	// Check that call was created
	callID := "invite-test-123"
	status, exists := server.calls[callID]
	if !exists {
		t.Fatalf("Call %s not created", callID)
	}
	if status != "connected" {
		t.Errorf("Wrong call status: got %s, want connected", status)
	}

	// Check the response - should be a 200 OK eventually
	// We might need to wait for all responses (100, 180, 200)
	time.Sleep(100 * time.Millisecond)

	sentData := mockConn.GetSentData()
	if sentData == nil {
		t.Fatal("No response was sent")
	}

	responseStr := string(sentData)
	if !strings.Contains(responseStr, "SIP/2.0 200 OK") {
		t.Error("Final response is not 200 OK")
	}
	if !strings.Contains(responseStr, "Call-ID: invite-test-123") {
		t.Error("Response has wrong Call-ID")
	}
}

func TestHandleBye(t *testing.T) {
	server := setupTestServer(t)
	mockConn, ok := server.conn.(*MockConn)
	if !ok {
		t.Fatal("Failed to cast server.conn to *MockConn")
	}

	// Create a test client address
	clientAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 12345,
	}

	// Add a call to the server
	callID := "bye-test-123"
	server.calls[callID] = "connected"

	// Create a BYE message
	byeMsg := NewMessage()
	byeMsg.StartLine = "BYE sip:bob@example.com SIP/2.0"
	byeMsg.Headers["Via"] = "SIP/2.0/UDP 127.0.0.1:12345;branch=z9hG4bK123"
	byeMsg.Headers["From"] = "<sip:alice@example.com>;tag=123"
	byeMsg.Headers["To"] = "<sip:bob@example.com>;tag=456"
	byeMsg.Headers["Call-ID"] = callID
	byeMsg.Headers["CSeq"] = "2 BYE"
	byeMsg.Headers["Content-Length"] = "0"

	// Handle the BYE message
	server.handleBye(clientAddr, byeMsg)

	// Check that call was removed
	_, exists := server.calls[callID]
	if exists {
		t.Errorf("Call %s not removed", callID)
	}

	// Check the response
	sentData := mockConn.GetSentData()
	if sentData == nil {
		t.Fatal("No response was sent")
	}

	responseStr := string(sentData)
	if !strings.Contains(responseStr, "SIP/2.0 200 OK") {
		t.Error("Response is not 200 OK")
	}
	if !strings.Contains(responseStr, "Call-ID: bye-test-123") {
		t.Error("Response has wrong Call-ID")
	}
}

func TestExtractSIPURI(t *testing.T) {
	testCases := []struct {
		header   string
		expected string
	}{
		{
			header:   "<sip:alice@example.com>",
			expected: "sip:alice@example.com",
		},
		{
			header:   "<sip:bob@example.com>;tag=123",
			expected: "sip:bob@example.com",
		},
		{
			header:   "\"Bob\" <sip:bob@example.com>",
			expected: "sip:bob@example.com",
		},
		{
			header:   "<tel:+12345678>",
			expected: "",
		},
		{
			header:   "sip:alice@example.com",
			expected: "sip:alice@example.com",
		},
	}

	for i, tc := range testCases {
		result := extractSIPURI(tc.header)
		if result != tc.expected {
			t.Errorf("Test case %d: expected %q, got %q", i, tc.expected, result)
		}
	}
}
