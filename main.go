package main

import (
  "bytes"
	"encoding/json"
  "fmt"
	"github.com/aws/aws-lambda-go/lambda"
  "github.com/joho/godotenv"
  "net/http"
  "os"
)

const IS_PROD_ENVIRONMENT = true

type LambdaEvent struct {
}

type LambdaResponse struct {
	Message string `json:"message"`
}

func GetBotToken(token_env_name string) (string) {
  if (IS_PROD_ENVIRONMENT) {
    return os.Getenv(token_env_name)
  } else {
    envFile, _ := godotenv.Read(".env")
    return envFile[token_env_name]
  }
}

func GetChatId(chatid_env_name string) (string) {
  if (IS_PROD_ENVIRONMENT) {
    return os.Getenv(chatid_env_name)
  } else {
    envFile, _ := godotenv.Read(".env")
    return envFile[chatid_env_name]
  }
}

func SendToTelegramBot(text string, bot string, chat_id string) {
	request_url := "https://api.telegram.org/" + bot + "/sendMessage"
	
	client := &http.Client{}
	values := map[string]string{"text": text, "chat_id": chat_id }
	json_paramaters, _ := json.Marshal(values)
	req, _:= http.NewRequest("POST", request_url, bytes.NewBuffer(json_paramaters))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	
	if(err != nil){
		fmt.Println(err)
	} else {
    fmt.Println(res.Status)
		defer res.Body.Close()
	}	
}

func HelloWorldMessage() {
	SendToTelegramBot("hello telegram", "bot" + GetBotToken("HOSHINO_AI_CRYPTO_BOT_TOKEN"), GetChatId("MY_CHAT_ID"))
}

func HandleLambdaEvent(event *LambdaEvent) (*LambdaResponse, error) {
  HelloWorldMessage()

	if event == nil {
		return &LambdaResponse{Message: "error"}, nil
	} else {
    return &LambdaResponse{Message: "done"}, nil
  }
}

func main() {
  if (IS_PROD_ENVIRONMENT) {
    lambda.Start(HandleLambdaEvent)
  } else {
    HelloWorldMessage()
  }
}