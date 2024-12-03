package fileio

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

const META_LENGTH int64 = 116
const SEEK_DELI int64 = 5
const META_START_OFFSET int64 = META_LENGTH + SEEK_DELI

func Deserialize[T any](filename string, offset int64, count int64, whence int) (T, error) {
	var result T

	if offset == 0 && count == 0 {
		return result, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return result, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = file.Seek(offset, whence)
	if err != nil {
		return result, fmt.Errorf("failed to seek to offset: %w", err)
	}

	limitedReader := io.LimitedReader{
		R: file,
		N: count,
	}

	decoder := gob.NewDecoder(&limitedReader)
	err = decoder.Decode(&result)
	if err != nil {
		return result, fmt.Errorf("failed to decode data: %w", err)
	}

	return result, nil
}

func DeserializeFromFile[T any](filename string, offset int64, count int64) (T, error) {
	return Deserialize[T](filename, offset, count, io.SeekStart)
}

func DeserializeRawString(filename string, offset int64, count int64, whence int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Seek to the specified offset
	_, err = file.Seek(offset, whence)
	if err != nil {
		return "", fmt.Errorf("failed to seek to offset: %w", err)
	}

	// Read the specified number of bytes
	buffer := make([]byte, count)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read from file: %w", err)
	}

	// Convert the read bytes to a string
	return string(buffer[:n]), nil
}