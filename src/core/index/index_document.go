package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/anikethz/HertzDB/src/core/utils"
)

const BUFFER_SIZE = 1000

type IndexMetadata struct {
	Start  ConstantInteger
	Length ConstantInteger
}

// Metadata starts at SeekStart + 9
// Map of field name
type IndexDocument struct {
	Name          string
	Json_Filename string
	Metadata      map[string]IndexMetadata
}

func purgeAndCreateNewFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return nil
}

func NewIndexDocument(name string, json_fileName string) (*IndexDocument, error) {

	meta := make(map[string]IndexMetadata)
	indexDocument := IndexDocument{Name: name, Json_Filename: json_fileName, Metadata: meta}
	purgeAndCreateNewFile(name)
	ioResult, err := serializeWithOffset(name, indexDocument, META_START_OFFSET, io.SeekStart)
	if err != nil {
		return nil, errors.New("failed creating index")
	}
	indexMetadata := IndexMetadata{Start: ItoC(ioResult.offset), Length: ItoC(ioResult.length)}

	ioResult, err = serializeToFile(name, indexMetadata)
	if err != nil {
		return nil, errors.New("failed creating index")
	}

	return &indexDocument, nil
}

// Get Field Index
func (document IndexDocument) GetFieldIndexMetadata(field string) (FieldIndexMetadata, error) {
	if metadata, ok := document.Metadata[field]; ok {
		return DeserializeFromFile[FieldIndexMetadata](document.Name, metadata.Start.CtoI(), metadata.Length.CtoI())
	}
	return FieldIndexMetadata{}, fmt.Errorf("field %s not found in document", field)
}

func (document *IndexDocument) updateFieldMetadata(fieldMetadata FieldIndexMetadata) (FileIOResult, error) {
	var ioResult FileIOResult
	var err error
	//Check if Not new
	meta := document.Metadata[fieldMetadata.Field]
	if meta.Start.CtoI() == 0 && meta.Length.CtoI() == 0 {
		ioResult, err = AppendToFile(fieldMetadata.Filename, fieldMetadata)
	} else {
		ioResult, err = AppendToFile[FieldIndexMetadata](fieldMetadata.Filename, fieldMetadata)

		if err != nil {
			return FileIOResult{}, fmt.Errorf("failed to index %v", err)
		} else {
			SetBlankBytes(fieldMetadata.Filename, meta.Start.CtoI(), meta.Length.CtoI())
		}

	}

	document.Metadata[fieldMetadata.Field] = IndexMetadata{Start: ItoC(ioResult.offset), Length: ItoC(ioResult.length)}
	return ioResult, err

}

// Convert to JSON util
// New Index
func (fieldIndex *FieldIndex) ingestDocument(document Document) {
	if _, ok := document.Doc[fieldIndex.Field]; !ok || document.Doc[fieldIndex.Field] == nil {
		// fmt.Printf("key not found: %v", document.Doc)
		return
	}
	tokens := utils.LowCaseTokenizer(document.Doc[fieldIndex.Field].(string))

	for i := range tokens {
		if fieldIndex.index[tokens[i]] == nil {
			fieldIndex.index[tokens[i]] = make([][2]int64, 0)
		}
		loc := [2]int64{document.offset, document.length}
		(fieldIndex.index)[tokens[i]] = append((fieldIndex.index)[tokens[i]], loc)
	}
}

func (indexDocument *IndexDocument) IndexTextFields(field string, documents *Documents) FieldIndexMetadata {

	//Get fieldMetadata
	fieldMetadata, err := indexDocument.GetFieldIndexMetadata(field)

	if err != nil {
		fieldMetadata = *NewFieldIndexMetadata(indexDocument.Name, field, "text")
	}

	fieldIndex := NewFieldIndex(field)

	_docs := *documents.docs

	for i := range _docs {
		fieldIndex.ingestDocument(_docs[i])
	}
	err = fieldIndex.UpdateFieldIndex(&fieldMetadata)

	if err != nil {
		fmt.Println(err)
	}

	_, err = indexDocument.updateFieldMetadata(fieldMetadata)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(indexDocument)
	UpdateIndexDocumentMetadata(fieldMetadata.Filename, indexDocument)

	return fieldMetadata
}

func (indexDocument *IndexDocument) ParseEntireFile(fields []string) {
	file, _ := os.Open(indexDocument.Json_Filename)
	//Read Document
	stringSlice := make([]string, 0, 1000)
	buffer := make([]byte, BUFFER_SIZE)

	counter := 0
	documentByteSize := BUFFER_SIZE * 10
	documentBytes := make([]byte, 0, documentByteSize)
	startOffset, endOffset := int64(0), int64(0)
	documents := NewDocuments()

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Something failed while reading the file: %v\n", err)
			}
			break
		}
		for i := 0; i < bytesRead; i++ {
			endOffset++
			documentBytes = append(documentBytes, buffer[i])

			// Check for the end of a JSON document (assuming newline as delimiter)
			if buffer[i] == '\n' {
				// Deserialize the JSON document
				var jsonDoc map[string]interface{}
				err := json.Unmarshal(documentBytes, &jsonDoc)
				if err != nil {
					fmt.Printf("Error unmarshaling JSON: %v\nInput JSON: %s\n", err, string(documentBytes))
				} else {
					_ = append(stringSlice, string(documentBytes))
					documents.ProcessNewDocumentAndIndex(string(documentBytes), startOffset, (endOffset - startOffset))
				}
				// Reset documentBytes for the next document
				documentBytes = make([]byte, 0)
				counter++
				startOffset = endOffset // Update startOffset to the next byte after the newline
				if counter%1000 == 0 {
					for _, field := range fields {
						indexDocument.IndexTextFields(field, documents)
					}
					documents = NewDocuments()
				}
			}
		}
	}
	for _, field := range fields {
		indexDocument.IndexTextFields(field, documents)
	}
	fmt.Printf("Number of documents: %v\n", counter)
}
