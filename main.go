package main

import (
	"log"
	"net/http"
	"text/template"
)

const port = "8080"
const html = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>File upload</title>
</head>
<body>
    <h1>File Upload</h1>
    <div>
        <form action="/upload" method="post">
            <input type="file" name="file" multiple id="">
            <button type="submit">Upload</button>
        </form>
    </div>
</body>
</html>`

func main() {
	// fmt.Println("Hello, World!")

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
	})

	log.Fatal(server.ListenAndServe())
}
