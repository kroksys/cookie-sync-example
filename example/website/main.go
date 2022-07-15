package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// This file is a web-server that serves content of ./public
// folder on the web.

//go:embed public
var publicFS embed.FS

// Port used to host the website. By default set to port 3000.
var port = flag.String("port", "3000", "port to serve on")
var debug = flag.Bool("d", false, "show debug info")

func main() {
	// Get ./public as root ./
	contentFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if *debug {
		r.Use(gin.LoggerWithWriter(gin.DefaultWriter))
	}
	r.Use(
		gin.Recovery(),
		static.Serve("/", embedFileSystem{http.FS(contentFS)}),
	)

	fmt.Println("Started hosting for https://website.com localy on :" + *port)
	r.Run(fmt.Sprintf(":%s", *port))
}

// interface implementation for gin to handle embed file system
type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}
