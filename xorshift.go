package main

import (
	"sort"
	"sync/atomic"
)

// seedNumber é o intervalor do seed entre os workers
const seedNumber = 5000000

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

func worker(done *int32, workerNumber int, numbersSum int, tentativas *uint64, numbers []int, res []int, sig chan int) {
	calc := NewCalc(int64(workerNumber * seedNumber))

	result := make([]byte, 6)
	var tries uint64

	for {
		if *done > 0 {
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
				*tentativas = tries

				sig <- 1
			}
		}
		tries++
	}
}

func xorshift(numbers []int) (uint64, []int) {
	numbersSum := 0
	for i, _ := range numbers {
		numbersSum += numbers[i]
	}

	maxProc := 5
	sig := make(chan int, 1)
	tentativas := make([]uint64, maxProc)

	var done int32 = 0
	result := make([]int, 6)
	for i := 0; i < maxProc; i++ {
		go worker(&done, i, numbersSum, &tentativas[i], numbers, result, sig)
	}
	<-sig
	atomic.AddInt32(&done, 1)
	var totalTentativas uint64
	for _, v := range tentativas {
		totalTentativas += v
	}
	return totalTentativas, result
}
