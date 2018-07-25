package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func CheckAnswerValidation(Challenger string) bool {
	if len(Challenger) != 3 {
		return false
	}

	for x := 0; x < 2; x++ {
		for y := 1; y < 3; y++ {
			if x == y {
				continue
			}

			if Challenger[x] == Challenger[y] {
				return false
			}
		}
	}

	return true
}

func Checker(Original string, Challenger string) (StrikeCount int, BallCount int) {

	// Strike Check
	for x := 0; x < 3; x++ {
		if Original[x] == Challenger[x] {
			StrikeCount++
		}
	}

	// Ball Check
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if x == y {
				continue
			}

			if Original[x] == Challenger[y] {
				BallCount++
			}
		}
	}
	return
}

func RandomOneToNine() string {
	// 현 시간을 시드로 받아옴.
	rand.Seed(time.Now().UnixNano())
	// 1-9까지의 숫자를 무작위 생성.
	IntNumber := rand.Intn(9) + 1
	// int를 str로 변환.
	StrNumber := strconv.Itoa(IntNumber)

	return StrNumber
}

func GetThreeRandomNumber() (AnswerNumber string) {
	var AnswerNumberSlice []string

	// 중복 상관 없이 무작위 숫자 3개를 뽑는다.
	for len(AnswerNumberSlice) < 3 {
		number := RandomOneToNine()
		AnswerNumberSlice = append(AnswerNumberSlice, number)
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
			for AnswerNumberSlice[x] == AnswerNumberSlice[y] {
				number := RandomOneToNine()
				AnswerNumberSlice[y] = number
			}
		}
	}
	AnswerNumber = AnswerNumberSlice[0] + AnswerNumberSlice[1] + AnswerNumberSlice[2]

	return

}

func play() {

	answer := GetThreeRandomNumber()
	fmt.Println("무작위 세 숫자를 뽑았습니다.")

	fmt.Println("추리를 시작하세요.")

	TryCount := 0
	var Challenger string

	for {
		fmt.Print("> ")
		fmt.Scanf("%v", &Challenger)

		isItValid := CheckAnswerValidation(Challenger)

		if !isItValid {
			fmt.Println("각기 다른 세자리 숫자를 입력해주세요!")
			continue
		}

		strike, ball := Checker(answer, Challenger)

		if strike == 3 {
			fmt.Print("정답! ", TryCount+1, "회의 시도 끝에 성공하였습니다.\n")
			break
		} else if strike == 0 && ball == 0 {
			fmt.Println("아웃!")
		} else {
			fmt.Println(strike, "스트라이크,", ball, "볼!")
		}

		TryCount++
	}
}

func main() {
	play()
}
