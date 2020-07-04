package logmine

import (
	"bufio"
	"fmt"
	"github.com/kpfaulkner/gologmine/pkg/logmine/tokenizers"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

type TokenizedLogEntry struct {
	Tokens []tokenizers.DataType
}

// LogMine .. initial implementation.
type LogMine struct {
	tokenizer Tokenizer

	clusterProcessor ClusterProcessor
	//
	tokenizedLogEntries []TokenizedLogEntry
}

func NewLogMine() LogMine {
	lm := LogMine{}
	lm.tokenizer = NewTokenizer()
	lm.clusterProcessor = NewClusterProcessor()

	return lm
}

func (lm *LogMine) ProcessLogsFromReader(reader io.Reader) error {

	// preprocess + datatype identification
	tokenizedLogEntries, err := lm.Preprocess(reader)
	if err != nil {
		return err
	}
	lm.tokenizedLogEntries = tokenizedLogEntries

	// generate clusters.
	err = lm.ClusterGeneration()
	if err != nil {
		return err
	}

	// now process/merge each cluster.
	for _,cluster := range lm.clusterProcessor.clusters {
		tokenizedLogEntry, err := lm.clusterProcessor.ProcessSingleCluster(cluster)
		if err != nil {
			return err
		}

		fmt.Printf("cluster %v\n", tokenizedLogEntry.Tokens)

	}

	// other stuff :)

	return nil
}

// willProcessLine.... eg dont proces if comments etc.
// For now, will just filter out lines where the first
// non white space is a #
func willProcessLine(l string) bool {
	if strings.TrimSpace(l)[0] == '#' {
		return false
	}

	return true
}

// Preprocess will read in ALL log entries from a file(reader)
// and process them. Will return a TokenizedLogEntry for each log line
// read.
func (lm *LogMine) Preprocess(reader io.Reader) ([]TokenizedLogEntry, error) {

	tokenizedLogEntries := []TokenizedLogEntry{}

	// read each log entry and preprocess them.
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := scanner.Text()
		if willProcessLine(l) {
			tokens, err := lm.tokenizer.Tokenize(l)
			if err != nil {
				log.Errorf("Error PreProcess %s\n", err.Error())
				return nil, err
			}
			te := TokenizedLogEntry{Tokens: tokens}
			tokenizedLogEntries = append(tokenizedLogEntries, te)
		}
	}

	return tokenizedLogEntries, nil
}

func (lm *LogMine) ClusterGeneration() error {

	for _, l := range lm.tokenizedLogEntries {
		err := lm.clusterProcessor.AddLogEntry(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (lm *LogMine) GenerateFinalOutput() error {

	return nil
}
