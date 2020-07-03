package tokenizers

import (
	"regexp"
)

const (
	//NOTSPACEREGEX string = "[-!$%#^&*()_+|~`{}[]:\";'<>?,./\\W]+"
	NOTSPACEREGEX string = ".*[!@#$%^&*()\"\\[\\]\"{}]+.*"
	//NOTSPACEREGEX string = "(\\w+\\S+|\\S+\\w+|\\w+\\S+\\w+)"
)

const (
	NOTSPACE DataType = "NOTSPACE"
)

type NotSpaceTokenizer struct {
	notSpaceRE *regexp.Regexp
}

func NewNotSpaceTokenizer() NotSpaceTokenizer {
	nt := NotSpaceTokenizer{}
	nt.notSpaceRE, _ = regexp.Compile(NOTSPACEREGEX)
	return nt
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (nt NotSpaceTokenizer) CheckToken(token string) bool {
	return nt.notSpaceRE.MatchString(token)
}

func (nt NotSpaceTokenizer) GetDataType() DataType {
	return NOTSPACE
}
