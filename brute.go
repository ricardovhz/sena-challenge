package main

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

func bruteWorker(wg *sync.WaitGroup, done *int32, start int, end int, numbers []int, numbersSum int, tentativas *uint64, result []int) {
	var checkA, checkB, checkC, checkD, checkE, checkF int
	var a, b, c, d, e, f int
	for a = start; a < end; a++ {
		checkA = a
		for b = 0; b < 61; b++ {
			checkB = checkA + b
			for c = 0; c < 61; c++ {
				checkC = checkB + c
				for d = 0; d < 61; d++ {
					checkD = checkC + d
					for e = 0; e < 61; e++ {
						checkE = checkD + e
						for f = 0; f < 61; f++ {
							if *done > 0 {
								*tentativas = calcTries(a-start, b, c, d, e, f)
								wg.Done()
								return
							}
							checkF = checkE + f

							if checkF == numbersSum {

								// probabilidade de serem os numeros escolhidos
								verify := []int{
									a, b, c, d, e, f,
								}
								sort.Ints(verify)
								if verify[0] == numbers[0] &&
									verify[1] == numbers[1] &&
									verify[2] == numbers[2] &&
									verify[3] == numbers[3] &&
									verify[4] == numbers[4] &&
									verify[5] == numbers[5] {

									result[0] = a
									result[1] = b
									result[2] = c
									result[3] = d
									result[4] = e
									result[5] = f
									atomic.AddInt32(done, 1)
									*tentativas = calcTries(a-start, b, c, d, e, f)
									wg.Done()
									return
								}
							}
						}
					}
				}
			}
		}
	}
}

// O calculo de tentativas por forca bruta pode ser definido como
// (a * (61^5)) + (b * (61^4)) + (c * (61^3)) + (d * 61^2) + (e * 61) + f
func calcTries(a int, b int, c int, d int, e int, f int) uint64 {
	var (
		a1 = uint64(a)
		b1 = uint64(b)
		c1 = uint64(c)
		d1 = uint64(d)
		e1 = uint64(e)
		f1 = uint64(f)
	)
	var result uint64

	result += a1 * uint64(math.Pow(61, 5))
	result += b1 * uint64(math.Pow(61, 4))
	result += c1 * uint64(math.Pow(61, 3))
	result += d1 * uint64(math.Pow(61, 2))
	result += e1 * 61
	result += f1
	return result
}

func brute(numbers []int) (uint64, []int) {
	numbersSum := 0
	for _, v := range numbers {
		numbersSum += v
	}
	result := make([]int, 6)
	tentativas := make([]uint64, 6)
	var done int32
	wg := sync.WaitGroup{}
	wg.Add(6)

	go bruteWorker(&wg, &done, 0, 10, numbers, numbersSum, &tentativas[0], result)
	go bruteWorker(&wg, &done, 10, 20, numbers, numbersSum, &tentativas[1], result)
	go bruteWorker(&wg, &done, 20, 30, numbers, numbersSum, &tentativas[2], result)
	go bruteWorker(&wg, &done, 30, 40, numbers, numbersSum, &tentativas[3], result)
	go bruteWorker(&wg, &done, 40, 50, numbers, numbersSum, &tentativas[4], result)
	go bruteWorker(&wg, &done, 50, 61, numbers, numbersSum, &tentativas[5], result)
	wg.Wait()

	var totalTentativas uint64
	for _, v := range tentativas {
		totalTentativas += v
	}

	return totalTentativas, result
}
