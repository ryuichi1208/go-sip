package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// Sample SIP client
func main() {
	// Parse command line arguments
	serverAddr := flag.String("server", "127.0.0.1:5060", "SIP server address")
	username := flag.String("user", "user1", "SIP username")
	flag.Parse()

	// Establish UDP connection
	addr, err := net.ResolveUDPAddr("udp", *serverAddr)
	if err != nil {
		log.Fatalf("Address resolution error: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("UDP connection error: %v", err)
	}
	defer conn.Close()

	fmt.Printf("SIP client started: %s\n", *username)
	fmt.Println("Commands: register, invite, bye, exit")

	// Goroutine for receiving responses
	go func() {
		buffer := make([]byte, 65535)
		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Response reading error: %v", err)
				continue
			}
			fmt.Printf("\nReceived response:\n%s\n", string(buffer[:n]))
		}
	}()

	// Command input loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		cmd := scanner.Text()
		if cmd == "exit" {
			break
		}

		switch cmd {
		case "register":
			sendRegister(conn, *username)
		case "invite":
			fmt.Print("Username to call: ")
			if !scanner.Scan() {
				break
			}
			callee := scanner.Text()
			sendInvite(conn, *username, callee)
		case "bye":
			sendBye(conn, *username)
		default:
			fmt.Println("Unknown command. Enter register, invite, bye, or exit.")
		}
	}
}

func sendRegister(conn *net.UDPConn, username string) {
	callID := generateCallID()

	msg := fmt.Sprintf("REGISTER sip:%s@localhost SIP/2.0\r\n"+
		"Via: SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK%s\r\n"+
		"From: <sip:%s@localhost>;tag=%s\r\n"+
		"To: <sip:%s@localhost>\r\n"+
		"Call-ID: %s\r\n"+
		"CSeq: 1 REGISTER\r\n"+
		"Contact: <sip:%s@127.0.0.1>\r\n"+
		"Max-Forwards: 70\r\n"+
		"User-Agent: Go-SIP-Client\r\n"+
		"Expires: 3600\r\n"+
		"Content-Length: 0\r\n\r\n",
		username, generateBranch(), username, generateTag(), username, callID, username)

	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("Message sending error: %v", err)
		return
	}

	fmt.Println("REGISTER sent")
}

func sendInvite(conn *net.UDPConn, caller, callee string) {
	callID := generateCallID()

	sdp := fmt.Sprintf("v=0\r\n"+
		"o=%s 123456 654321 IN IP4 127.0.0.1\r\n"+
		"s=SIP Call\r\n"+
		"c=IN IP4 127.0.0.1\r\n"+
		"t=0 0\r\n"+
		"m=audio 49170 RTP/AVP 0\r\n"+
		"a=rtpmap:0 PCMU/8000\r\n", caller)

	msg := fmt.Sprintf("INVITE sip:%s@localhost SIP/2.0\r\n"+
		"Via: SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK%s\r\n"+
		"From: <sip:%s@localhost>;tag=%s\r\n"+
		"To: <sip:%s@localhost>\r\n"+
		"Call-ID: %s\r\n"+
		"CSeq: 1 INVITE\r\n"+
		"Contact: <sip:%s@127.0.0.1>\r\n"+
		"Content-Type: application/sdp\r\n"+
		"Content-Length: %d\r\n\r\n%s",
		callee, generateBranch(), caller, generateTag(), callee, callID, caller, len(sdp), sdp)

	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("Message sending error: %v", err)
		return
	}

	fmt.Println("INVITE sent")
}

func sendBye(conn *net.UDPConn, username string) {
	callID := generateCallID()

	msg := fmt.Sprintf("BYE sip:server@localhost SIP/2.0\r\n"+
		"Via: SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK%s\r\n"+
		"From: <sip:%s@localhost>;tag=%s\r\n"+
		"To: <sip:server@localhost>;tag=as6f4bc61\r\n"+
		"Call-ID: %s\r\n"+
		"CSeq: 1 BYE\r\n"+
		"Content-Length: 0\r\n\r\n",
		generateBranch(), username, generateTag(), callID)

	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("Message sending error: %v", err)
		return
	}

	fmt.Println("BYE sent")
}

func generateCallID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateTag() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%100000)
}

func generateBranch() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
}
