package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {

	listner, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Something went wrong with tcp connection.")
		return
	}
	defer listner.Close()

	for true {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("something wrong with accept!")
			break
		}
		fmt.Println("connection accepted!")
		messageChan := getLinesChannel(conn)
		for msg := range messageChan {
			fmt.Println("msg:", msg)
		}
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	buffer := make([]byte, 8)

	currentLine := ""
	go func() {
		for true {
			size, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					// file done.
					defer fmt.Println("read: end")
					defer close(ch)
					defer f.Close()
					break
				}
				fmt.Println("somthing wrong with reading file bytes.")
				break
			}

			rByte := string(buffer[:size])

			segs := strings.Split(rByte, "\n")
			currentLine += segs[0]
			for _, s := range segs[1:] {
				ch <- currentLine
				currentLine = s
			}

		}
	}()

	return ch
}
