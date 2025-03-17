package main

import (
	"io"
	"log"
	"os"
)

type PathTransFormFunc func(string) string

type StoreOpts struct {
	PathTransFormFunc PathTransFormFunc
}

var DefaultPathTransformFunc = func(key  string)string{
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader)error{
	pathName:= s.PathTransFormFunc(key)
	if err := os.MkdirAll(pathName,os.ModePerm);err!=nil{
		return err
	}
	filename:= "somefilename"
	pathAndFilename:= pathName + "/" + filename
	f,err := os.Create(pathAndFilename)
	if err!=nil{
		return err
	}

	n,err := io.Copy(f,r)
	if err!=nil{
		return err
	}
	log.Printf("written (%d) bytes to disk: %s",n,pathAndFilename)
	return nil
}
