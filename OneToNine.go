package main

import (
	"fmt"
	"math"
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

func ScoreCalculater(StartTime time.Time, EndTime time.Time, TryCount int) int {
	SpendedTimeFloat := EndTime.Sub(StartTime).Seconds()

	SpendedTime := int(math.Round(SpendedTimeFloat))

	//(180-경과 시간(초)-횟수*5)*100
	fmt.Println("SpendedTime:", SpendedTime)
	fmt.Println("TryCount:", TryCount)
	score := (180 - SpendedTime - TryCount*5) * 100

	return score
}

func GetThreeRandomNumber() (AnswerNumber string) {
	n := 0
	rand.Seed(time.Now().UnixNano())
	r := rand.Perm(9)

	for n < 3 {
		AnswerNumber += strconv.Itoa(1 + r[n])
		n += 1
	}

	return

}

func play() {

	answer := GetThreeRandomNumber()
	fmt.Println("무작위 세 숫자를 뽑았습니다.")

	fmt.Println("추리를 시작하세요.")

	TryCount := 0
	var Challenger string

	StartTime := time.Now()

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
			EndTime := time.Now()
			score := ScoreCalculater(StartTime, EndTime, TryCount+1)

			fmt.Print("정답! 내 점수는 ", score, "입니다.\n")

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
