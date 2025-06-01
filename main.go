package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		EncKey:             newEncryptionKey(),
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
	s2 := makeServer(":7000","")
	s3 := makeServer(":5000",":3000",":7000")

	go func() {
		log.Fatal(s1.Start())
	}()
		time.Sleep(time.Millisecond * 500)
		go func() {
		log.Fatal(s2.Start())
	}()


	time.Sleep(2*time.Second)

	go s3.Start()
	time.Sleep(2*time.Second)
  
 for i:=0;i<20;i++{
		key:=fmt.Sprintf("picture_%d.png",i)
		data := bytes.NewReader([]byte("my big data file here!")) 
		s3.Store(key,data)

		if err := s3.store.Delete(key);err!=nil{
			log.Fatal(err)
		}
		
		r,err := s3.Get(key)
		if err !=nil{
			log.Fatal(err)
		}

		b,err := ioutil.ReadAll(r)
		if err !=nil{
			log.Fatal(err)
		}

		fmt.Println(string(b))
  }
}




