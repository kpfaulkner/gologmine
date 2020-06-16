package logmine

import "io"

type LogMine struct {

}

func NewLogMine() LogMine {
	lm := LogMine{}

	return lm
}


func (lm *LogMine) Preprocess(reader io.Reader) error {

	// tokenize

	// datatype identification

	return nil
}

func (lm *LogMine) PatternTreeGeneration() error {

	return nil
}

func (lm *LogMine) GenerateFinalOutput() error {

	return nil
}
