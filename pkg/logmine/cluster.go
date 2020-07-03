package logmine

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
type Cluster struct {


}


func NewCluster() Cluster {
	c := Cluster{}

	return c
}

