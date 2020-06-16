package tokenizers

import (
	"regexp"
)

const (
  //WORDREGEX string = "[a-zA-Z_0-9]+"

  // I'm starting to think that WORD means a combination of letters and at least 1 number.
  // Going by videos, ACM paper and PPT.  But still is never specific said anywhere.
  WORDREGEX string = "(?:[0-9]+[a-z]|[a-z]+[0-9])[a-z0-9]*"
)

const (
	WORD DataType = "WORD"
)

type WordTokenizer struct {
  wordRE *regexp.Regexp

}

func NewWordTokenizer() WordTokenizer{
	wt := WordTokenizer{}
	wt.wordRE,_ = regexp.Compile(WORDREGEX)
	return wt
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (wt WordTokenizer) CheckToken(token string ) bool {
	return wt.wordRE.MatchString(token)
}

func (wt WordTokenizer) GetDataType() DataType {
	return WORD
}
