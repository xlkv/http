package server

import (
	"net"
	"strconv"
	"sync/atomic"

	"http.xlkv.io/internal/request"
	"http.xlkv.io/internal/response"
)

type Server struct {
	Listener net.Listener
	IsClosed atomic.Bool
	Handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (err HandlerError) Error(w response.Writer) (int, error) {
	w.WriteStatusLine(err.StatusCode)
	w.WriteHeaders(w.GetDefaultHeaders(len(err.Message)))
	return w.WriteBody([]byte(err.Message))
}

type Handler func(w *response.Writer, req *request.Request)

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			s.IsClosed.Store(true)
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	writer := response.NewWriter(conn)
	request, err := request.RequestFromReader(conn)
	if err != nil {
		writer.WriteStatusLine(response.BadRequest)
		return
	}
	s.Handler(writer, request)
}

func (s *Server) Close() error {
	s.IsClosed.Store(true)
	return s.Listener.Close()
}

func Serve(port int, handler Handler) (*Server, error) {
	listner, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	server := &Server{
		Listener: listner,
		Handler:  handler,
	}
	go server.listen()
	return server, err
}
