package pop3


import (
	"net"
	"os"
	"bufio"
	"strings"
)

const (
	CRLF = "\r\n" 
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

func NewClient(conn net.Conn, name string) (*Client, os.Error) {
	client := new(Client)

	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	client.stream = stream

	msg, err := client.ReadMessage(false)

	if err != nil {
		return nil, err
	}
	
	client.Greetings = msg
	return client, nil
}

func (client *Client) WriteMessage(message string) os.Error {
	if client == nil {
		return os.NewError("Connection hasn't been established")
	}
	
	tmp := message + CRLF
	_, err1 := client.stream.WriteString(tmp)
	client.stream.Flush()

	return err1
}

func (client *Client) ReadMessage(multiLine bool) (string, os.Error) {
	if client == nil {
		return "", os.NewError("Connection hasn't been established")
	}
	
	msg, err := client.stream.ReadString('\n')

	if err != nil {
		return "", err
	}

	//Check, whether the response starts with "+OK" or "-ERR", otherwise return an error
	if strings.HasPrefix(msg, "+OK") {
		msg = msg[4:]
	} else if strings.HasPrefix(msg, "-ERR") {
		return "", os.NewError(msg[5:])
	} else {
		return "", os.NewError("Unkown Response received")
	}

	if multiLine {

		for true {
			line, err1 := client.stream.ReadString('\n')

			if err1 != nil {
				return "", err1
			}
			
			if line == "." + CRLF {
				break;
			}
			
			msg += line
		}
	}
	
	return msg[4:], nil
}
