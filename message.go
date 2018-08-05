package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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
		if err := Db.First(&playing).Error; err != nil {
			playing.AnswerNumber = GetThreeRandomNumber()
			Db.Create(&playing)
			sendMessage(w, "게임이 시작되었습니다. 추리를 시작해주세요.")
			return

		} else {
			sendMessage(w, "이미 게임이 진행 중입니다.")
			return
		}
	case "순위":

		var record Record
		if err := Db.Where("userkey = $1", post.Userkey).First(&record).Error; err != nil {
			sendMessage(w, "일단 게임을 하고 오시는 게 어떨까요?")
			return
		}

		//rows, err := Db.Raw("select userkey, nickname, score, row_number () over (order by score desc) from records").Where("userkey = $1", post.Userkey).Rows()
		rows, err := Db.Table("records").Select("userkey, nickname, score, row_number () over (order by score desc)").Rows()
		if err != nil {
			log.Println(err)
			sendMessage(w, "일단 게임은 하고 오는 게 어때?")
			return
		}

		var rankers []RecordForShow
		var ranker RecordForShow

		for rows.Next() {
			rows.Scan(&ranker.Userkey, &ranker.Nickname, &ranker.Score, &ranker.RowNumber)
			rankers = append(rankers, ranker)
		}

		var myRecord RecordForShow
		var rankersByTemplate string

		recordTemplate := "%v등: %v | %v점"
		splitLine := "-------------"

		for i, recordforshow := range rankers {
			if recordforshow.Userkey == post.Userkey {
				if i > 3 {
					myRecord = recordforshow
					break
				}

			}

			if i <= 3 {

				rankerByTemplate := fmt.Sprintf(recordTemplate, recordforshow.RowNumber, recordforshow.Nickname, recordforshow.Score)
				rankersByTemplate += rankerByTemplate + "\\n"
			}

		}

		var content string

		if myRecord.Userkey != "" {
			myRecordString := fmt.Sprintf(recordTemplate, myRecord.RowNumber, myRecord.Nickname, myRecord.Score)

			content = rankersByTemplate + splitLine + "\\n" + myRecordString
		} else {
			content = rankersByTemplate

		}

		sendMessage(w, content)
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

		if err := Db.First(&playing).Error; err == nil {
			playing.TryCount++

			strike, ball := Checker(playing.AnswerNumber, post.Content)

			if strike == 3 {
				now := time.Now()

				score, SpendedTime := ScoreCalculater(playing.CreatedAt, now, playing.TryCount-1)

				var record Record
				if err := Db.Where("userkey = $1", playing.Userkey).First(&record).Error; err != nil {
					record.Score = score
					record.SpendedTime = SpendedTime
					record.TryCount = playing.TryCount
				} else {
					Db.Delete(Record{}, "userkey LIKE ?", post.Userkey)
					if record.Score < score {
						record.Score = score
						record.SpendedTime = SpendedTime
						record.TryCount = playing.TryCount
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
