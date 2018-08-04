package main

import (
	"math/rand"
	"strconv"
	"time"
)

func CheckAnswerValidation(Challenger string) bool {
	if len(Challenger) != 3 {
		return false
	}

	for i := range Challenger {
		if Challenger[i] == byte('0') {
			return false
		}
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

func ScoreCalculater(StartTime time.Time, EndTime time.Time, TryCount int) (int, int) {
	SpendedTimeFloat := EndTime.Sub(StartTime).Seconds()

	SpendedTime := int(SpendedTimeFloat)

	//(180-경과 시간(초)-횟수*5)*100
	//fmt.Println("SpendedTime:", SpendedTime)
	//fmt.Println("TryCount:", TryCount)
	score := (180 - SpendedTime - TryCount*5) * 100

	return score, SpendedTime
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
