package main

import (
	"net/http" 
	"html/template"
	"fmt"
	"log"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path"
	)


type Gif struct {
	Title string
	Path string
}

func baseHandler(writer http.ResponseWriter, request *http.Request) {
	gifid := request.URL.Path[len("/"):]
	gif := &Gif{Title: gifid, Path: "/static/gifs/"+gifid+".gif"}
	basetemplate, _ := template.ParseFiles("templates/base.html")
	basetemplate.Execute(writer, gif)
}

func uploadHandler(writer http.ResponseWriter, request *http.Request){
	/* TODO
	Implement filetype checking:

	/ image formats and magic numbers
	var magicTable = map[string]string{
		 "\xff\xd8\xff":      "image/jpeg",
		 "\x89PNG\r\n\x1a\n": "image/png",
		 "GIF87a":            "image/gif",
		 "GIF89a":            "image/gif",
	}

	// mimeFromIncipit returns the mime type of an image file from its first few
	// bytes or the empty string if the file does not look like a known file type
	func mimeFromIncipit(incipit []byte) string {
		 incipitStr := []byte(incipit)
		 for magic, mime := range magicTable {
			  if strings.HasPrefix(incipitStr, magic) {
					return mime
			  }
		 }

		 return ""
	}
	*/

	srcFile, _, err := request.FormFile("file")
	if err != nil {
		log.Print("Failed to acquire the uploaded content")
		writer.WriteHeader(http.StatusInternalServerError)
		//writeError(writer, err)
		return
	}

	defer srcFile.Close()
	//log.Print(info)

	/*
	size,err := getSize(srcFile)
	if err != nil {
		log.Print("Failed to get size of file") 
		writer.Writeheader(http.StatusInternalServerError)
		//writeError(writer,err)
		return
	}
	*/

	body, err := ioutil.ReadAll(srcFile)
	if err != nil {
		log.Print("Failed to read the uploaded content") 
		writer.WriteHeader(http.StatusInternalServerError)
		//writeError(writer,err)
		return
	}
	filehash := md5.Sum(body)
	filehashDigest := hex.EncodeToString(filehash[:])
	/*
	Make it return if it already has that gif, reducing disk writes.
	*/

	dstPath := path.Join("static/gifs", filehashDigest)
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		log.Print("Failed to open file!")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dstFile.Close()

	_ , err = dstFile.Write(body); if err != nil {
		log.Print("Failed to write content!")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(writer, filehashDigest)
	}
}

func main() {
	http.HandleFunc("/",baseHandler)
	http.HandleFunc("/upload/",uploadHandler)

	http.Handle("/static/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":9090",nil))
}
