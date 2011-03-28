package pop3


import (
	"net"
	"os"
	"bufio"
	"strings"
)

const (
	//Carriage Return + Line Feed
	//CRLF is appended at the end of each commands
	CRLF = "\r\n" 
)

type Client struct {
	conn       net.Conn
	stream     *bufio.ReadWriter
	ServerName string
	Greeting  string
} 

//Returns a new Client connected to a POP3 server at addr.
//The format of addr is "ip:port"
func Dial(addr string) (client *Client, err os.Error) {
	conn, err := net.Dial("tcp", "", addr)
	if err != nil {
		return nil, err
	}
	host := addr[:strings.Index(addr, ":")]
	return NewClient(conn, host)

}

//NewClient returns a new Client using an existing connection
//name is used as the Servername
func NewClient(conn net.Conn, name string) (*Client, os.Error) {
	client := new(Client)

	//Create a new ReadWriter and store it
	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	client.stream = stream

	//Download the greeting from the POP3 server
	msg, err := client.ReadMessage(false)

	if err != nil {
		return nil, err
	}
	
	client.ServerName = name
	client.Greeting = msg
	
	return client, nil
}

//WriteMessage sends the message to the POP3 server
func (client *Client) WriteMessage(message string) os.Error {
	if client == nil {
		return os.NewError("Connection hasn't been established")
	}
	
	tmp := message + CRLF
	_, err1 := client.stream.WriteString(tmp)
	client.stream.Flush()

	return err1
}

//ReadMessage reads a single or multiline response from the POP3 server
//It doesnt finish, until it has received a message
func (client *Client) ReadMessage(multiLine bool) (string, os.Error) {

	//Check, whether the client connection has already been 
	if client == nil {
		return "", os.NewError("Connection hasn't been established")
	}
	
	//Get first line of the response
	msg, err := client.stream.ReadString('\n')

	if err != nil {
		return "", err
	}

	//Check, whether the response starts with "+OK" or "-ERR", otherwise return an error
	if strings.HasPrefix(msg, "+OK") {
		msg = msg[4:]
		
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
	
	} else if strings.HasPrefix(msg, "-ERR") {
		return "", os.NewError(msg[5:])
	} else {
		return "", os.NewError("Unkown Response received")
	}
	
	return msg, nil
}
