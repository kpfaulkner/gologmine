package logmine

import (
	"bufio"
	"fmt"
	"github.com/kpfaulkner/gologmine/pkg/logmine/tokenizers"
	log "github.com/sirupsen/logrus"
	"io"
	"sort"
	"strings"
)

type TokenizedLogEntry struct {
	Tokens []tokenizers.DataType

	// number of entries from previous level?
	NumberOfPreviousEntries int
}

func (te TokenizedLogEntry) ToString() string {
	resultStrings := make([]string, len(te.Tokens))
	for i, t := range te.Tokens {
	  resultStrings[i] = string(t)
	}
	return strings.Join(resultStrings," ")
}

// LogMine .. initial implementation.
type LogMine struct {
	tokenizer Tokenizer

	clusterProcessor ClusterProcessor
	//
	tokenizedLogEntries []TokenizedLogEntry

	// distances used for calculations
	distances []float64
}

func NewLogMine(distances []float64) LogMine {
	lm := LogMine{}
	lm.tokenizer = NewTokenizer()
	lm.distances = distances
	lm.clusterProcessor = NewClusterProcessor(lm.distances)


	return lm
}

func (lm *LogMine) ProcessLogsFromReader(reader io.Reader, maxLevel int) error {

	// preprocess + datatype identification
	tokenizedLogEntries, err := lm.Preprocess(reader)
	if err != nil {
		return  err
	}

	// loop through all the levels.
	for level:=0 ; level <= maxLevel; level++ {

		// generate clusters.
		err = lm.ClusterGeneration(tokenizedLogEntries, level)
		if err != nil {
			return  err
		}

		newTokenizedLogEntries := []TokenizedLogEntry{}

		// now process/merge each cluster in the cluster "level"
		for index, cluster := range lm.clusterProcessor.clusters[level] {
			tokenizedLogEntry, err := lm.clusterProcessor.ProcessSingleCluster(cluster)
			if err != nil {
				return  err
			}

			// if first level (ie ALL logs available) then store number of logs.
			if level == 0{
				tokenizedLogEntry.NumberOfPreviousEntries = len(cluster.logsInCluster)
			}

			// record pattern for cluster.
			cluster.PatternForCluster = *tokenizedLogEntry
			lm.clusterProcessor.clusters[level][index] = cluster
			newTokenizedLogEntries = append(newTokenizedLogEntries, *tokenizedLogEntry)
		}
		tokenizedLogEntries = newTokenizedLogEntries
	}

	return  nil
}

// willProcessLine.... eg dont proces if comments etc.
// For now, will just filter out lines where the first
// non white space is a #
func willProcessLine(l string) bool {
	if len(l) > 0 && strings.TrimSpace(l)[0] == '#' {
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

func (lm *LogMine) ClusterGeneration(logs []TokenizedLogEntry, level int) error {

	for _, l := range logs {
		err := lm.clusterProcessor.AddLogEntry(l, level)
		if err != nil {
			return err
		}
	}
	return nil
}

func (lm *LogMine) GenerateUsefulOutput() ([]TokenizedLogEntry, error) {

	lastLevel := len(lm.clusterProcessor.clusters) - 1

	clusters := lm.clusterProcessor.clusters[lastLevel]
	sort.Slice( clusters, func (i int, j int) bool {
		return clusters[i].PatternForCluster.NumberOfPreviousEntries > clusters[j].PatternForCluster.NumberOfPreviousEntries
	})

	results := make([]TokenizedLogEntry, len(clusters))
	for i,c := range lm.clusterProcessor.clusters[lastLevel] {
		results[i] = c.PatternForCluster
	}

	return results, nil
}

// just display to stdout for now.
// order by fewest entries to most.
func (lm *LogMine) DisplayFinalOutput() error {

	lastLevel := len(lm.clusterProcessor.clusters) - 1

	clusters := lm.clusterProcessor.clusters[lastLevel]
	sort.Slice( clusters, func (i int, j int) bool {
		return clusters[i].PatternForCluster.NumberOfPreviousEntries < clusters[j].PatternForCluster.NumberOfPreviousEntries
	})

	for _,c := range lm.clusterProcessor.clusters[lastLevel] {
		fmt.Printf("count %d : pattern %s\n",c.PatternForCluster.NumberOfPreviousEntries, c.PatternForCluster.Tokens)
	}

	return nil
}
