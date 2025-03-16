package p2p

import (
	"fmt"
	"net"
	"sync"
)

//TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct{
	//conn is the underlying connection of the peer
	conn net.Conn
	//if we dial and retrieve a conn =>outbound==true
	//if we accept and rettrive a conn =>outbound ==false

	outbound bool
}

func NewTCPPeer(conn net.Conn,outbound bool)*TCPPeer{
	return &TCPPeer{
		conn: conn,
		outbound: outbound,
	}
}

//Close implements the peer interface
func (p *TCPPeer)Close()error{
	return p.conn.Close()
}

type TCPTransportOpts struct{
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	Onpeer        func(Peer)error
}
  
type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	rpcch         chan RPC

	mu     sync.RWMutex
}


func NewTCPTransport(opts TCPTransportOpts)*TCPTransport{
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

//consume implements the Transport interface,which will return read-only channel
// for reading the incoming messages recieved from another peer in the network
func (t *TCPTransport)Consume()<-chan RPC{
	return t.rpcch
}

func (t *TCPTransport)ListenAndAccept()error{
	var err error
	t.listener,err = net.Listen("tcp",t.ListenAddr)
	if err!=nil{
		return err
	}
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport)startAcceptLoop(){
	for{
		conn,err := t.listener.Accept()
		if err!=nil{
			fmt.Printf("TCP accept error: %s\n",err)
		}
	  fmt.Printf("new incoming connection %+v\n",conn)	
	   go t.handleConn(conn)
	}
}


func (t *TCPTransport) handleConn(conn net.Conn){
	  var err error
		defer func ()  {
			fmt.Printf("dropping peer connetion: %s",err)
			conn.Close()
		}()

		peer:=NewTCPPeer(conn,true)

		if err:= t.HandshakeFunc(peer);err!=nil{
			return 
		}

		if t.Onpeer !=nil{
			if err = t.Onpeer(peer);err!=nil{
				return
			}
		}

		//Read loop
		rpc := RPC{}
		for {
			
			 err = t.Decoder.Decode(conn,&rpc)
			 if err!=nil{
				return
			}
			rpc.From = conn.RemoteAddr()
			t.rpcch <- rpc
		}
}   