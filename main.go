package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// this program watches the current directory
// if there are any .mangrove files it build them
// into the .school directory and servers them at port 8071

// there should be a subdirectory "theme"

func main() {
	os.Mkdir(".school", os.ModePerm)

	schoolBuilder("/Users/jimhurd/yakdb/gomgen/sample", ".school")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, ".school/"+r.URL.Path[1:])
	})
	_, e := copyFile("./style.css", ".school/style.css")
	if e != nil {
		panic(e)
	}
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
