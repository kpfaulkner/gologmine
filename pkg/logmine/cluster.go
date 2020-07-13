package logmine

import (
	"github.com/kpfaulkner/gologmine/pkg/logmine/tokenizers"
	"math"
)

type Cluster struct {
	score         float64 // floats ok for this?
	logsInCluster []TokenizedLogEntry
	PatternForCluster TokenizedLogEntry
}

type ClusterGeneration struct {
	Count int // number of entries in this cluster generation
}

// Rules for clustering.
// This is taken from the powerpoint.
// Unsure what --- is???
// Think it's from the smither & waterman algorithm used to
// fill in positions where the 2 lists are different lengths
//
// DATETIME + NOTSPACE = *
// WORD + NOTSPACE = NOTSPACE
// WORD + NUMBER = NOTSPACE
// IPV4 + WORD = NOTSPACE
// IPV4 + NUMBER = NOTSPACE
// IPV4 = DATETIME = *
// --- + WORD = *
// --- + NUMBER = *
// --- + NOTSPACE = *
// --- + DATETIME = *
type ClusterProcessor struct {

	// list of list of clusters.
	// first index is level dealing with, second index are all the clusters in that level.
	// clusters[0][2] is level 0 clusters... and getting the 3rd of that level.
	clusters [][]Cluster

	//MaxDistance float64
	distances []float64
}

func NewClusterProcessor(distances []float64) ClusterProcessor {
	c := ClusterProcessor{}
	//c.MaxDistance = 0.01 // just going off presentation  for now. Will need to figure this out.
	c.distances = distances
	return c
}

func (cp *ClusterProcessor) ClearXX() {
  //cp.clusters = []Cluster{}
}

func score(e1 tokenizers.DataType, e2 tokenizers.DataType, level int) float64 {

	if level == 0 {
		if e1 == e2 {
			return 1
		}
		return 0
	}

	// if a generic data type, return 1.
  newToken := ConvertTokenDataType(e1, e2)
  if newToken != "" {
  	return 1
  }

  // otherwise do normal test.
  // This is just testing out an idea.
	if e1 == e2 {
		return 1
	}
	return 0

}

func LogDistance(log1 TokenizedLogEntry, log2 TokenizedLogEntry, level int) float64 {

	log1Len := float64(len(log1.Tokens))
	log2Len := float64(len(log2.Tokens))
	maxLen := math.Max(log1Len, log2Len)
	minLen := math.Min(log1Len, log2Len)

	total := 0.0
	for i := 0; i < int(minLen); i++ {
		s := score(log1.Tokens[i], log2.Tokens[i], level) / maxLen
		total += s
	}

	return 1 - total
}

func (cp *ClusterProcessor) AddLogEntry(l TokenizedLogEntry, level int) error {

	addedToCluster := false
	indexOfClosestCluster := -1
  closestDistance := 100.0

  // add empty cluster list if required
  if len(cp.clusters) < level +1 {
  	cp.clusters = append(cp.clusters, []Cluster{})
  }

	// calculate which cluster it can go into.
	for index, cluster := range cp.clusters[level] {

		// just get distance between new log entry and first element in cluster.
		// just use first entry for appropriate level.
		dist := LogDistance(cluster.logsInCluster[0], l, level)

		if dist <= cp.distances[level] && dist <= closestDistance {
			indexOfClosestCluster = index
			closestDistance = dist
			addedToCluster = true
		}
	}

	// haven't added to cluster yet, so make a new one.
	if !addedToCluster {
		c := Cluster{}
		c.logsInCluster = append(c.logsInCluster, l)
		cp.clusters[level] = append(cp.clusters[level], c)
	} else {
		cp.clusters[level][indexOfClosestCluster].logsInCluster = append(cp.clusters[level][indexOfClosestCluster].logsInCluster, l)
	}

	return nil
}

func ConvertTokenDataType( token1 tokenizers.DataType, token2 tokenizers.DataType) tokenizers.DataType {

	tokenToUse := tokenizers.DataType("")

	if token1 == tokenizers.DATE && token2 == tokenizers.NOTSPACE {
		tokenToUse = tokenizers.ANYDATA
	}

	if token1 == tokenizers.WORD && token2 == tokenizers.NOTSPACE {
		tokenToUse = tokenizers.NOTSPACE
	}

	if token1 == tokenizers.WORD && token2 == tokenizers.NUMBER{
		tokenToUse = tokenizers.NOTSPACE
	}

	if token1 == tokenizers.IPV4 && token2 == tokenizers.WORD{
		tokenToUse = tokenizers.NOTSPACE
	}

	if token1 == tokenizers.IPV4 && token2 == tokenizers.NUMBER{
		tokenToUse = tokenizers.NOTSPACE
	}

	if token1 == tokenizers.IPV4 && token2 == tokenizers.NUMBER{
		tokenToUse = tokenizers.NOTSPACE
	}

	if token1 == tokenizers.IPV4 && token2 == tokenizers.DATE{
		tokenToUse = tokenizers.NOTSPACE
	}

	if token1 == tokenizers.IPV4 && token2 == tokenizers.TIME{
		tokenToUse = tokenizers.NOTSPACE
	}

	if token1 == tokenizers.ALIGNER && token2 == tokenizers.WORD{
		tokenToUse = tokenizers.ANYDATA
	}

	if token1 == tokenizers.ALIGNER && token2 == tokenizers.NUMBER{
		tokenToUse = tokenizers.ANYDATA
	}

	if token1 == tokenizers.ALIGNER && token2 == tokenizers.NOTSPACE{
		tokenToUse = tokenizers.ANYDATA
	}

	if token1 == tokenizers.ALIGNER && token2 == tokenizers.DATE{
		tokenToUse = tokenizers.ANYDATA
	}

	if token1 == tokenizers.ALIGNER && token2 == tokenizers.TIME{
		tokenToUse = tokenizers.ANYDATA
	}

	return tokenToUse
}

func mergeToken(t1 tokenizers.DataType, t2 tokenizers.DataType, e1 tokenizers.DataType, e2 tokenizers.DataType, replacementToken tokenizers.DataType, existingToken tokenizers.DataType) tokenizers.DataType {

	if t1 == e1 && t2 == e2 ||
		 t1 == e2 && t2 == e1 {
		return replacementToken
	}

	// special case for aligners.... always produce an ANY data entry?
	if t1 == tokenizers.ALIGNER || t2 == tokenizers.ALIGNER {
		return tokenizers.ANYDATA
	}

	return existingToken
}

// mergeAlignedLogs 2 log entries should have same length (aligned)
// now go through rules to determine what the merged version looks like:
func mergeAlignedLogs( align1 []tokenizers.DataType, align2 []tokenizers.DataType) ([]tokenizers.DataType, error) {

	result := make([]tokenizers.DataType, len(align1))

	var tokenToUse tokenizers.DataType
	for i:=0;i<len(align1);i++ {
		token1 := align1[i]
		token2 := align2[i]

		// default to word? does this make ANY sense?
		tokenToUse = tokenizers.WORD

		if token1 == token2 {
			tokenToUse = token1
		}

		tokenToUse = mergeToken(token1, token2, tokenizers.DATE, tokenizers.NOTSPACE, tokenizers.ANYDATA, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.WORD, tokenizers.NOTSPACE, tokenizers.NOTSPACE, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.WORD, tokenizers.NUMBER, tokenizers.NOTSPACE, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.IPV4, tokenizers.WORD, tokenizers.NOTSPACE, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.IPV4, tokenizers.NUMBER, tokenizers.NOTSPACE, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.IPV4, tokenizers.DATE, tokenizers.ANYDATA, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.IPV4, tokenizers.TIME, tokenizers.ANYDATA, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.ALIGNER, tokenizers.WORD, tokenizers.ANYDATA, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.ALIGNER, tokenizers.NUMBER, tokenizers.ANYDATA, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.ALIGNER, tokenizers.NOTSPACE, tokenizers.ANYDATA, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.ALIGNER, tokenizers.DATE, tokenizers.ANYDATA, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.ALIGNER, tokenizers.TIME, tokenizers.ANYDATA, tokenToUse)
		tokenToUse = mergeToken(token1, token2, tokenizers.ALIGNER, tokenizers.ANYDATA, tokenizers.ANYDATA, tokenToUse)


		result[i] = tokenToUse

	}

	return result, nil
}

// Process a cluster (collection of TokenizedLogEntry) and generate a generic
// TokenizedLogEntry that will represent the entire cluster.
func (cp *ClusterProcessor) ProcessSingleCluster( cluster Cluster) (*TokenizedLogEntry, error) {

	existingEntry := cluster.logsInCluster[0].Tokens
	numPreviousLogs := cluster.logsInCluster[0].NumberOfPreviousEntries
	for _,entry := range cluster.logsInCluster[1:] {

		// align the 2 logs.
		align1, align2, err := SmithWaterman(existingEntry, entry.Tokens)
		if err != nil {
			return nil, err
		}

		// merge the alignments
		mergedResult, err := mergeAlignedLogs( align1, align2)
		if err != nil {
			return nil, err
		}

		existingEntry = mergedResult

		// keep count of all logs that make up this cluster
		numPreviousLogs += entry.NumberOfPreviousEntries
	}

  tle := TokenizedLogEntry{}
  tle.Tokens = existingEntry
  tle.NumberOfPreviousEntries = numPreviousLogs
  return &tle,nil
}

// SimplifyClusterLevel consolidates contiguous *,WORD or NOSPACE together
// This is purely for display purposes and would mean any further processing would be null (nil)
// and void :P
func (cp *ClusterProcessor) CreateSimplifedPatternForClusterLevel(level int) ([]TokenizedLogEntry,error) {

	existingClusters := cp.clusters[level]

	newTokenizedLogEntrySlice := []TokenizedLogEntry{}
	hasGenericToken := false
	for _, t := range existingClusters{
		pattern := t.PatternForCluster
		newTokenSlice := []tokenizers.DataType{}
		for _,t := range pattern.Tokens {
			if t == tokenizers.NOTSPACE || t == tokenizers.WORD || t == tokenizers.ANYDATA {
				if !hasGenericToken {
					newTokenSlice = append(newTokenSlice, tokenizers.ANYDATA)
				}
				hasGenericToken = true
			} else {
				hasGenericToken = false
				newTokenSlice = append(newTokenSlice, t)
			}
		}

		tle := TokenizedLogEntry{}
		tle.Tokens = newTokenSlice
		tle.NumberOfPreviousEntries = pattern.NumberOfPreviousEntries
		newTokenizedLogEntrySlice = append(newTokenizedLogEntrySlice, tle)
	}

	// have single slice with simplified entries.
	/// remove dupes (but increase NumberOfPreviousEntries...
	newTokenizedLogEntrySlice = mergeSimplifiedTokens(newTokenizedLogEntrySlice)

  return newTokenizedLogEntrySlice, nil
}


func mergeSimplifiedTokens(logEntries []TokenizedLogEntry) []TokenizedLogEntry {

	filteredSlice := []TokenizedLogEntry{}
  tokenLUT := make(map[string]TokenizedLogEntry)

	for _,le := range logEntries {
		entry, ok := tokenLUT[le.ToString()]
		if ok {
			entry.NumberOfPreviousEntries += le.NumberOfPreviousEntries
		} else {
			entry = le
		}
		tokenLUT[le.ToString()] = entry
	}

	for _,v := range tokenLUT {
		filteredSlice = append(filteredSlice, v)
	}

	return filteredSlice
}
