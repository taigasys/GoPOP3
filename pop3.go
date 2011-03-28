package pop3


import (
	"net"
	"os"
	"bufio"
	"strings"
)

type Client struct {
	conn       net.Conn
	stream     *bufio.ReadWriter
	ServerName string
	Greetings  string
} 

func Dial(addr string) (client *Client, err os.Error) {
	conn, err := net.Dial("tcp", "", addr)
	if err != nil {
		return nil, err
	}
	host := addr[:strings.Index(addr, ":")]
	return NewClient(conn, host)

}

func NewClient(conn net.Conn, name string) (client *Client, err os.Error) {
	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	line, err := stream.ReadString('\n')

	if err != nil {
		return nil, err
	}

	client.Greetings = line

	return client, nil
}

func (client *Client) ReadMessage(multiLine bool) (string, os.Error) {
	if client == nil {
		return "", os.NewError("Connection hasn't been established")
	}

	message, err := client.stream.ReadString('\n')

	if err != nil {
		return "", err
	}

	if multiLine {

		for true {
			line, err1 := client.stream.ReadString('\n')

			if err1 != nil {
				return "", err1
			}
			
			if line != ".\n" {
				break;
			}
		}
	}
	
	return message, nil
}
