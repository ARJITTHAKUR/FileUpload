package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

const port = "8080"
const form = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>File Upload</title>
</head>
<body>
<h1>File Upload</h1>
<div>
	<form action="/upload" method="post" enctype="multipart/form-data" id="uploadForm">
		<input type="file" name="file" multiple id="files">
		<button type="submit">Upload</button>
	</form>
</div>
<div>
	<h2>selected files</h2>
	<div id="list">

	</div>
</div>
<script>
	const form = document.querySelector("#uploadForm")
	const listElement = document.querySelector("#list")

	form.addEventListener("change",handleChange)

	function handleChange(e){
		const fileList = Array.from(e.target.files)
		fileList.forEach(file => {
			const li = document.createElement("div")
			li.textContent = file.name
			listElement.append(li)
		});
	}
</script>
</body>
</html>
`
const uploaded = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    <h1>Files uploaded !</h1>
    <button onclick="window.history.back()">Go back</button>
</body>
</html>
`

func main() {

	mux := http.DefaultServeMux

	server := http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index").Parse(form)
		if err != nil {
			w.Write([]byte("some error occurec"))
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		multipleFileReader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error()+":2", http.StatusBadRequest)
			return
		}

		for {
			part, err := multipleFileReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, err.Error()+":3", http.StatusBadRequest)
				return
			}
			path := filepath.Join("./uploaded/" + part.FileName())
			fmt.Println(path)
			file, err := os.Create(path)
			if err != nil {
				http.Error(w, err.Error()+":4", http.StatusBadRequest)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, part)

			if err != nil {
				http.Error(w, err.Error()+":5", http.StatusBadRequest)
				return
			}

		}
		t, err := template.New("uploaded").Parse(uploaded)
		if err != nil {
			w.Write([]byte("some error occurec"))
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	err := os.Mkdir("uploaded", 0777)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	log.Fatal(server.ListenAndServe())
}
