package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// seedNumber é o intervalor do seed entre os workers
const seedNumber = 5000000

var (
	t1   time.Time
	t2   time.Time
	done int32 = 0
)

type randctx struct {
	a uint64
	b uint64
}

// xorshift+
func (c *randctx) nextrand() uint64 {
	t := c.a
	s := c.b

	c.a = s
	t ^= t << 23
	t ^= t >> 17
	t ^= s ^ (s >> 26)
	c.b = t
	return t + s
}

func (c *randctx) nextval() byte {
	return byte(c.nextrand())
}

func NewCalc(seed int64) *randctx {
	return &randctx{
		a: uint64(1 + seed),
		b: uint64(1 + seed),
	}
}

func worker(workerNumber int, numbersSum int, tentativas *int, numbers []int, res []int, sig chan int) {
	calc := NewCalc(int64(workerNumber * seedNumber))

	result := make([]byte, 6)

	for {
		if done > 0 {
			break
		}
		checksum := 0
		for i := 0; i < 6; i++ {

			// obter novo numero sorteado, sem repetir os anteriores
			for {
				goback := false
				result[i] = calc.nextval() % 61
				for j := 0; j < i; j++ {
					if result[i] == result[j] {
						goback = true
					}
				}
				if !goback {
					break
				}
			}
			checksum += int(result[i])
		}

		if checksum == numbersSum {
			// probabilidade de serem os numeros escolhidos
			verify := make([]int, 6)
			for i, v := range result {
				verify[i] = int(v)
			}
			sort.Ints(verify)
			if verify[0] == numbers[0] &&
				verify[1] == numbers[1] &&
				verify[2] == numbers[2] &&
				verify[3] == numbers[3] &&
				verify[4] == numbers[4] &&
				verify[5] == numbers[5] {

				for i, v := range result {
					res[i] = int(v)
				}

				sig <- 1
			}
		}
		*tentativas++
	}
}

func main() {
	numbersStr := strings.Split(os.Args[1], "-")
	if len(numbersStr) != 6 {
		log.Fatalf("Precisam ser 6 numeros: %d", len(numbersStr))
	}
	numbers := make([]int, len(numbersStr))
	numbersSum := 0
	for i, v := range numbersStr {
		numbers[i], _ = strconv.Atoi(v)
		numbersSum += numbers[i]
	}
	sort.Ints(numbers)

	maxProc := 5
	sig := make(chan int, 1)
	tentativas := make([]int, maxProc)

	result := make([]int, 6)
	t1 = time.Now() // inicio do temporizador
	for i := 0; i < maxProc; i++ {
		go worker(i, numbersSum, &tentativas[i], numbers, result, sig)
	}
	<-sig
	atomic.AddInt32(&done, 1)
	t2 = time.Now() // fim do temporizador

	fmt.Printf("%d-%d-%d-%d-%d-%d\n", result[0], result[1], result[2], result[3], result[4], result[5])

	totalTentativas := 0
	for _, v := range tentativas {
		totalTentativas += v
	}

	timeElapsed := t2.Sub(t1).Milliseconds()
	ops := float64(totalTentativas) / float64(timeElapsed)

	fmt.Printf("Tentativas: %d - Tempo: %dms - Ops/ms: %f", totalTentativas, timeElapsed, ops)
}
