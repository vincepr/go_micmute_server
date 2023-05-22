/*		Main entry point.
*		- hard coded Listen-Port 3003 right now, we use a Docker container anyway
*		- We also setup the Routes here (each handler-function handles one Path like "/login")
*/
package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	// Top level context for graceful shuttdown:
	rootCtx := context.Background()
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	initApi(ctx)
	log.Fatal(http.ListenAndServe(":3003", nil))
}

func initApi(ctx context.Context) {
	manager := NewManager(ctx)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/login_receiver", manager.loginReceiverHandler)
	http.HandleFunc("/receiver", manager.serveReceiverHandler)
	http.HandleFunc("/controller", manager.ControllerRequestHandler)

	http.HandleFunc("./debug", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Receivers:",w, len(manager.receivers))
	})
}