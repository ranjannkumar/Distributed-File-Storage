package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "ggnetwork"

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
	//Root is the folder name of the root ,containing all the folders/files of the system
	Root              string
	PathTransFormFunc PathTransFormFunc
}

var DefaultPathTransformFunc = func(key  string)PathKey{
	return PathKey{
		PathName: key,
		Filename: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransFormFunc==nil{
		opts.PathTransFormFunc=DefaultPathTransformFunc
	}
	if len(opts.Root)==0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string)bool{
	pathKey := s.PathTransFormFunc(key)
	rootPathWithColon := strings.Split(s.Root, ":")
	rootDir := rootPathWithColon[1]
	fullPathWithRoot := fmt.Sprintf("%s/%s",rootDir,pathKey.FullPath())

	_,err := os.Stat(fullPathWithRoot)
	return !errors.Is(err,os.ErrNotExist)
}


func (s *Store)Clear()error{
	return os.RemoveAll(s.Root)
}

func (s *Store)Delete(key string)error{
	pathKey := s.PathTransFormFunc(key)
	rootPathWithColon := strings.Split(s.Root, ":")
	rootDir := rootPathWithColon[1]

	defer func(){
		log.Printf("deleted [%s] from disk",pathKey.Filename)
	}()
	firstPathNameWithRoot := fmt.Sprintf("%s/%s",rootDir,pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store)Write(key string,r io.Reader)(int64,error){
	return s.writeStream(key,r)
}

func (s *Store) Read(key string)(int64, io.Reader, error){
	return s.readStream(key)
}

func (s *Store)readStream(key string)(int64, io.ReadCloser,error){
	pathKey := s.PathTransFormFunc(key)
	rootPathWithColon := strings.Split(s.Root, ":")
	rootDir := rootPathWithColon[1]
	fullPathWithRoot := fmt.Sprintf("%s/%s",rootDir,pathKey.FullPath())

	file,err := os.Open(fullPathWithRoot)
	if err !=nil{
		return 0,nil,err
	}

	fi,err := file.Stat()
	if err!=nil{
		return 0,nil,err
	}
	return fi.Size(),file,nil
}

func (s *Store) writeStream(key string, r io.Reader)(int64,error){
			pathKey:= s.PathTransFormFunc(key)
			rootPathWithColon := strings.Split(s.Root, ":")
			rootDir := rootPathWithColon[1]
			pathNameWithRoot := fmt.Sprintf("%s/%s",rootDir,pathKey.PathName)
			if err := os.MkdirAll(pathNameWithRoot,os.ModePerm);err!=nil{
				return 0,err
			}

			fullPathWithRoot:= fmt.Sprintf("%s/%s",rootDir,pathKey.FullPath())

			f,err := os.Create(fullPathWithRoot)
			if err!=nil{
				return 0,err
			}
			defer f.Close()

			n,err := io.Copy(f,r)
			if err!=nil{
				return 0,err
			}
			// log.Printf("written (%d) bytes to disk: %s",n,fullPathWithRoot)
			return n,nil
}
