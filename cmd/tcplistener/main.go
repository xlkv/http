package main

import (
	"fmt"
	"net"

	"http.xlkv.io/internal/request"
)

func main() {
	listner, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Something went wrong with tcp connection:", err)
		return
	}
	defer listner.Close()

	for true {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("something wrong with accept!")
			break
		}
		req, err := request.RequestFromReader(conn)
		if err != nil {
			conn.Write([]byte(err.Error()))
			return
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			str := fmt.Sprintf("- %v: %v", key, value)
			fmt.Println(str)
		}
	}
}

// func readParams(f io.ReadCloser) <-chan string {

// 	ch := make(chan string)

// 	go func() {
// 		defer close(ch)
// 		defer f.Close()
// 		buffer := make([]byte, 8)

// 		var data []byte

// 		for true {
// 			size, err := f.Read(buffer)

// 			if err != nil {
// 				fmt.Println("something went wrong with reading data from conn.")
// 				return
// 			}

// 			data = append(data, buffer[:size]...)

// 			headersBytes, bodyBytes, ok := bytes.Cut(data, []byte("\r\n\r\n"))

// 			if !ok {
// 				continue
// 			}

// 			contentLength := 0

// 			headers := strings.Split(string(headersBytes), "\r\n")

// 			for _, header := range headers {
// 				parts := strings.Split(header, ": ")
// 				if len(parts) != 2 {
// 					continue
// 				}
// 				if parts[0] == "Content-Length" {
// 					contentLength, err = strconv.Atoi(parts[1])
// 				}
// 			}
// 			if len(bodyBytes) < contentLength {
// 				continue
// 			}

// 			for _, header := range headers {
// 				ch <- header
// 			}

// 			ch <- string(bodyBytes[:contentLength])
// 			return
// 		}

// 	}()

// 	return ch
// }

// func bodyReadChan(f io.ReadCloser) (string, <-chan string) {
// 	buffer := make([]byte, 4096)
// 	size, err := f.Read(buffer)
// 	if err != nil {
// 		if errors.Is(err, io.EOF) {

// 		}
// 		return "", nil
// 	}
// 	body := strings.SplitN(string(buffer[:size]), "\n", 2)[1]

// 	return body, getLinesChannel(f)
// }

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	ch := make(chan string)

// 	go func() {
// 		defer close(ch)
// 		defer f.Close()

// 		buffer := make([]byte, 8)

// 		var data []byte

// 		headers := make(map[string]string)

// 		currentLine := ""

// 		for true {
// 			size, err := f.Read(buffer)
// 			if err != nil {
// 				fmt.Println("somthing wrong with reading file bytes.")
// 				break
// 			}

// 			data = append(data, buffer[:size]...)

// 			idx := bytes.Index(data, []byte("\r\n\r\n"))

// 			if idx == -1 {
// 				continue
// 			}

// 			headerBytes := data[:idx]
// 			bodyBytes := data[idx+4:]

// 			// headers := strings.Split(string(headerBytes), "\r\n")

// 			rByte := string(buffer[:size])

// 			segs := strings.Split(rByte, "\r\n")
// 			currentLine += segs[0]
// 			for _, s := range segs[1:] {
// 				ch <- currentLine
// 				parts := strings.Split(currentLine, ": ")
// 				if len(parts) == 2 {
// 					headers[parts[0]] = parts[1]
// 				}
// 				currentLine = s
// 				// if currentLine == "" {
// 				// 	fmt.Println("true")
// 				// 	body := make([]byte, length)
// 				// 	size, err := io.ReadFull(f, body)
// 				// 	if err != nil {
// 				// 		return
// 				// 	}
// 				// ch <- string(body[:size])
// 				// }
// 			}

// 			lengthStr, ok := headers["Content-Length"]
// 			if !ok {
// 				return
// 			}
// 			length, err := strconv.Atoi(lengthStr)
// 			if err != nil {
// 				return
// 			}

// 			bodyBuf := make([]byte, length)

// 			size, err = f.Read(bodyBuf)

// 			if err != nil {
// 				return
// 			}

// 			parts := bytes.Split(bodyBuf, []byte("\r\n\r\n"))

// 			if len(parts) != 2 {
// 				return
// 			}

// 			ch <- string(parts[1])
// 		}
// 	}()

// 	return ch
// }
