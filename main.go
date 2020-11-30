package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kpfaulkner/gologmine/pkg/logmine"
	log "github.com/sirupsen/logrus"
)

func generateMaxDistances(maxDist string) ([]float64, error) {
	sp := strings.Split(maxDist, ",")
	distances := []float64{}

	for _, d := range sp {
		f, err := strconv.ParseFloat(d, 64)
		if err != nil {
			return nil, err
		}
		distances = append(distances, f)
	}

	return distances, nil
}

func main() {
	fmt.Printf("so it begins...\n")

	start := time.Now()
	fmt.Printf("start %s\n", start)
	maxDist := flag.String("maxdist", "0.01,0.1,0.3,0.9", "Max distances for clustering. Comma separated decimals. eg 0.01,0.05 etc")
	maxLevel := flag.Int("maxlevel", 3, "Max level 0-3")
	file := flag.String("file", "test.log", "Log file to process")
	simplify := flag.Bool("simplify", false, "Simplyify output. Compact sequential WORD, NOSPACE and * outputs to single entry")
	flag.Parse()

	distances, err := generateMaxDistances(*maxDist)

	lm := logmine.NewLogMine(distances)
	f, _ := os.Open(*file)

	err = lm.ProcessLogsFromReader(f, *maxLevel)
	if err != nil {
		log.Fatalf("error while processing. %s\n", err.Error())
	}

	//lm.DisplayFinalOutput(false)

	lm.DisplayFinalOutput(*simplify)

  end := time.Now()
	fmt.Printf("end %s\n", end)
  fmt.Printf("took %dms\n", end.Sub(start).Milliseconds())


}
