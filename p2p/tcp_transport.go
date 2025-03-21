package p2p

import (
	"errors"
	"fmt"
	"log"
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

//RemoteAddr implements the Peer ingterface and will return the
//remote address of its underlying connection.
func (p *TCPPeer)RemoteAddr()net.Addr{
	return p.conn.RemoteAddr()
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
  
//A TCPTransport is responsible for managing the overall TCP transport layer, including listening for connections, accepting connections, and handling peer communication.
//the TCPTransport creates the connections, and the TCPPeer represents those individual connections.
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

//Close implements the Transport interface
func (t *TCPTransport)Close()error{
	return t.listener.Close()
}

//Dial implements the Transport ineterface
func(t *TCPTransport)Dial(addr string)error{
	conn,err:= net.Dial("tcp",addr)
	if err!=nil{
		return err
	}

	go t.handleConn(conn,true)

	return nil
}

func (t *TCPTransport)ListenAndAccept()error{
	var err error
	t.listener,err = net.Listen("tcp",t.ListenAddr)
	if err!=nil{
		return err
	}
	go t.startAcceptLoop()
	log.Printf("TCP transport listening on port: %s\n",t.ListenAddr)
	return nil
}

func (t *TCPTransport)startAcceptLoop(){
	for{
		conn,err := t.listener.Accept()
		if errors.Is(err,net.ErrClosed){
			return
		}
		if err!=nil{
			fmt.Printf("TCP accept error: %s\n",err)
		}
	  fmt.Printf("new incoming connection %+v\n",conn)	
	   go t.handleConn(conn,false)
	}
}


func (t *TCPTransport) handleConn(conn net.Conn,outbound bool){
	  var err error
		defer func ()  {
			fmt.Printf("dropping peer connetion: %s",err)
			conn.Close()
		}()

		peer:=NewTCPPeer(conn,outbound)

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
			//This line sends the decoded rpc message into the rpcch channel.
			t.rpcch <- rpc
		}
}   