package logmine

import (
	"fmt"
	"github.com/kpfaulkner/gologmine/pkg/logmine/tokenizers"
	"math"
)

type Cluster struct {
	score         float64 // floats ok for this?
	logsInCluster []TokenizedLogEntry
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
	clusters []Cluster

	//MaxDistance float64
	distances []float64
}

func NewClusterProcessor(distances []float64) ClusterProcessor {
	c := ClusterProcessor{}
	//c.MaxDistance = 0.01 // just going off presentation  for now. Will need to figure this out.
	c.distances = distances
	return c
}

func (cp *ClusterProcessor) Clear() {
  cp.clusters = []Cluster{}
}

func score(e1 tokenizers.DataType, e2 tokenizers.DataType) float64 {
	if e1 == e2 {
		return 1
	}
	return 0
}

func LogDistance(log1 TokenizedLogEntry, log2 TokenizedLogEntry) float64 {

	log1Len := float64(len(log1.Tokens))
	log2Len := float64(len(log2.Tokens))
	maxLen := math.Max(log1Len, log2Len)
	minLen := math.Min(log1Len, log2Len)

	total := 0.0
	for i := 0; i < int(minLen); i++ {
		s := score(log1.Tokens[i], log2.Tokens[i]) / maxLen
		total += s
	}

	return 1 - total
}

func (cp *ClusterProcessor) AddLogEntry(l TokenizedLogEntry, level int) error {

	addedToCluster := false
	indexOfClosestCluster := -1
  closestDistance := 100.0

  if l.Tokens[6] == tokenizers.NOTSPACE {
  	fmt.Printf("boom\n")
  }

	// calculate which cluster it can go into.
	for index, cluster := range cp.clusters {
		// just get distance between new log entry and first element in cluster.
		dist := LogDistance(cluster.logsInCluster[0], l)

		if dist <= cp.distances[level] && dist <= closestDistance {
			indexOfClosestCluster = index
			closestDistance = dist
			addedToCluster = true
		}

		/*
		// add to first cluster that meets criteria <--- mistake I think.
		if dist <= cp.distances[level] {
			indexOfClosestCluster = index
			cp.clusters[index].logsInCluster = append(cp.clusters[index].logsInCluster, l)
			addedToCluster = true
			break
		} */

	}

	// haven't added to cluster yet, so make a new one.
	if !addedToCluster {
		c := Cluster{}
		c.logsInCluster = append(c.logsInCluster, l)
		cp.clusters = append(cp.clusters, c)
	} else {
		cp.clusters[indexOfClosestCluster].logsInCluster = append(cp.clusters[indexOfClosestCluster].logsInCluster, l)
	}

	return nil
}

// mergeAlignedLogs 2 log entries should have same length (aligned)
// now go through rules to determine what the merged version looks like:
func mergeAlignedLogs( align1 []tokenizers.DataType, align2 []tokenizers.DataType) ([]tokenizers.DataType, error) {

	result := make([]tokenizers.DataType, len(align1))

	var tokenToUse tokenizers.DataType
	for i:=0;i<len(align1);i++ {
		token1 := align1[i]
		token2 := align2[i]

		if token1 == "Connect" {
		  fmt.Printf("oops\n")
		}

		if token1 == token2 {
			tokenToUse = token1
		}

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

		result[i] = tokenToUse

	}

	return result, nil
}

// Process a cluster (collection of TokenizedLogEntry) and generate a generic
// TokenizedLogEntry that will represent the entire cluster.
func (cp *ClusterProcessor) ProcessSingleCluster( cluster Cluster) (*TokenizedLogEntry, error) {

	existingEntry := cluster.logsInCluster[0].Tokens
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
	}

  tle := TokenizedLogEntry{}
  tle.Tokens = existingEntry
  return &tle,nil
}
