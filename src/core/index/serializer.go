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

	return serializeWithOffset(name, item, 0, io.SeekEnd)
}

func setBlankBytes(filename string, offset int64, length int64) (FileIOResult, error) {
	blankBytes := make([]byte, length)
	return serializeWithOffset(filename, blankBytes, offset, io.SeekStart)
}

func updateMetaLength(filename string, offset int64, length int64) (FileIOResult, error) {

	metadata, err := deserializeMetaLength(filename)
	if err != nil {
		return FileIOResult{}, fmt.Errorf("failed to get meta length: %w", err)
	}

	//Update old space with blank bytes
	// blankBytes := make([]byte, metadata.Length.CtoI())
	// serializeWithOffset(filename, blankBytes, metadata.Start.CtoI(), io.SeekStart)
	setBlankBytes(filename, metadata.Start.CtoI(), metadata.Length.CtoI())

	metadata.Start = ItoC(offset)
	metadata.Length = ItoC(length)
	return serializeWithOffset(filename, metadata, 0, io.SeekStart)

}

func UpdateIndexDocumentMetadata(filename string, indexDocument IndexDocument) (FileIOResult, error) {

	ioResult, err := AppendToFile(filename, indexDocument)
	if err != nil {
		return FileIOResult{}, fmt.Errorf("failed to write delimiter: %w", err)
	}
	_, err = updateMetaLength(filename, ioResult.offset, ioResult.length)

	if err != nil {
		return FileIOResult{}, fmt.Errorf("failed to update metadata document: %w", err)
	}

	return ioResult, nil
}
