package index

type Document struct {
	Doc    map[string]interface{}
	offset int64
	length int64
}

func NewDocument(_doc map[string]interface{}, offset int64, length int64) *Document {

	doc := Document{Doc: _doc, offset: offset, length: length}
	return &doc
}
