package fileio

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

type FileIOResult struct {
	Offset int64
	Length int64
}

func SerializeWithOffset[T interface{}](name string, item T, offset int64, whence int) (FileIOResult, error) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return FileIOResult{}, err
	}
	defer file.Close()

	// Seek to the given offset
	fileOffset, err := file.Seek(offset, whence)
	if err != nil {
		return FileIOResult{}, fmt.Errorf("failed to seek to offset: %w", err)
	}

	byteCountWriter := ByteCounterWriter{Writer: file}
	encoder := gob.NewEncoder(&byteCountWriter)
	err = encoder.Encode(item)
	if err != nil {
		return FileIOResult{}, err
	}

	return FileIOResult{Offset: fileOffset, Length: byteCountWriter.Count}, nil
}

func SerializeToFile[T interface{}](name string, item T) (FileIOResult, error) {
	return SerializeWithOffset[T](name, item, 0, io.SeekStart)
}

func AppendToFile[T interface{}](name string, item T) (FileIOResult, error) {

	return SerializeWithOffset(name, item, 2, io.SeekEnd)
}

// @deprecated
func setBlankBytes(filename string, offset int64, length int64) (FileIOResult, error) {
	blankBytes := make([]byte, length)
	return SerializeWithOffset(filename, blankBytes, offset, io.SeekStart)
}

func SetBlankBytes(filename string, offset int64, length int64) (int, error) {
	// Open the file for reading and writing
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	blankData := make([]byte, length)

	if _, err := file.Seek(offset, 0); err != nil {
		return 0, fmt.Errorf("failed to seek to offset: %v", err)
	}

	bytes, err := file.Write(blankData)
	if err != nil {
		return 0, fmt.Errorf("failed to write blank bytes: %v", err)
	}

	return bytes, nil
}
