package index

import (
	"fmt"
	"io"

	fileio "github.com/anikethz/HertzDB/src/core/file_io"
)

func SearchTerm(filename string, field string, term string) ([][2]int64, error) {

	indexDocument, _ := DeserializeIndexDocumentMeta(filename)
	fieldIndexMetadata, _ := indexDocument.GetFieldIndexMetadata(field)
	rIndex, err := fieldIndexMetadata.GetRuneIndex(rune(term[0]))
	if err != nil {
		fmt.Println(err)
	}
	return rIndex[term], nil

}

func GetDocument(filename string, locs [][2]int64) ([]string, error) {

	documents := make([]string, 0)

	for _, _loc := range locs {

		doc, err := fileio.DeserializeRawString(filename, _loc[0], _loc[1], io.SeekStart)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		documents = append(documents, doc)
	}
	return documents, nil
}
