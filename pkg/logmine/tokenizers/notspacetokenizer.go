package tokenizers

import (
	"regexp"
)

const (
  NOTPSACERREGEX string = "^[\t\n\x0B\f\r]+"
)

const (
	NOTSPACE DataType = "NOTSPACE"
)

type NotSpaceTokenizer struct {
  notSpaceRE *regexp.Regexp

}

func NewNotSpaceTokenizer() NotSpaceTokenizer{
	nt := NotSpaceTokenizer{}
	nt.notSpaceRE,_ = regexp.Compile(NOTPSACERREGEX)
	return nt
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (nt NotSpaceTokenizer) CheckToken(token string ) bool {
	return nt.notSpaceRE.MatchString(token)
}

func (nt NotSpaceTokenizer) GetDataType() DataType {
	return NOTSPACE
}
