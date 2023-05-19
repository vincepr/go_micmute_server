/*
*
*
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

	testing()

	initApi(ctx)
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func initApi(ctx context.Context) {
	manager := NewManager(ctx)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/login_controller", manager.loginControllerHandler)
	http.HandleFunc("/login_receiver", manager.loginReceiverHandler)
	http.HandleFunc("/receiver", manager.serveReceiverHandler)
	http.HandleFunc("/controller", manager.serveControllerHandler)

	http.HandleFunc("./debug", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Controllers:",w, len(manager.controllers))
		log.Println("Receivers:",w, len(manager.receivers))
	})
}

// TODO: remove this
func testing() {

}
