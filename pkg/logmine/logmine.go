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
	"time"
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

func (lm *LogMine) ProcessLogsFromSlice(logEntries []string, maxLevel int) error {
	// preprocess + datatype identification
	tokenizedLogEntries, err := lm.PreprocessFromSlice(logEntries)
	if err != nil {
		return err
	}

	err = lm.processTokenizedLogEntries(tokenizedLogEntries, maxLevel)
	return err
}

// processLogsFromChannel batch process from channel.
// Read as many entries (up to a limit or timeout) from the channel
// and then send them for processing.
// Don't want to process single tle at a time.
func (lm *LogMine) processLogsFromChannelORIG(outCh chan TokenizedLogEntry, maxLevel int) error {

	done := false
	batchDone := false

	// while counter > 0 keep processing. if gets to 0... no content for a while... bail
	counter := 5

	for !done && counter > 0 {

		tokens := []TokenizedLogEntry{}

		for !batchDone {
			select {
			case te := <-outCh:
				counter = 5
				tokens = append(tokens, te)
				if len(tokens) > 100 {
					batchDone = true
				}
			case <-time.After(5 * time.Second):
				batchDone = true
				counter--
			}
		}

		if len(tokens) > 0 {
			err := lm.processTokenizedLogEntries(tokens, maxLevel)
			if err != nil {
				log.Errorf("Unable to process log batch : %s", err.Error())
			}
		}
		batchDone = false
	}

	return nil
}

// processLogsFromChannel batch process from channel.
// Read as many entries (up to a limit or timeout) from the channel
// and then send them for processing.
// Don't want to process single tle at a time.
func (lm *LogMine) processLogsFromChannel(outCh chan TokenizedLogEntry, maxLevel int) error {

	tokens := []TokenizedLogEntry{}
	totalTokensObserved := 0

	for te := range outCh {
		tokens = append(tokens, te)
		totalTokensObserved++
		if len(tokens) > 1000 {
			err := lm.generateAllClusterLevels(tokens, maxLevel)
			if err != nil {
				log.Errorf("Unable to process log batch : %s", err.Error())
			}
			tokens = []TokenizedLogEntry{}
			//fmt.Printf("number of tokens %d\n", totalTokensObserved)
		}
	}

	// catch anything thats been put in the slice but hasn't been processed yet.
	if len(tokens) > 1000 {
		err := lm.generateAllClusterLevels(tokens, maxLevel)
		if err != nil {
			log.Errorf("Unable to process log batch : %s", err.Error())
		}
	}

	// now
	return nil
}

// ProcessLogsAsRead reads the log entries from the file, sends to
// another function for tokenization and gets the results back as soon as possible
// ie it does not have to wait for entire file to be processed.
// It then passes the tokenized version to the LogMine routines.
func (lm *LogMine) ProcessLogsAsRead(reader io.Reader, maxLevel int) error {

	outCh := make(chan TokenizedLogEntry, 100000)
	inCh := make(chan string, 100000)

	// should we spawn more of these?
	// Currently this is definitely NOT the bottleneck.
	go lm.processLogsFromChannel(outCh, maxLevel)

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go lm.doTokenization(inCh, outCh, &wg)
	}

	// read each log entry and preprocess them.
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := scanner.Text()
		if willProcessLine(l) {
			inCh <- l
		}
	}
	close(inCh)
	wg.Wait()
	close(outCh)


	return nil
}

func (lm *LogMine) ProcessLogsFromReader(reader io.Reader, maxLevel int) error {

	// preprocess + datatype identification
	tokenizedLogEntries, err := lm.PreprocessFromReader(reader)
	if err != nil {
		return err
	}

	fmt.Printf("tokenize: %s  :Read file and initial tokenization\n", time.Now())
	err = lm.processTokenizedLogEntries(tokenizedLogEntries, maxLevel)
	return err
}

func (lm *LogMine) generateAllClusterLevels(tokenizedLogEntries []TokenizedLogEntry, maxLevel int) error {

	start1 := time.Now()
	// loop through all the levels.
	for level := 0; level <= maxLevel; level++ {
		fmt.Printf("level %d\n", level)
		// generate clusters.
		err := lm.ClusterGeneration(tokenizedLogEntries, level)
		if err != nil {
			return err
		}
	}
	end1 := time.Now()
	fmt.Printf("generateAllClusterLevels took %d ms\n", end1.Sub(start1).Milliseconds())
	return nil
}

func (lm *LogMine) processTokenizedLogEntries(tokenizedLogEntries []TokenizedLogEntry, maxLevel int) error {

	start1 := time.Now()
	// loop through all the levels.
	for level := 0; level <= maxLevel; level++ {
		fmt.Printf("level %d\n", level)
		// generate clusters.
		err := lm.ClusterGeneration(tokenizedLogEntries, level)
		if err != nil {
			return err
		}

		newTokenizedLogEntries := []TokenizedLogEntry{}

		start := time.Now()
		// now process/merge each cluster in the cluster "level"
		for index, cluster := range lm.clusterProcessor.clusters[level] {
			tokenizedLogEntry, err := lm.clusterProcessor.ProcessSingleCluster(cluster)
			if err != nil {
				return err
			}

			// if first level (ie ALL logs available) then store number of logs.
			if level == 0 {
				tokenizedLogEntry.NumberOfPreviousEntries = len(cluster.logsInCluster)
			}

			// record pattern for cluster.
			cluster.PatternForCluster = *tokenizedLogEntry
			lm.clusterProcessor.clusters[level][index] = cluster
			newTokenizedLogEntries = append(newTokenizedLogEntries, *tokenizedLogEntry)
		}
		tokenizedLogEntries = newTokenizedLogEntries
		end := time.Now()
		fmt.Printf("processing cluster level took %d ms\n", end.Sub(start).Milliseconds())
	}

	end1 := time.Now()
	fmt.Printf("processTokenizedLogEntries took %d ms\n", end1.Sub(start1).Milliseconds())
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

// PreprocessFromReader will read in ALL log entries from a file(reader)
// and process them. Will return a TokenizedLogEntry for each log line
// read.
func (lm *LogMine) PreprocessFromReader(reader io.Reader) ([]TokenizedLogEntry, error) {

	tokenizedLogEntries := []TokenizedLogEntry{}
	outCh := make(chan TokenizedLogEntry, 100000)
	inCh := make(chan string, 100000)
	go func() {
		for te := range outCh {
			tokenizedLogEntries = append(tokenizedLogEntries, te)
		}
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go lm.doTokenization(inCh, outCh, &wg)
	}

	// read each log entry and preprocess them.
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := scanner.Text()
		if willProcessLine(l) {
			inCh <- l
		}
	}
	close(inCh)
	wg.Wait()
	close(outCh)
	return tokenizedLogEntries, nil
}

func (lm *LogMine) doTokenization(inCh chan string, outCh chan TokenizedLogEntry, wg *sync.WaitGroup) {
	for l := range inCh {
		tokens, err := lm.tokenizer.Tokenize(l)
		if err != nil {
			log.Errorf("Error PreProcess %s\n", err.Error())
			continue
		}
		te := TokenizedLogEntry{Tokens: tokens}
		outCh <- te
	}
	wg.Done()
}

// PreprocessFromSlice will read in ALL log entries from a file(reader)
// and process them. Will return a TokenizedLogEntry for each log line
// read.
func (lm *LogMine) PreprocessFromSlice(logEntries []string) ([]TokenizedLogEntry, error) {

	tokenizedLogEntries := []TokenizedLogEntry{}

	// read each log entry and preprocess them.
	for _, l := range logEntries {
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

	tokens, err := lm.GenerateFinalOutput(simplify)
	if err != nil {
		return err
	}

	for _, t := range tokens {
		fmt.Printf("count %d : pattern %s\n", t.NumberOfPreviousEntries, t.ToString())
	}
	return nil
}

func (lm *LogMine) GenerateFinalOutput(simplify bool) ([]TokenizedLogEntry, error) {

	lastLevel := len(lm.clusterProcessor.clusters) - 1

	clusters := lm.clusterProcessor.clusters[lastLevel]
	sort.Slice(clusters, func(i int, j int) bool {
		return clusters[i].PatternForCluster.NumberOfPreviousEntries < clusters[j].PatternForCluster.NumberOfPreviousEntries
	})

	tokens, _ := lm.clusterProcessor.CreateSimplifedPatternForClusterLevel(lastLevel)
	sort.Slice(tokens, func(i int, j int) bool {
		return tokens[i].NumberOfPreviousEntries < tokens[j].NumberOfPreviousEntries
	})

	return tokens, nil
}
