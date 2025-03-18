package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string)PathKey{
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen  := len(hashStr)/blocksize
	paths := make([]string,sliceLen)

	for i:=0;i<sliceLen;i++{
		from,to := i*blocksize,(i*blocksize)+blocksize
		paths[i]=hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths,"/"),
		Filename: hashStr,
	}
}

type PathTransFormFunc func(string) PathKey

type PathKey struct{
	PathName   string
	Filename   string
}

func(p PathKey)FirstPathName()string{
	paths:= strings.Split(p.PathName,"/")
	if len(paths)==0{
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath()string{
	return fmt.Sprintf("%s/%s",p.PathName,p.Filename)
}

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

func (s *Store) Has(key string)bool{
	pathKey := s.PathTransFormFunc(key)

	_,err := os.Stat(pathKey.FullPath())
	if err == fs.ErrNotExist{
		return false
	}
	return true
}

func (s *Store)Delete(key string)error{
	pathKey := s.PathTransFormFunc(key)

	defer func(){
		log.Printf("deleted [%s] from disk",pathKey.Filename)
	}()

	return os.RemoveAll(pathKey.FirstPathName())
}

func (s *Store) Read(key string)(io.Reader,error){
	f,err := s.readStream(key)
	if err!=nil{
		return nil,err
	}
	defer f.Close()
	
	buf := new(bytes.Buffer)
	_,err = io.Copy(buf,f)

	return buf,err
}

func (s *Store)readStream(key string)(io.ReadCloser,error){
	PathKey := s.PathTransFormFunc(key)
	return os.Open(PathKey.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader)error{
			pathKey:= s.PathTransFormFunc(key)
			if err := os.MkdirAll(pathKey.PathName,os.ModePerm);err!=nil{
				return err
			}

			fullPath:= pathKey.FullPath()

			f,err := os.Create(fullPath)
			if err!=nil{
				return err
			}
			defer f.Close()

			n,err := io.Copy(f,r)
			if err!=nil{
				return err
			}
			log.Printf("written (%d) bytes to disk: %s",n,fullPath)
			return nil
}
