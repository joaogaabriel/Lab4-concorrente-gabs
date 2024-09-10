package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string, resultChan chan<- struct {
	path string
	sum  int
	err  error
}) {
	data, err := readFile(filePath)
	if err != nil {
		resultChan <- struct {
			path string
			sum  int
			err  error
		}{path: filePath, sum: 0, err: err}
		return
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	resultChan <- struct {
		path string
		sum  int
		err  error
	}{path: filePath, sum: _sum, err: nil}
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var totalSum int64
	sums := make(map[int][]string)

	resultChan := make(chan struct {
		path string
		sum  int
		err  error
	})

	for _, path := range os.Args[1:] {
		go sum(path, resultChan)
	}

	for i := 0; i < (len(os.Args) - 1); i++ {
		result := <-resultChan
		if result.err != nil {
			continue
		}

		totalSum += int64(result.sum)
		sums[result.sum] = append(sums[result.sum], result.path)
	}
	fmt.Println(totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}
