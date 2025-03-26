package main

import (
	"bytes"
	"log"
	"time"

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

	time.Sleep(1*time.Second)
	go s2.Start()
	time.Sleep(1*time.Second)
	data := bytes.NewReader([]byte("my big data file here!")) 

	s2.StoreData("myprivatedata",data)
	select {}
}

//5:01:51 