package tokenizers

import (
	"regexp"
)

const (
  NUMBERREGEX string = "^(-[0-9]*[.])?([0-9]+)|(0x[0-9A-F]+)$"
)

const (
	NUMBER DataType = "NUMBER"
)

type NumberTokenizer struct {
  numberRE *regexp.Regexp

}

func NewNumberTokenizer() NumberTokenizer{
	nt := NumberTokenizer{}
	nt.numberRE,_ = regexp.Compile(NUMBERREGEX)
	return nt
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (nt NumberTokenizer) CheckToken(token string ) bool {
	return nt.numberRE.MatchString(token)
}

func (nt NumberTokenizer) GetDataType() DataType {
	return NUMBER
}
