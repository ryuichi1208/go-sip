package sip

import (
	"strings"
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage()
	if msg == nil {
		t.Fatal("NewMessage() returned nil")
	}
	if msg.Headers == nil {
		t.Error("Headers map not initialized")
	}
}

func TestParseMessage(t *testing.T) {
	// Test valid SIP message
	validMsg := "REGISTER sip:test@example.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK123\r\n" +
		"From: <sip:test@example.com>;tag=123\r\n" +
		"To: <sip:test@example.com>\r\n" +
		"Call-ID: test123\r\n" +
		"CSeq: 1 REGISTER\r\n" +
		"Contact: <sip:test@127.0.0.1>\r\n" +
		"Content-Length: 0\r\n\r\n"

	msg, err := ParseMessage(validMsg)
	if err != nil {
		t.Fatalf("Failed to parse valid message: %v", err)
	}

	if msg.StartLine != "REGISTER sip:test@example.com SIP/2.0" {
		t.Errorf("Wrong start line: %s", msg.StartLine)
	}

	if msg.Headers["Call-ID"] != "test123" {
		t.Errorf("Wrong Call-ID: %s", msg.Headers["Call-ID"])
	}

	if msg.Headers["From"] != "<sip:test@example.com>;tag=123" {
		t.Errorf("Wrong From: %s", msg.Headers["From"])
	}

	if msg.Body != "" {
		t.Errorf("Body should be empty, got: %s", msg.Body)
	}

	// Test message with body
	msgWithBody := "INVITE sip:test@example.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK123\r\n" +
		"From: <sip:caller@example.com>;tag=123\r\n" +
		"To: <sip:test@example.com>\r\n" +
		"Call-ID: test123\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length: 19\r\n\r\n" +
		"v=0\r\n" +
		"o=test 123 456"

	msg, err = ParseMessage(msgWithBody)
	if err != nil {
		t.Fatalf("Failed to parse message with body: %v", err)
	}

	if msg.Body != "v=0\r\no=test 123 456" {
		t.Errorf("Wrong body: %s", msg.Body)
	}

	// Test empty message
	_, err = ParseMessage("")
	if err == nil {
		t.Error("Expected error for empty message, but got nil")
	}

	// Test malformed message
	_, err = ParseMessage("NOT A SIP MESSAGE")
	if err == nil {
		t.Error("Expected error for malformed message, but got nil")
	}
}

func TestString(t *testing.T) {
	msg := NewMessage()
	msg.StartLine = "SIP/2.0 200 OK"
	msg.Headers["Via"] = "SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK123"
	msg.Headers["From"] = "<sip:test@example.com>;tag=123"
	msg.Headers["To"] = "<sip:test@example.com>"
	msg.Headers["Call-ID"] = "test123"
	msg.Headers["CSeq"] = "1 REGISTER"
	msg.Headers["Content-Length"] = "0"

	msgStr := msg.String()

	// Check that all headers are present
	if !strings.Contains(msgStr, "SIP/2.0 200 OK\r\n") {
		t.Error("Missing start line in message string")
	}
	if !strings.Contains(msgStr, "Via: SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK123\r\n") {
		t.Error("Missing Via header in message string")
	}
	if !strings.Contains(msgStr, "From: <sip:test@example.com>;tag=123\r\n") {
		t.Error("Missing From header in message string")
	}
	if !strings.Contains(msgStr, "To: <sip:test@example.com>\r\n") {
		t.Error("Missing To header in message string")
	}
	if !strings.Contains(msgStr, "Call-ID: test123\r\n") {
		t.Error("Missing Call-ID header in message string")
	}
	if !strings.Contains(msgStr, "CSeq: 1 REGISTER\r\n") {
		t.Error("Missing CSeq header in message string")
	}
	if !strings.Contains(msgStr, "Content-Length: 0\r\n") {
		t.Error("Missing Content-Length header in message string")
	}

	// Check body and end of message
	if !strings.HasSuffix(msgStr, "\r\n\r\n") {
		t.Error("Message does not end with \\r\\n\\r\\n")
	}

	// Test with body
	msg.Body = "test body"
	msg.Headers["Content-Length"] = "9"
	msgStr = msg.String()

	if !strings.HasSuffix(msgStr, "\r\n\r\ntest body") {
		t.Error("Message body not correctly appended")
	}
}

func TestNewResponse(t *testing.T) {
	request := NewMessage()
	request.StartLine = "REGISTER sip:test@example.com SIP/2.0"
	request.Headers["Via"] = "SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK123"
	request.Headers["From"] = "<sip:test@example.com>;tag=123"
	request.Headers["To"] = "<sip:test@example.com>"
	request.Headers["Call-ID"] = "test123"
	request.Headers["CSeq"] = "1 REGISTER"

	response := NewResponse("200", "OK", request)

	if response.StartLine != "SIP/2.0 200 OK" {
		t.Errorf("Wrong start line: %s", response.StartLine)
	}

	if response.Headers["Via"] != "SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK123" {
		t.Error("Via header not copied from request")
	}

	if response.Headers["From"] != "<sip:test@example.com>;tag=123" {
		t.Error("From header not copied from request")
	}

	if response.Headers["To"] != "<sip:test@example.com>" {
		t.Error("To header not copied from request")
	}

	if response.Headers["Call-ID"] != "test123" {
		t.Error("Call-ID header not copied from request")
	}

	if response.Headers["CSeq"] != "1 REGISTER" {
		t.Error("CSeq header not copied from request")
	}

	if response.Headers["Server"] != "Go-SIP-Server" {
		t.Error("Server header not set")
	}

	if response.Headers["Content-Length"] != "0" {
		t.Error("Content-Length header not set")
	}

	if _, exists := response.Headers["Date"]; !exists {
		t.Error("Date header not set")
	}
}
