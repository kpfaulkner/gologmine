package main

import (
	"flag"
	"fmt"
	"github.com/kpfaulkner/gologmine/pkg/logmine"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

// hopefully this ends up as a working implementation of LogMine
// See https://www.cs.unm.edu/~mueen/Papers/LogMine.pdf and https://www.youtube.com/watch?v=URQnTPTxBbA for details.
func testTokenizer() {

	tokenizer := logmine.NewTokenizer()

	res, _ := tokenizer.Tokenize("2017/02/04 09:01:00 login 127.0.0.1 user=bear12")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:02:00 DB Connect 127.0.0.1 user=bear12")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:03:00 DB Disconnect 127.0.0.1 user=bear12")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:04:00 logout 127.0.0.1 user=bear12")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:05:00 login 127.0.0.1 user=bear34")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:06:00 DB Connect 127.0.0.1 user=bear34")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:07:00 DB Disconnect 127.0.0.1 user=bear34")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:08:00 logout 127.0.0.1 user=bear34")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:05:00 login 127.0.0.1 user=bear#1")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:06:00 DB Connect 127.0.0.1 user=bear#1")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:07:00 DB Disconnect 127.0.0.1 user=bear#1")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:08:00 logout 127.0.0.1 user=bear#1")
	fmt.Printf("res is %v\n", res)
}

func generateMaxDistances( maxDist string) ([]float64, error) {
	sp := strings.Split(maxDist,",")
	distances := []float64{}

  for _,d := range sp {
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

	maxDist := flag.String("maxdist", "0.01,0.1", "Max distances for clustering. Comma separated decimals. eg 0.01,0.05 etc")
	file := flag.String("file", "test.log", "Log file to process")
	flag.Parse()

	distances, err := generateMaxDistances(*maxDist)

	lm := logmine.NewLogMine(distances)
  f, _ := os.Open(*file)

  res, err := lm.ProcessLogsFromReader(f)
  if err != nil {
  	log.Fatalf("error while processing. %s\n", err.Error())
  }

  fmt.Printf("======================\n")
  for _,e := range res {
		fmt.Printf("%v\n", e.Tokens)
  }

  res, err = lm.ProcessAgain(res)
	fmt.Printf("======================\n")
	for _,e := range res {
		fmt.Printf("%v\n", e.Tokens)
	}

}
