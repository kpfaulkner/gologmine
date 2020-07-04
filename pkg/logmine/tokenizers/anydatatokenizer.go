package tokenizers

import (
	"regexp"
)

const (
	ANYDATAREGEX string = ".*"
)

const (
	ANYDATA DataType = "*"
)

type AnyDataTokenizer struct {
	anyDataRE *regexp.Regexp
}

func NewAnyDataTokenizer() AnyDataTokenizer {
	at := AnyDataTokenizer{}
	at.anyDataRE, _ = regexp.Compile(ANYDATAREGEX)
	return at
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (at AnyDataTokenizer) CheckToken(token string) bool {
	return at.anyDataRE.MatchString(token)
}

func (at AnyDataTokenizer) GetDataType() DataType {
	return ANYDATA
}
