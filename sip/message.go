package sip

import (
	"fmt"
	"strings"
	"time"
)

// Message represents a SIP message structure
type Message struct {
	StartLine string
	Headers   map[string]string
	Body      string
}

// NewMessage creates a new SIP message
func NewMessage() *Message {
	return &Message{
		Headers: make(map[string]string),
	}
}

// ParseMessage parses a SIP message from a string
func ParseMessage(data string) (*Message, error) {
	// Check for empty message
	if data == "" {
		return nil, fmt.Errorf("empty message")
	}

	// Split into lines, ensuring we have \r\n line endings
	data = strings.Replace(data, "\n", "\r\n", -1)
	data = strings.Replace(data, "\r\r\n", "\r\n", -1)
	lines := strings.Split(data, "\r\n")

	// Must have at least a start line
	if len(lines) < 1 {
		return nil, fmt.Errorf("empty message")
	}

	// Create new message
	msg := NewMessage()
	msg.StartLine = lines[0]

	// Validate start line
	if !strings.Contains(msg.StartLine, "SIP/2.0") {
		return nil, fmt.Errorf("invalid SIP message: missing SIP/2.0 in start line")
	}

	// Find boundary between headers and body
	bodyStart := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			bodyStart = i + 1
			break
		}

		parts := strings.SplitN(lines[i], ":", 2)
		if len(parts) != 2 {
			continue // Skip invalid header lines
		}

		headerName := strings.TrimSpace(parts[0])
		headerValue := strings.TrimSpace(parts[1])
		msg.Headers[headerName] = headerValue
	}

	// Get body if present
	if bodyStart != -1 && bodyStart < len(lines) {
		msg.Body = strings.Join(lines[bodyStart:], "\r\n")
	}

	return msg, nil
}

// String converts the message to a string representation
func (m *Message) String() string {
	var sb strings.Builder
	sb.WriteString(m.StartLine + "\r\n")

	for name, value := range m.Headers {
		sb.WriteString(name + ": " + value + "\r\n")
	}

	sb.WriteString("\r\n")
	if m.Body != "" {
		sb.WriteString(m.Body)
	}

	return sb.String()
}

// NewResponse generates a SIP response message
func NewResponse(statusCode string, statusText string, request *Message) *Message {
	resp := NewMessage()
	resp.StartLine = fmt.Sprintf("SIP/2.0 %s %s", statusCode, statusText)

	// Copy headers from request
	headersToCopy := []string{"Call-ID", "From", "To", "CSeq", "Via"}
	for _, header := range headersToCopy {
		if val, ok := request.Headers[header]; ok {
			resp.Headers[header] = val
		}
	}

	// Add server info and timestamp
	resp.Headers["Server"] = "Go-SIP-Server"
	resp.Headers["Date"] = time.Now().Format(time.RFC1123)
	resp.Headers["Content-Length"] = "0"

	return resp
}
