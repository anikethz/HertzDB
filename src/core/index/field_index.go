package index

import (
	"fmt"
)

// Array of [startOffset, length]
type FieldIndexMetadata struct {
	Filename  string
	Index     map[rune]IndexMetadata
	Field     string
	IndexType string
}

type FieldIndex struct {
	Field string
	index map[string][][2]int64
}

func NewFieldIndexMetadata(filename string, field string, indexType string) *FieldIndexMetadata {

	newIndex := make(map[rune]IndexMetadata)
	for ch := 'a'; ch <= 'z'; ch++ {
		newIndex[ch] = IndexMetadata{Start: ItoC(0), Length: ItoC(0)}
	}
	return &FieldIndexMetadata{Filename: filename, Index: newIndex, Field: field, IndexType: indexType}
}

func NewFieldIndex(field string) FieldIndex {
	return FieldIndex{Field: field, index: make(map[string][][2]int64)}
}

func (fieldIndex *FieldIndex) UpdateFieldIndex(fieldMetadata *FieldIndexMetadata) error {

	//ToDo: Check FieldIndexMetadata matches FieldIndex
	if fieldMetadata.Field != fieldIndex.Field {
		return fmt.Errorf("failed to index : invalid fieldMetadata: %v for the given fieldIndex: %v", fieldMetadata.Field, fieldIndex.Field)
	}

	//Sort First Character wise
	runeMap := make(map[rune]map[string][][2]int64)
	for k, v := range fieldIndex.index {
		if runeMap[rune(k[0])] == nil {
			runeMap[rune(k[0])] = make(map[string][][2]int64)
		}
		runeMap[rune(k[0])][k] = v
	}

	for k, v := range runeMap {

		x, err := fieldMetadata.updateRuneIndex(k, v)
		if err != nil {
			fmt.Printf("Error: %v :: %v", k, err)
			return err
		}
		fmt.Printf("Updated Rune index : %c with %v\n", k, x)

	}

	return nil

}

// Get Rune Index
func (fieldMetadata FieldIndexMetadata) GetRuneIndex(ch rune) (map[string][][2]int64, error) {

	return DeserializeFromFile[map[string][][2]int64](fieldMetadata.Filename, fieldMetadata.Index[ch].Start.CtoI(), fieldMetadata.Index[ch].Length.CtoI())
}

// Save Rune Index
func (fieldMetadata *FieldIndexMetadata) updateRuneIndex(ch rune, index map[string][][2]int64) (FileIOResult, error) {

	var ioResult FileIOResult
	var err error
	//Check if Not new
	meta := fieldMetadata.Index[ch]
	if meta.Start.CtoI() == 0 && meta.Length.CtoI() == 0 {
		ioResult, err = AppendToFile(fieldMetadata.Filename, index)
		if err != nil {
			return FileIOResult{}, fmt.Errorf("failed to append rune '%c' index for field %v: %v", ch, fieldMetadata.Field, err)
		}
	} else {
		rIndex, err := fieldMetadata.GetRuneIndex(ch)
		if err != nil {
			fmt.Println(meta)
			return FileIOResult{}, fmt.Errorf("failed to retrieve rune '%c' index for field %v: %v", ch, fieldMetadata.Field, err)
		}

		for k := range index {
			rIndex[k] = append(rIndex[k], index[k]...)
		}

		ioResult, err = AppendToFile[map[string][][2]int64](fieldMetadata.Filename, rIndex)

		if err != nil {
			return FileIOResult{}, fmt.Errorf("failed to index %v", err)
		}
		y, err := SetBlankBytes(fieldMetadata.Filename, meta.Start.CtoI(), meta.Length.CtoI())
		if err != nil || y != int(meta.Length.CtoI()) {
			fmt.Printf("Failed to save_++_ rune %c: %v\n", 'D', err)
		}

	}

	fieldMetadata.Index[ch] = IndexMetadata{Start: ItoC(ioResult.offset), Length: ItoC(ioResult.length)}

	return ioResult, err
}
