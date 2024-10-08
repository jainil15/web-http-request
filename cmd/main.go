package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"webhttprequest/pkg/config"
	"webhttprequest/pkg/response"
	"webhttprequest/views"

	"github.com/a-h/templ"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

type ForwardRequest struct {
	url     string
	body    map[string]interface{}
	headers map[string]string
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wrappedResponseWriter{w, http.StatusOK}
		h.ServeHTTP(ww, r)
		log.Printf("%v %v %v %v", ww.StatusCode, r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc(
		"GET /health",
		func(w http.ResponseWriter, r *http.Request) {
			response.ResponseHandler(w, &response.Response{Message: "Server Running"})
			return
		},
	)
	// router.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
	// 	im, err := os.ReadFile("web/img/doomguy.ico")
	// 	if err != nil {
	// 		log.Println("Error ->", err)
	// 		response.ErrorHandler(
	// 			w,
	// 			&response.Error{
	// 				StatusCode: http.StatusNotFound,
	// 				Error:      err,
	// 				Message:    "Favicon not found",
	// 			},
	// 		)
	// 		return
	// 	}
	// 	w.Header().Add("Content-type", "image/x-icon")
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write(im)
	// 	return
	// })
	router.Handle("/", templ.Handler(views.Hello("Hellow")))
	port := config.Env.Port
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: LoggingMiddleware(router),
	}
	log.Println("Server starting on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Server Crashed", err)
	}
}
