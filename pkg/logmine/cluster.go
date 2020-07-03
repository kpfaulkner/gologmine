package logmine

import (
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

	maxDistance float64
}

func NewClusterProcessor() ClusterProcessor {
	c := ClusterProcessor{}
	c.maxDistance = 0.01 // just going off presentation  for now. Will need to figure this out.
	return c
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

func (cp *ClusterProcessor) AddLogEntry(l TokenizedLogEntry) error {

	addedToCluster := false
	// calculate which cluster it can go into.
	for _, cluster := range cp.clusters {
		// just get distance between new log entry and first element in cluster.
		dist := LogDistance(cluster.logsInCluster[0], l)

		// add to first cluster that meets criteria
		if dist <= cp.maxDistance {
			cluster.logsInCluster = append(cluster.logsInCluster, l)
			addedToCluster = true
			break
		}
	}

	// haven't added to cluster yet, so make a new one.
	if !addedToCluster {
		c := Cluster{}
		c.logsInCluster = append(c.logsInCluster, l)
		cp.clusters = append(cp.clusters, c)
	}

	return nil
}
