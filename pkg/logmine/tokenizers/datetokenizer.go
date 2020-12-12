package tokenizers

import (
	"regexp"
)

const (

	// dd/mm/yyyy
	DATE5REGEX string = "(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)"

	// mm/dd/yyyy
	DATE6REGEX string = "(0?[1-9]|1[012])/(0?[1-9]|[12][0-9]|3[01])/((19|20)\\d\\d)"

	// dd-mm-yyyy
	DATE3REGEX string = "(0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])-((19|20)\\d\\d)"

	// mm-dd-yyyy
	DATE4REGEX string = "(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])-((19|20)\\d\\d)"

	// yyyy-mm-dd
	DATE1REGEX string = "((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])"

	// yyyy/mm/dd
	DATE2REGEX string = "((19|20)\\d\\d)/(0?[1-9]|1[012])/(0?[1-9]|[12][0-9]|3[01])"
)

const (
	DATE DataType = "DATE"
)

type RegexValidator struct {
	Re *regexp.Regexp
	Min int
	Max int
}

func NewRegexValidator( s string, min int, max int) RegexValidator {
	rev := RegexValidator{}
	rev.Re,_ = regexp.Compile(s)
	rev.Min = min
	rev.Max = max
	return rev
}

type DateTokenizer struct {
	dateRESlice []RegexValidator
}

func NewDateTokenizer() DateTokenizer {
	dt := DateTokenizer{}
	dt.dateRESlice = []RegexValidator{
		NewRegexValidator(DATE1REGEX,8,10),
		NewRegexValidator(DATE2REGEX,8,10),
		NewRegexValidator(DATE3REGEX,8,10),
		NewRegexValidator(DATE4REGEX,8,10),
		NewRegexValidator(DATE5REGEX,8,10),
		NewRegexValidator(DATE6REGEX,8,10),
	}

	return dt
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (dt DateTokenizer) CheckToken(token string) bool {
	for _, re := range dt.dateRESlice {
		// check length is valid at least.
		if len(token) >= re.Min && len(token) <= re.Max {
			if re.Re.MatchString(token) {
				return true
			}
		}
	}
	return false
}

func (dt DateTokenizer) GetDataType() DataType {
	return DATE
}
