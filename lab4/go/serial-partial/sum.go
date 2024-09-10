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
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	return _sum, nil
}

func chunkSum(chunk []byte) int {
	sum := 0
	for _, b := range chunk {
		sum += int(b)
	}
	return sum
}

func similarity(file1Chunks, file2Chunks []int) float64 {
	counter := 0
	file2Copy := make([]int, len(file2Chunks))
	copy(file2Copy, file2Chunks)

	for _, sum := range file1Chunks {
		for i, val := range file2Copy {
			if sum == val {
				counter++
				file2Copy = append(file2Copy[:i], file2Copy[i+1:]...)
				break
			}
		}
	}
	return float64(counter) / float64(len(file1Chunks))
}

// main handles the command line arguments and file comparisons
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	fileFingerprints := make(map[string][]int)

	for _, path := range os.Args[1:] {
		data, err := readFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			continue
		}
		fingerprints := []int{}
		for _, b := range data {
			fingerprints = append(fingerprints, chunkSum([]byte{b}))
		}

		fileFingerprints[path] = fingerprints
	}

	for i := 0; i < len(os.Args)-1; i++ {
		for j := i + 1; j < len(os.Args)-1; j++ {
			file1 := os.Args[i+1]
			file2 := os.Args[j+1]
			fingerprint1 := fileFingerprints[file1]
			fingerprint2 := fileFingerprints[file2]
			similarityScore := similarity(fingerprint1, fingerprint2)
			fmt.Printf("Similarity between %s and %s: %.2f%%\n", file1, file2, similarityScore*100)
		}
	}
}
