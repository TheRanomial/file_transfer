package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request){

	t,err:=template.ParseFiles("index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w,nil)
}

func fileUploader(w http.ResponseWriter, r *http.Request){
	
	if r.Method=="GET"{
		t,err:=template.ParseFiles("index.html")
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}
		t.Execute(w,nil)

	} else if r.Method=="POST"{

		err:=r.ParseMultipartForm(32<<20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file,header,err:=r.FormFile("file")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err=os.MkdirAll("./uploads",os.ModePerm)
		if err!=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dst,err:=os.Create("./uploads/"+header.Filename)
		if err!=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _,err:=io.Copy(dst,file); err!=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename)
		fmt.Println("file uploaded successfully")

	} else{
		http.Error(w, "Method not supported", http.StatusInternalServerError)
	}
}

func main(){

	http.HandleFunc("/",indexHandler)
	http.HandleFunc("/upload",fileUploader)

	err:=http.ListenAndServe(":8080",nil)

	if err!=nil{
		fmt.Println("there is an error")
		return
	}
}