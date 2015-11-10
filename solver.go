package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"
)

var cores int
var solved chan string

type solutionAtN struct {
	n       int
	count   int
	bitmask int
	started time.Time
}

func makeSolution(n int) *solutionAtN {
	bitmask := int(math.Pow(2, float64(n)) - 1)
	newS := solutionAtN{n: n, started: time.Now(), bitmask: bitmask}
	return &newS
}

func (s *solutionAtN) Solve() {
	s.try(0, 0, 0)
	solved <- fmt.Sprintf("There are %5d solutions to %2d queens problem, found in %s\n", s.count, s.n, time.Since(s.started))
}

func (s *solutionAtN) try(ld int, cols int, rd int) {
	if cols == s.bitmask {
		s.count++
		return
	}
	poss := ^(ld | cols | rd) & s.bitmask
	for poss != 0 {
		bit := poss & -poss
		poss -= bit
		s.try((ld|bit)<<1, (cols | bit), (rd|bit)>>1)
	}
}

func main() {
	maxN, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	solved = make(chan string, maxN)

	cores := runtime.NumCPU()
	workers := make(chan struct{}, cores)
	start := time.Now()

	go func() {
		for i := 1; i <= maxN; i++ {
			workers <- struct{}{}
			go makeSolution(i).Solve()
		}
	}()

	for i := 1; i <= maxN; i++ {
		solution := <-solved
		<-workers
		fmt.Println(solution)
	}

	fmt.Printf("Found all solution counts in %s", time.Since(start))
}
