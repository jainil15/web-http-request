package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"webhttprequest/pkg/config"
	parseJSON "webhttprequest/pkg/parse"
	"webhttprequest/pkg/response"
	"webhttprequest/views/components"
	forms "webhttprequest/views/forms"

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
	method  string
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

func SendRequest(req ForwardRequest) (*http.Response, error) {
	marshalled, err := json.Marshal(req.body)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(req.method, req.url, bytes.NewBuffer(marshalled))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	return response, nil
}

//go:embed static
var Static embed.FS

func main() {
	router := http.NewServeMux()
	router.HandleFunc(
		"GET /health",
		func(w http.ResponseWriter, r *http.Request) {
			response.ResponseHandler(w, &response.Response{Message: "Server Running"})
			return
		},
	)
	router.HandleFunc(
		"POST /request",
		func(w http.ResponseWriter, r *http.Request) {
			url := r.FormValue("url")
			keys := r.Form["key"]
			values := r.Form["value"]
			method := r.FormValue("method")
			requestBody := make(map[string]interface{})
			for i, key := range keys {
				requestBody[key] = values[i]
			}
			req := ForwardRequest{
				url:  url,
				body: requestBody,
				headers: map[string]string{
					"Content-type": "application/json",
				},
				method: method,
			}
			log.Printf("Url: %v\n", url)
			log.Printf("%v: %v\n", keys, values)
			res, err := SendRequest(req)
			if err != nil {
				response.ErrorHandler(w, &response.Error{
					StatusCode: 500,
					Error:      err,
					Message:    fmt.Sprintf("%v\n", err),
				})
				return
			}
			resBody := make(map[string]interface{})
			parseJSON.Decode2(&resBody, res.Body)
			marshalled, err := json.MarshalIndent(resBody, "", "  ")
			if err != nil {
				response.ErrorHandler(w, &response.Error{
					StatusCode: 500,
					Error:      err,
					Message:    fmt.Sprintf("%v\n", err),
				})
				return
			}
			log.Printf("Reponse: \n%v\n", string(marshalled))
			w.Header().Add("Content-type", "text/html")
			w.WriteHeader(200)
			components.Response(marshalled).Render(context.Background(), w)
			return
		},
	)

	// Views
	homePage := forms.RequestForm()
	router.Handle("/", templ.Handler(homePage))
	router.Handle("/static/", http.FileServer(http.FS(Static)))
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
