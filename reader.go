package main

import (
	"errors"
	"io"
	"strings"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer file.Close()
		defer close(ch)

		currentLine := ""
		for {
			buf := make([]byte, 8)
			n, err := file.Read(buf)

			if n > 0 {
				str := string(buf[:n])
				parts := strings.Split(str, "\n")
				for i := 0; i < len(parts)-1; i++ {
					ch <- currentLine + parts[i]
					currentLine = ""
				}
				currentLine += parts[len(parts)-1]
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					if currentLine != "" {
						ch <- currentLine
					}
					return
				}
				return
			}
		}
	}()

	return ch
}
