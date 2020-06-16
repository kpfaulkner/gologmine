package tokenizers

import (
	"regexp"
)

const (
	// probably need more.. but will do for moment
  TIMEREGEX string = "(00|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9]):(0[0-9]|[0-5][0-9])"
)

const (
	TIME DataType = "TIME"
)

type TimeTokenizer struct {
  timeRE *regexp.Regexp

}

func NewTimeTokenizer() TimeTokenizer{
	tt := TimeTokenizer{}
	tt.timeRE,_ = regexp.Compile(TIMEREGEX)
	return tt
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (tt TimeTokenizer) CheckToken(token string ) bool {
	return tt.timeRE.MatchString(token)
}

func (tt TimeTokenizer) GetDataType() DataType {
	return TIME
}
