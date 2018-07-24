package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func RandomOneToNine() string {
	// 현 시간을 시드로 받아옴.
	rand.Seed(time.Now().UnixNano())
	// 1-9까지의 숫자를 무작위 생성.
	IntNumber := rand.Intn(9) + 1
	// int를 str로 변환.
	StrNumber := strconv.Itoa(IntNumber)

	return StrNumber
}

func GetThreeRandomNumber() {
	var AnswerNumber []string

	// 중복 상관 없이 무작위 숫자 3개를 뽑는다.
	for len(AnswerNumber) < 3 {
		number := RandomOneToNine()
		AnswerNumber = append(AnswerNumber, number)
	}

	// 숫자들의 중복 확인을 한 후 중복된 숫자를 바꾼다.

	// x의 인덱스가 2까지 갈 필요는 없으니...
	for x := 0; x < 2; x++ {
		// y는 1부터 2까지 있으면 된다.
		for y := 1; y < 3; y++ {
			// 같은 인덱스끼리 비교하면 무한 루프니까 continue
			if x == y {
				continue
			}

			// 만약 [x]와 [y]가 같을 경우 y의 숫자를 다시 뽑는다.
			// 그래도 같으면 반복.
			for AnswerNumber[x] == AnswerNumber[y] {
				number := RandomOneToNine()
				AnswerNumber[y] = number
			}
		}
	}

	fmt.Println(AnswerNumber)
}

func main() {

	// 테스트용 무한 루프.
	for {
		GetThreeRandomNumber()
		// Seed가 달라야하니 1초 쉬기
		time.Sleep(time.Second)
	}
}
