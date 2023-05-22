/*		Main entry point.
*		- hard coded Listen-Port 3003 right now, we use a Docker container anyway
*		- We also setup the Routes here (each handler-function handles one Path like "/login")
*/
package main

import (
	"flag"
	"context"
	"log"
	"net/http"
)

// Optional Flags, with default values if not set.
var (
	PORT = flag.String("port", "3003", "the port the server listens to")
	PATH_TO_PUBLICFOLDER = flag.String("files", "/var/www/go_micmute_server", "expected path to the /public folder with index.html")
)

func main() {
	// Top level context for graceful shuttdown:
	rootCtx := context.Background()
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	flag.Parse()
	log.Println("Serving on Port: "+*PORT)
	log.Println("Expect files at in:"+*PATH_TO_PUBLICFOLDER)

	initApi(ctx)
	log.Fatal(http.ListenAndServe(":"+*PORT, nil))
}

func initApi(ctx context.Context) {
	manager := NewManager(ctx)

	http.Handle("/", http.FileServer(http.Dir(*PATH_TO_PUBLICFOLDER)))
	http.HandleFunc("/login_receiver", manager.loginReceiverHandler)
	http.HandleFunc("/receiver", manager.serveReceiverHandler)
	http.HandleFunc("/controller", manager.ControllerRequestHandler)

	http.HandleFunc("./debug", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Receivers:",w, len(manager.receivers))
	})
}