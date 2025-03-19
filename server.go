package main

import (
	"fmt"
	"log"

	"github.com/ranjannkumar/distributedFileStorage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransFormFunc PathTransFormFunc
	Transport         p2p.Transport
	BootstrapNodes     []string
}

type FileServer struct {
	FileServerOpts
	store *Store
	quitch chan  struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransFormFunc: opts.PathTransFormFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
	}
}

func (s *FileServer) Stop(){
	close(s.quitch)
}

func (s *FileServer) loop(){
	defer func(){
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()

	for{
		select{
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}
}

func(s *FileServer)bootstrapNetwork()error{
	for _,addr := range s.BootstrapNodes{
		go func (addr string)  {
			fmt.Println("attempting to connecct with remote: ",addr)

				if err := s.Transport.Dial(addr);err!=nil{
					log.Println("dial error: ",err)
				}
		}(addr)

	}
	return nil
}

func (s *FileServer)Start()error{
	if err:= s.Transport.ListenAndAccept();err!=nil{
		return err
	}
	s.bootstrapNetwork()
	s.loop()
	return nil
}