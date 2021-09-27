package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// custom FS handler to restrict dir access
type restrictedFS struct {
	fs http.FileSystem
}

func (rfs restrictedFS) Open(path string) (http.File, error) {
	// filesystem handle to the 'path' (relative to 'assetsDir')
	f, err := rfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	// get file system info for path
	s, err := f.Stat()
	if s.IsDir() {
		// if path is a directory
		// check if an 'index (default)' page exists
		indxFile := filepath.Join(path, config().DefaultPage)
		// try open file
		if _, err := rfs.fs.Open(indxFile); err != nil {
			// file exist close the file handle
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			// else return resource not found error
			return nil, err

		}
	}
	// if reached here pass on th file handle
	return f, nil
}

func main() {
	// multiplexer to route req to hanlder functions
	muxHandler := http.NewServeMux()

	// handler to serve pages / static content
	// use 'assets_dir' as the root of FS object for handler
	// decorate default FS object to restric dir access
	fs := http.FileServer(restrictedFS{http.Dir(config().AssetsDir)})
	// struct to encapsulate handler methods
	rh := reqHandler{filehandler: fs}

	// handle base path
	muxHandler.HandleFunc("/", rh.checkAndServeFile)
	// handle /login POST
	muxHandler.HandleFunc("/submit", rh.submit)

	// dummy api for testing html response
	muxHandler.HandleFunc("/api/hello", rh.hello)
	// dummy api for testing json reponse
	muxHandler.HandleFunc("/api/datetime", rh.datetime)

	// http server instance
	server := http.Server{
		Addr: fmt.Sprintf(":%s", config().Port),
		// our mux as handler for http requests
		Handler:      muxHandler,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 10,
	}
	log.Printf("Starting http server on PORT: %s\n", config().Port)
	// start server and listen for req to serve
	// this is a 'blocking' call
	err := server.ListenAndServe()
	// check type of termination error
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}

}
