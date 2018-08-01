package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	Userkey string `json:"user_key"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Playing struct {
	Userkey      string    `sql:"not null"`
	CreatedAt    time.Time `sql:"not null"`
	TryCount     int       `sql:"not null"`
	AnswerNumber string    `sql:"not null"`
}

type Record struct {
	Userkey   string    `sql:"not null"`
	CreatedAt time.Time `sql:"not null"`
	Nickname  string    `sql:"not null"`
	Score     int       `sql:"not null"`
}

type UserInfo struct {
	Userkey     string    `sql:"not null"`
	CreatedAt   time.Time `sql:"not null"`
	Nickname    string
	IsItUpdated bool `sql:"not null"`
}

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("postgres", "user=golang dbname=onetonine password=ilovegolang sslmode=disable")

	if err != nil {
		panic(err)
	}

	Db.AutoMigrate(&Playing{}, &Record{}, &UserInfo{})

}

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

func ScoreCalculater(StartTime time.Time, EndTime time.Time, TryCount int) int {
	SpendedTimeFloat := EndTime.Sub(StartTime).Seconds()

	SpendedTime := int(SpendedTimeFloat)

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

func sendMessage(w http.ResponseWriter, message ...string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	template := `{
	"message":{
		"text" : "%v"
	}
}`

	response := fmt.Sprintf(template, strings.Join(message, ""))
	w.Write([]byte(response))
	return
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	var playing Playing

	var userinfo UserInfo

	var post Post
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)

	playing.Userkey = post.Userkey

	if err != nil {
		w.WriteHeader(400)
		return
	}

	if post.Type != "text" {
		sendMessage(w, "텍스트 전송만 지원됩니다.")
		return
	}

	if err := Db.Where("userkey = $1", post.Userkey).First(&userinfo).Error; err != nil {
		userinfo.Userkey = post.Userkey
		userinfo.IsItUpdated = false
		Db.Create(&userinfo)

		sendMessage(w, "순위에 사용될 닉네임을 입력해주세요. (필수)")
		return
	}

	if err := Db.Where("userkey = $1", post.Userkey).First(&userinfo).Error; err == nil && userinfo.IsItUpdated == false {
		Db.Delete(UserInfo{}, "userkey LIKE ?", post.Userkey)

		userinfo.Nickname = post.Content
		userinfo.IsItUpdated = true

		Db.Create(&userinfo)
		sendMessage(w, userinfo.Nickname, "을(를) 닉네임으로 설정하였습니다. 수정하려면 '수정'이라고 입력해주세요.")
		return
	}

	switch post.Content {
	case "시작":
		if err := Db.Where("userkey = $1", playing.Userkey).First(&playing).Error; err != nil {
			playing.AnswerNumber = GetThreeRandomNumber()
			Db.Create(&playing)
			sendMessage(w, "게임이 시작되었습니다. 추리를 시작해주세요.")
			return

		} else {
			sendMessage(w, "이미 게임이 진행 중입니다.")
			return
		}
	case "순위":
		sendMessage(w, "순위 기능은 아직 미구현!")
		return
	case "수정":
		Db.Where("userkey = $1", post.Userkey).First(&userinfo)
		Db.Delete(UserInfo{}, "userkey LIKE ?", post.Userkey)

		userinfo.IsItUpdated = false

		Db.Create(&userinfo)

		sendMessage(w, "예전 닉네임: ", userinfo.Nickname, ", 수정할 닉네임을 입력해주세요.")
		return
	default:
		if _, err := strconv.Atoi(post.Content); err != nil {
			sendMessage(w, "'시작', '순위', '수정'을 입력하실 수 있습니다.")
			return
		}

		if isItValid := CheckAnswerValidation(post.Content); !isItValid {
			sendMessage(w, "모두 다른 세개의 숫자를 입력해주세요.")
			return
		}

		if err := Db.Where("userkey = $1", playing.Userkey).First(&playing).Error; err == nil {
			playing.TryCount++

			strike, ball := Checker(playing.AnswerNumber, post.Content)

			if strike == 3 {
				now := time.Now()

				score := ScoreCalculater(playing.CreatedAt, now, playing.TryCount-1)

				var record Record
				if err := Db.Where("userkey = $1", playing.Userkey).First(&record).Error; err != nil {
					record.Score = score
				} else {
					Db.Delete(Record{}, "userkey LIKE ?", post.Userkey)
					if record.Score < score {
						record.Score = score
					}
				}

				record.Userkey = post.Userkey

				Db.Where("userkey = $1", playing.Userkey).First(&userinfo)
				record.Nickname = userinfo.Nickname

				Db.Create(&record)
				Db.Delete(Playing{}, "userkey LIKE ?", post.Userkey)

				sendMessage(w, "정답입니다! 당신의 점수는 "+strconv.Itoa(score)+"입니다!")
				return

			} else if strike == 0 && ball == 0 {
				sendMessage(w, "아웃!")
				Db.Delete(Playing{}, "userkey LIKE ?", post.Userkey)
				Db.Create(&playing)
				return

			} else {
				sendMessage(w, strconv.Itoa(strike), " 스트라이크, ", strconv.Itoa(ball), " 볼!")
				Db.Delete(Playing{}, "userkey LIKE ?", post.Userkey)
				Db.Create(&playing)
				return
			}
		} else {
			sendMessage(w, "아직 게임이 시작하지 않았습니다. '시작'을 입력해주세요.")
			return
		}
		sendMessage(w, "잠깐... 예기치 못한 오류!")
		log.Println("오류 발생! userkey:", post.Userkey)
		return
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
