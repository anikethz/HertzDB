package index

import "github.com/anikethz/HertzDB/src/core/utils"

type Documents struct {
	docs *[]Document
}

func NewDocuments() *Documents {
	documentArray := make([]Document, 0, 10000)
	return &Documents{docs: &documentArray}
}

func (documents *Documents) ProcessNewDocumentAndIndex(serializedJson string, offset int64, length int64) {
	document := NewDocument(utils.ConvertStoMap(serializedJson), offset, length)
	// Initialize docs if nil
	if documents.docs == nil {
		docs := make([]Document, 0)
		documents.docs = &docs
	}
	*documents.docs = append(*documents.docs, *document)
}
