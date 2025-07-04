package util

import (
	"bufio"
	"os"
)

type Result struct {
	Scanner  *bufio.Scanner
	Err      error
	Filepath string
	Close    func()
}

func CreateFileReaderIterator(path string) (<-chan Result, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	ch := make(chan Result)
	go func() {
		defer close(ch)

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filePath := path + "/" + entry.Name()
			file, err := os.Open(filePath)
			if err != nil {
				ch <- Result{Err: err}
				continue
			}

			scanner := bufio.NewScanner(file)
			ch <- Result{Scanner: scanner, Filepath: filePath, Close: func() { file.Close() }}
		}
	}()
	return ch, nil
}
