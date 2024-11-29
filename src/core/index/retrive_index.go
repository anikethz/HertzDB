package index

import (
	"fmt"
	"io"
)

func SearchTerm(filename string, field string, term string) ([][2]int64, error) {

	indexDocument, _ := DeserializeIndexDocumentMeta(filename)
	fieldIndexMetadata, _ := indexDocument.GetFieldIndexMetadata(field)

	rIndex, _ := fieldIndexMetadata.GetRuneIndex(rune(term[0]))
	return rIndex[term], nil

}

func GetDocument(filename string, locs [][2]int64) ([]string, error) {

	documents := make([]string, 0)

	for _, _loc := range locs {
		// doc, err := DeserializeFromFile[map[string]interface{}](filename, _loc[0], _loc[1])
		doc, err := DeserializeRawString(filename, _loc[0], _loc[1], io.SeekStart)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println(doc)
		documents = append(documents, doc)
	}
	return documents, nil
}
