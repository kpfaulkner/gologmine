package logmine

import (
	"github.com/kpfaulkner/gologmine/pkg/logmine/tokenizers"
	"strings"
)

type Tokenizer struct {

	// Special Key phrases that we want to not get rid of.
	specialKeyPhrases []string

	// special symbols that we want to surround with spaces to it becomes a token itself (eg, '=')
	tokenDelimiters []string

	// tokenizers that will process each token.
	tokenizerCheckers []tokenizers.TokenizerChecker
}

func NewTokenizer() Tokenizer {
	t := Tokenizer{}
	t.specialKeyPhrases = []string{"user"}
	t.tokenDelimiters = []string{"="}

	// list all the specifictokenizers that we're going to use.
	t.tokenizerCheckers = []tokenizers.TokenizerChecker{tokenizers.NewDateTokenizer(), tokenizers.NewTimeTokenizer(),
		tokenizers.NewIPV4Tokenizer(), tokenizers.NewNumberTokenizer(),
		tokenizers.NewWordTokenizer(), tokenizers.NewNotSpaceTokenizer()}

	//																					 tokenizers.NewAnyDataTokenizer()}
	return t
}

// addSpacesToLog just adds spaces around specific characters
// to help with tokenization/datatype identification.
// Just equals (=) for now...but will expand this to cover many chars
func (t Tokenizer) addSpacesToLog(log string) string {
	l := log
	for _, ss := range t.tokenDelimiters {
		l = strings.ReplaceAll(l, ss, " "+ss+" ")
	}
	return l
}

// processToken does a number of things.
// This runs the token string through all specific tokenizers
// and results in the associated DataType
func (t Tokenizer) processToken(token string) (tokenizers.DataType, error) {

	for _, tokenizer := range t.tokenizerCheckers {
		if tokenizer.CheckToken(token) {
			return tokenizer.GetDataType(), nil
		}
	}

	// no specific datatype, so return token string as DataType
	return tokenizers.DataType(token), nil
}

func (t Tokenizer) Tokenize(log string) ([]tokenizers.DataType, error) {

	// add spaces around equals signs... probably more.
	modifiedLog := t.addSpacesToLog(log)
	tokens := strings.Split(modifiedLog, " ")

	dataTypeArray := []tokenizers.DataType{}
	for _, token := range tokens {
		//fmt.Printf("token: %s\n", token)
		trimmedToken := strings.TrimSpace(token)
		dt, err := t.processToken(trimmedToken)
		if err != nil {
			return nil, err
		}

		dataTypeArray = append(dataTypeArray, dt)
	}

	return dataTypeArray, nil
}
