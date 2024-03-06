package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SignUpStruct struct {
	Name          string
	TelegramLogin string
	Password      string
}

var SignUpSlice = []SignUpStruct{} // ? empty
func main() {
	r := gin.Default()

	r.Use(Cors)
	r.POST("/signup", SignUp)
	go Recovery()

	r.Run(":3434")
}

func Recovery() {
	// fmt.Println("Saaassas")
	ReadUser()
	res, err := tgbotapi.NewBotAPI("6847516848:AAFsEZJyI2LD5sJt-lkIxdUtp_KQlN5_KNY")

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	updates := tgbotapi.NewUpdate(0)                      //
	allUpdates, updatesErr := res.GetUpdatesChan(updates) //

	if updatesErr != nil {
		fmt.Printf("updatesErr: %v\n", updatesErr)
	}

	isExisting := false
	for eachUpdate := range allUpdates {
		if eachUpdate.Message.IsCommand() {
			if eachUpdate.Message.Command() == "reset" {
				for _, item := range SignUpSlice {
					if item.TelegramLogin == eachUpdate.Message.Chat.UserName {
						isExisting = true
					}
				}

				if isExisting {
					msg := tgbotapi.NewMessage(eachUpdate.Message.Chat.ID, "Enter your pass")
					res.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(eachUpdate.Message.Chat.ID, "Can't find username")
					res.Send(msg)
				}
			}
		} else {
			if isExisting {
				fmt.Println(eachUpdate.Message.Text)
				for index, item := range SignUpSlice {
					if item.TelegramLogin == eachUpdate.Message.Chat.UserName {
						SignUpSlice[index].Password = eachUpdate.Message.Text
						isExisting = false
						WriteUser()
					}
				}
			}
		}
	}
	
}

func SignUp(c *gin.Context) {
	var SignUpTemp SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Password == "" || SignUpTemp.TelegramLogin == "" {
		c.JSON(404, "Empty field")
	} else {
		ReadUser()
		SignUpSlice = append(SignUpSlice, SignUpTemp)
		WriteUser()
	}
}

func WriteUser() {
	marshalledData, _ := json.Marshal(SignUpSlice)
	ioutil.WriteFile("db.json", marshalledData, 0644)
}

func ReadUser() {
	readByte, _ := ioutil.ReadFile("db.json")
	json.Unmarshal(readByte, &SignUpSlice)
}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.246:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}
