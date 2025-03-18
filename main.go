package main

import (
	"fmt"
	"log"

	"github.com/ranjannkumar/distributedFileStorage/p2p"
)
func Onpeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main() {

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":4000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		Onpeer: Onpeer,
	}

	tr:=p2p.NewTCPTransport(tcpOpts)

	go func ()  {
		for{
			msg := <-tr.Consume()
			fmt.Printf("%+v\n",msg)
		}
	}()
	
	if err:= tr.ListenAndAccept();err!=nil{
		log.Fatal(err)
	}
	select {}
}

//3"04:34