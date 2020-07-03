package tokenizers

import (
	"regexp"
)

const (

	// not using one from ppt slide deck.. had issues in go and cant be bothered debugging it (yet) :)
	// stole from SO
	IPV4REGEX string = "(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])"
)

const (
	IPV4 DataType = "IPV4"
)

type IPV4Tokenizer struct {
	ipv4RE *regexp.Regexp
}

func NewIPV4Tokenizer() IPV4Tokenizer {
	it := IPV4Tokenizer{}
	it.ipv4RE, _ = regexp.Compile(IPV4REGEX)
	return it
}

// CheckDate checks a number of different date formats and indicates if a match is found.
func (it IPV4Tokenizer) CheckToken(token string) bool {
	return it.ipv4RE.MatchString(token)
}

func (it IPV4Tokenizer) GetDataType() DataType {
	return IPV4
}
