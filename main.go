package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	t1 time.Time
	t2 time.Time
)

type calcFunc func(numbers []int) (uint64, []int)

func main() {
	numbersStr := strings.Split(os.Args[1], "-")
	if len(numbersStr) != 6 {
		log.Fatalf("Precisam ser 6 numeros: %d", len(numbersStr))
	}
	numbers := make([]int, len(numbersStr))
	for i, v := range numbersStr {
		numbers[i], _ = strconv.Atoi(v)
	}

	sort.Ints(numbers)

	var function = getCalcFunction()

	t1 = time.Now() // inicio do temporizador
	tentativas, result := function(numbers)
	t2 = time.Now() // fim do temporizador

	fmt.Printf("%d-%d-%d-%d-%d-%d\n", result[0], result[1], result[2], result[3], result[4], result[5])

	timeElapsed := t2.Sub(t1).Milliseconds()
	ops := float64(tentativas) / float64(timeElapsed)

	fmt.Printf("Tentativas: %d - Tempo: %dms - Ops/ms: %f", tentativas, timeElapsed, ops)
}

func getCalcFunction() calcFunc {
	switch os.Getenv("ALG") {
	case "brute":
		return brute
	case "xorshift":
		return xorshift
	default:
		return xorshift
	}
}
