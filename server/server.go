package server

import (
	"fmt"
	"net/http"
	"time"
	"video-streaming/router"
)

func StartServer(port int, timeout time.Duration) error {

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadTimeout:       timeout * time.Second,
		ReadHeaderTimeout: timeout * time.Second,
		WriteTimeout:      timeout * time.Second,
		IdleTimeout:       timeout * time.Second,
		Handler:           router.RouteHandler(),
	}

	fmt.Printf("-------------- Server starting at port %d --------------\n", port)

	return srv.ListenAndServe()
}
