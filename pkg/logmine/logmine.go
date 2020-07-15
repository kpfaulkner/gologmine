package logmine

import (
	"bufio"
	"fmt"
	"github.com/kpfaulkner/gologmine/pkg/logmine/tokenizers"
	log "github.com/sirupsen/logrus"
	"io"
	"sort"
	"strings"
	"sync"
)

type ClusterMessage struct {
	cluster Cluster
	originalIndex int
}

type TokenizedLogEntryMessage struct {
  tokenizedLogEntry TokenizedLogEntry
  originalClusterIndex int
}

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
	return strings.Join(resultStrings, " ")
}

// Simplify the entries. This is to reduce the visual clutter (and would make any further analysis probably invalid/impossible
// Basically remove continguous NOSPACE, WORD and * entries into a single one.
// So we dont get entries such as WORD WORD NOSPACE WORD NOSPACE *  etc.
// This is really just to make it a bit more human readable.
func (te *TokenizedLogEntry) Simplify() error {
	newTokenList := []tokenizers.DataType{}

	hasGenericToken := false
	for _, t := range te.Tokens {

		if t == tokenizers.NOTSPACE || t == tokenizers.WORD || t == tokenizers.ANYDATA {
			if !hasGenericToken {
				newTokenList = append(newTokenList, tokenizers.ANYDATA)
			}
			hasGenericToken = true
		} else {
			hasGenericToken = false
			newTokenList = append(newTokenList, t)
		}
	}

	te.Tokens = newTokenList

	return nil
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

// processCluster reads in a channel of clusters and returns via a channel of TokenizedLogEntries
func (lm *LogMine) processCluster(ch chan ClusterMessage, out chan TokenizedLogEntryMessage  ) {

	// need to check if this is the best way of doing this.
	for cm := range ch {
		cluster := cm.cluster
		tokenizedLogEntry, err := lm.clusterProcessor.ProcessSingleCluster(cluster)
		if err != nil {
			log.Errorf("Unable to process cluster %v : error %s\n", cluster, err.Error())
			continue
		}

		tlem := TokenizedLogEntryMessage{ tokenizedLogEntry: *tokenizedLogEntry, originalClusterIndex: cm.originalIndex}
		out <- tlem
	}
}


func (lm *LogMine) ProcessLogsFromReader(reader io.Reader, maxLevel int, noProcessors int) error {

	// preprocess + datatype identification
	tokenizedLogEntries, err := lm.Preprocess(reader)
	if err != nil {
		return err
	}

	ch := make(chan ClusterMessage, 100)
	resultCh := make(chan TokenizedLogEntryMessage, 100)

	for i:=0;i<noProcessors;i++ {
		go lm.processCluster(ch, resultCh)
	}

	fmt.Printf("Preprocessing complete\n")
	// loop through all the levels.
	for level := 0; level <= maxLevel; level++ {
		fmt.Printf("level %d\n", level)
		// generate clusters.
		err = lm.ClusterGeneration(tokenizedLogEntries, level)
		if err != nil {
			return err
		}

		newTokenizedLogEntries := []TokenizedLogEntry{}

		// put all clusters on channel to be read and processed by a higher power
		for index, cluster := range lm.clusterProcessor.clusters[level] {
			cm := ClusterMessage{ cluster:cluster, originalIndex: index}
			ch <- cm
		}

		// we know exactly how many responses we should get, just loop the
		// appropriate number of times. Make this more error proof! TODO(kpfaulkner)
		for index := 0 ; index < len(lm.clusterProcessor.clusters[level]); index++ {

			tlem := <- resultCh
			tokenizedLogEntry := tlem.tokenizedLogEntry

			cluster := lm.clusterProcessor.clusters[level][tlem.originalClusterIndex]

			// if first level (ie ALL logs available) then store number of logs.
			if level == 0 {
				tokenizedLogEntry.NumberOfPreviousEntries = len(cluster.logsInCluster)
			}

			// record pattern for cluster.
			cluster.PatternForCluster = tokenizedLogEntry
			lm.clusterProcessor.clusters[level][tlem.originalClusterIndex] = cluster
			newTokenizedLogEntries = append(newTokenizedLogEntries, tokenizedLogEntry)
		}

		tokenizedLogEntries = newTokenizedLogEntries
	}

	return nil
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

func (lm *LogMine) PreprocessLine(ch chan string, out chan TokenizedLogEntry) {

	for line := range ch {
		tokens, err := lm.tokenizer.Tokenize(line)
		if err != nil {
			log.Errorf("Error PreProcess %s\n", err.Error())
			continue
		}
		te := TokenizedLogEntry{Tokens: tokens}
		out <- te
	}
}


// Preprocess will read in ALL log entries from a file(reader)
// and process them. Will return a TokenizedLogEntry for each log line
// read.
func (lm *LogMine) Preprocess(reader io.Reader) ([]TokenizedLogEntry, error) {

	tokenizedLogEntries := []TokenizedLogEntry{}

	ch := make(chan string, 10000)
	resultCh := make(chan TokenizedLogEntry, 10000)

	wg := sync.WaitGroup{}

	for i:=0;i<10;i++ {
		go lm.PreprocessLine(ch, resultCh)
	}

	go func(){
	  for res := range resultCh {
		  tokenizedLogEntries = append(tokenizedLogEntries, res)
		  wg.Done()
	  }
	}()

	// read each log entry and preprocess them.
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := scanner.Text()
		if willProcessLine(l) {
			wg.Add(1)
      ch <- l
		}
	}

	fmt.Printf("waiting on preprocessing\n")
	wg.Wait()
	fmt.Printf("completeds preprocessing\n")

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

func (lm *LogMine) GenerateUsefulOutput(simplify bool) ([]TokenizedLogEntry, error) {

	lastLevel := len(lm.clusterProcessor.clusters) - 1

	clusters := lm.clusterProcessor.clusters[lastLevel]
	sort.Slice(clusters, func(i int, j int) bool {
		return clusters[i].PatternForCluster.NumberOfPreviousEntries > clusters[j].PatternForCluster.NumberOfPreviousEntries
	})

	results := make([]TokenizedLogEntry, len(clusters))
	for i, c := range lm.clusterProcessor.clusters[lastLevel] {
		results[i] = c.PatternForCluster
	}

	return results, nil
}

// just display to stdout for now.
// order by fewest entries to most.
func (lm *LogMine) DisplayFinalOutput(simplify bool) error {

	lastLevel := len(lm.clusterProcessor.clusters) - 1

	clusters := lm.clusterProcessor.clusters[lastLevel]
	sort.Slice(clusters, func(i int, j int) bool {
		return clusters[i].PatternForCluster.NumberOfPreviousEntries < clusters[j].PatternForCluster.NumberOfPreviousEntries
	})

	tokens,_ := lm.clusterProcessor.CreateSimplifedPatternForClusterLevel(lastLevel)
	sort.Slice(tokens, func(i int, j int) bool {
		return tokens[i].NumberOfPreviousEntries < tokens[j].NumberOfPreviousEntries
	})

	for _,t := range tokens {
		fmt.Printf("count %d : pattern %s\n", t.NumberOfPreviousEntries, t.ToString())
	}
	return nil
}
