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

type DateTokenizer struct {
	date1RE *regexp.Regexp
	date2RE *regexp.Regexp
	date3RE *regexp.Regexp
	date4RE *regexp.Regexp
	date5RE *regexp.Regexp
	date6RE *regexp.Regexp

	dateRESlice []*regexp.Regexp
}

func NewDateTokenizer() DateTokenizer {
	dt := DateTokenizer{}
	dt.date1RE, _ = regexp.Compile(DATE1REGEX)
	dt.date2RE, _ = regexp.Compile(DATE2REGEX)
	dt.date3RE, _ = regexp.Compile(DATE3REGEX)
	dt.date4RE, _ = regexp.Compile(DATE4REGEX)
	dt.date5RE, _ = regexp.Compile(DATE5REGEX)
	dt.date6RE, _ = regexp.Compile(DATE6REGEX)

	dt.dateRESlice = []*regexp.Regexp{dt.date1RE, dt.date2RE, dt.date3RE, dt.date4RE, dt.date5RE, dt.date6RE}

	return dt
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (dt DateTokenizer) CheckToken(token string) bool {
	for _, re := range dt.dateRESlice {
		if re.MatchString(token) {
			return true
		}
	}
	return false
}

func (dt DateTokenizer) GetDataType() DataType {
	return DATE
}
