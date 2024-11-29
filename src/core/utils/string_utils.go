package utils

import (
	"github.com/bzick/tokenizer"
)

func LowCaseTokenizer(input string) []string {
	tokens := make([]string, 0)
	parser := tokenizer.New()
	parser.AllowKeywordUnderscore()

	stream := parser.ParseString(input)
	defer stream.Close()

	for stream.IsValid() {
		if stream.CurrentToken().Is(tokenizer.TokenKeyword) {
			tokens = append(tokens, stream.CurrentToken().ValueString())
		}
		stream.GoNext()
	}

	return tokens
}
