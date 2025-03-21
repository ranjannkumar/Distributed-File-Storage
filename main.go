package main

import (
	"log"

	"github.com/ranjannkumar/distributedFileStorage/p2p"
)

func makeServer(ListenAddr string,nodes ...string)*FileServer{
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: ListenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		//TODO ; onPeer funnc
	}

	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:        ListenAddr + "_network",
		PathTransFormFunc:  CASPathTransformFunc,
		Transport:          tcpTransport,
		BootstrapNodes:      nodes,
	}

	s:= NewFileServer(fileServerOpts)

	tcpTransport.Onpeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000","")
	s2 := makeServer(":4000",":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()
}

//4:12:23