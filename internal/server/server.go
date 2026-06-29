package server

import (
	"bytes"
	"io"
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

func (err HandlerError) Error(w io.Writer) (int, error) {
	response.WriteStatusLine(w, err.StatusCode)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(err.Message)))
	return w.Write([]byte(err.Message))
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

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
	buf := bytes.NewBuffer([]byte{})
	request, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.BadRequest)
		return
	}
	handlerErr := s.Handler(buf, request)
	if handlerErr != nil {
		handlerErr.Error(conn)
		return
	}
	b := buf.Bytes()
	headers := response.GetDefaultHeaders(len(b))
	response.WriteStatusLine(conn, response.OK)
	response.WriteHeaders(conn, headers)
	conn.Write(b)
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
