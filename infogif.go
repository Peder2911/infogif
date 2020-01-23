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
	"time"
	)


type Gif struct {
	Title string
	Path string
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

	uploadedFile, _, err := request.FormFile("file")
	if err != nil {
		log.Print("Failed to acquire the uploaded content")
		writer.WriteHeader(http.StatusInternalServerError)
		//writeError(writer, err)
		return
	}
	defer uploadedFile.Close()
	// Read content
	body, err := ioutil.ReadAll(uploadedFile)
	if err != nil {
		log.Print("Failed to read the uploaded content") 
		writer.WriteHeader(http.StatusInternalServerError)
		//writeError(writer,err)
		return
	}

	// Gif content is hashed to assert identity.
	filehash := md5.Sum(body)
	filehashDigest := hex.EncodeToString(filehash[:])

	// Check if file already exists
	dstPath := path.Join("static/gifs", filehashDigest)
	if !fileExists(dstPath){
		// Write the data to a new file 
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
		} 
		defer cleanup(dstPath)
	} else {
	}
	fmt.Fprintf(writer, filehashDigest)
}


func cleanup(path string){
	deleteFile:= func(){
		log.Print(fmt.Sprintf("Deleting %s",path))
		err := os.Remove(path)
		if err != nil {
			log.Print(err)
		}
	}
	// Deletes a file after X minutes, since this is not a file hosting service.
	time.AfterFunc(5*time.Minute, deleteFile) 
}


func main() {
	http.HandleFunc("/",baseHandler)
	http.HandleFunc("/upload/",uploadHandler)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":9090",nil))
}
