package main

import (
	"log"
	"time"

	"github.com/ranjannkumar/distributedFileStorage/p2p"
)


func main() {

	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: ":4000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		//TODO ; onPeer funnc
	}

	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:        "4000_network",
		PathTransFormFunc:  CASPathTransformFunc,
		Transport:          tcpTransport,
	}

	s := NewFileServer(fileServerOpts)

	go func(){
		time.Sleep(time.Second*3)
		s.Stop()
	}()

	if err:= s.Start();err!=nil{
		log.Fatal(err)
	}

}

//3"04:34