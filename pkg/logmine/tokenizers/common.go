package tokenizers

type DataType string

type TokenizerChecker interface {
	CheckToken(token string) bool
	GetDataType() DataType
}
