package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"http.xlkv.io/internal/headers"
	"http.xlkv.io/internal/request"
	"http.xlkv.io/internal/response"
	"http.xlkv.io/internal/server"
)

const port = 42069

func main() {

	var handler = func(w *response.Writer, req *request.Request) {
		switch req.RequestLine.RequestTarget {
		case "/video":
			vidByte, err := os.ReadFile("/Users/azizbek/gopro/http/assets/vim.mp4")
			if err != nil {
				fmt.Println("error:", err)
				w.WriteStatusLine(response.ServerError)
				return
			}
			fmt.Println(len(vidByte))
			w.WriteStatusLine(response.OK)
			headers := w.GetDefaultHeaders(len(vidByte))
			headers.Override("Content-Type", "video/mp4")
			w.WriteHeaders(headers)
			w.WriteBody(vidByte)
		case "/yourproblem":
			body := []byte(`<html>
							  <head>
							    <title>400 Bad Request</title>
							  </head>
							  <body>
							    <h1>Bad Request</h1>
							    <p>Your request honestly kinda sucked.</p>
							  </body>
							</html>`)
			w.WriteStatusLine(response.BadRequest)
			headers := w.GetDefaultHeaders(len(body))
			headers.Override("Content-Type", "text/html")
			w.WriteHeaders(headers)
			w.WriteBody(body)
		case "/myproblem":
			body := []byte(`<html>
								  <head>
								    <title>500 Internal Server Error</title>
								  </head>
								  <body>
								    <h1>Internal Server Error</h1>
								    <p>Okay, you know what? This one is on me.</p>
								  </body>
								</html>`)
			w.WriteStatusLine(response.ServerError)
			headers := w.GetDefaultHeaders(len(body))
			headers.Override("Content-Type", "text/html")
			w.WriteHeaders(headers)
			w.WriteBody(body)
		default:
			if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {

				baseURL := "https://httpbin.org"
				reqUrl := fmt.Sprintf("%v%v", baseURL, strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin"))
				header := w.GetDefaultHeaders(0)
				var full []byte
				buf := make([]byte, 1024)
				trailers := headers.NewHeaders()

				resp, err := http.Get(reqUrl)
				if err != nil {
					return
				}

				header.Delete("content-length")
				header.Add("transfer-encoding", "chunked")
				header.Add("trailer", "X-Content-SHA256, X-Content-Length")

				w.WriteStatusLine(response.OK)
				w.WriteHeaders(header)

				headerTrailer := resp.Header.Get("Trailer")
				if headerTrailer != "" {
					header.Add("Trailer", headerTrailer)
				}

				for {
					n, err := resp.Body.Read(buf)
					if err != nil {
						if errors.Is(err, io.EOF) {
							w.WriteChunkedBodyDone()
							hash := sha256.Sum256(full)
							len := len(full)
							trailers.Add("X-Content-SHA256", fmt.Sprintf("%x", hash))
							trailers.Add("X-Content-Length", fmt.Sprintf("%v", len))
							w.WriteTrailers(trailers)
							return
						}
					}
					full = append(full, buf[:n]...)
					w.WriteChunkedBody(buf[:n])
				}
			}
			body := []byte(`<html>
							  <head>
							    <title>200 OK</title>
							  </head>
							  <body>
							    <h1>Success!</h1>
							    <p>Your request was an absolute banger.</p>
							  </body>
							</html>`)
			w.WriteStatusLine(response.OK)
			headers := w.GetDefaultHeaders(len(body))
			headers.Override("Content-Type", "text/html")
			w.WriteHeaders(headers)
			w.WriteBody(body)
		}
	}

	server, err := server.Serve(port, handler)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
