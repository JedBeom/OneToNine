package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("postgres", "user=golang dbname=onetonine password=ilovegolang sslmode=disable")

	if err != nil {
		panic(err)
	}

	Db.AutoMigrate(&Playing{}, &Record{}, &UserInfo{})

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
			score, _ := ScoreCalculater(StartTime, EndTime, TryCount+1)

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

func keyboardHandler(w http.ResponseWriter, r *http.Request) {
	template := `{
	"type": "text"
}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write([]byte(template))
}

func main() {
	server := http.Server{
		Addr: ":80",
	}

	http.HandleFunc("/keyboard", keyboardHandler)
	http.HandleFunc("/message", messageHandler)

	server.ListenAndServe()
}
