package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

const port = "8080"
const html = `
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
        <form action="/upload" method="post" enctype="multipart/form-data">
            <input type="file" name="file" multiple id="files" >
            <button type="submit">Upload</button>
        </form>
    </div>
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
		// send html file
		t, err := template.New("index").Parse(html)
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

			file, err := os.Create(part.FileName())
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
		fmt.Fprintf(w, "file upload completed!")
	})

	log.Fatal(server.ListenAndServe())
}
