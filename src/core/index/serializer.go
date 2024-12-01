package index

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

type FileIOResult struct {
	offset int64
	length int64
}

func serializeWithOffset[T interface{}](name string, item T, offset int64, whence int) (FileIOResult, error) {
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

	return FileIOResult{offset: fileOffset, length: byteCountWriter.Count}, nil
}

func serializeToFile[T interface{}](name string, item T) (FileIOResult, error) {
	return serializeWithOffset[T](name, item, 0, io.SeekStart)
}

func AppendToFile[T interface{}](name string, item T) (FileIOResult, error) {

	return serializeWithOffset(name, item, 2, io.SeekEnd)
}

// @deprecated
func setBlankBytes(filename string, offset int64, length int64) (FileIOResult, error) {
	blankBytes := make([]byte, length)
	return serializeWithOffset(filename, blankBytes, offset, io.SeekStart)
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

func updateMetaLength(filename string, offset int64, length int64) (FileIOResult, error) {

	metadata, err := deserializeMetaLength(filename)
	if err != nil {
		return FileIOResult{}, fmt.Errorf("failed to get meta length: %w", err)
	}

	//Update old space with blank bytes
	// blankBytes := make([]byte, metadata.Length.CtoI())
	// serializeWithOffset(filename, blankBytes, metadata.Start.CtoI(), io.SeekStart)
	SetBlankBytes(filename, metadata.Start.CtoI(), metadata.Length.CtoI())

	metadata.Start = ItoC(offset)
	metadata.Length = ItoC(length)
	return serializeWithOffset(filename, metadata, 0, io.SeekStart)

}

func UpdateIndexDocumentMetadata(filename string, indexDocument *IndexDocument) (FileIOResult, error) {

	ioResult, err := AppendToFile(filename, indexDocument)
	if err != nil {
		return FileIOResult{}, fmt.Errorf("failed to write delimiter: %w", err)
	}
	x, err := updateMetaLength(filename, ioResult.offset, ioResult.length)

	fmt.Println("{", x.offset, ",", x.length, "}")

	if err != nil {
		return FileIOResult{}, fmt.Errorf("failed to update metadata document: %w", err)
	}

	return ioResult, nil
}
