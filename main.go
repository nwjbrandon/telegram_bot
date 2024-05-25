package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"
)

const IS_PROD_ENVIRONMENT = true

func GetBotToken() string {
	env_name := "HOSHINO_AI_CRYPTO_BOT_TOKEN"
	if IS_PROD_ENVIRONMENT {
		return os.Getenv(env_name)
	} else {
		envFile, _ := godotenv.Read(".env")
		return envFile[env_name]
	}
}

func GetChatId() string {
	env_name := "MY_CHAT_ID"
	if IS_PROD_ENVIRONMENT {
		return os.Getenv(env_name)
	} else {
		envFile, _ := godotenv.Read(".env")
		return envFile[env_name]
	}
}

func GetCoinMarketCapToken() string {
	env_name := "X_CMC_PRO_API_KEY"
	if IS_PROD_ENVIRONMENT {
		return os.Getenv(env_name)
	} else {
		envFile, _ := godotenv.Read(".env")
		return envFile[env_name]
	}
}

func SendToTelegramBot(text string, bot string, chat_id string) {
	request_url := "https://api.telegram.org/" + bot + "/sendMessage"

	client := &http.Client{}
	values := map[string]string{"text": text, "chat_id": chat_id}
	json_paramaters, _ := json.Marshal(values)
	req, _ := http.NewRequest("POST", request_url, bytes.NewBuffer(json_paramaters))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.Status)
		defer res.Body.Close()
	}
}

func HelloWorldMessage() {
	SendToTelegramBot("hello telegram", "bot"+GetBotToken(), GetChatId())
}

type CryptoCurrencyQuoteUSD struct {
	Price float64 `json:price`
}

type CryptoCurrencyQuote struct {
	USD CryptoCurrencyQuoteUSD `json:USD`
}

type CryptoCurrencyLatestListingResponse struct {
	Data []struct {
		Name  string              `json:"name"`
		Quote CryptoCurrencyQuote `json:"quote"`
	} `json:"data"`
}

func CryptoCurrencyPriceUpdates() {
	// Create GET Request
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	req.Header.Add("X-CMC_PRO_API_KEY", GetCoinMarketCapToken())

	// Invoke GET Request
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Status Code: %d\n", res.StatusCode)
		return
	}

	// Parse Response
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var cryptoCurrencyLatestListingResponse CryptoCurrencyLatestListingResponse
	if err := json.Unmarshal(body, &cryptoCurrencyLatestListingResponse); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return
	}

	// Sort By Name
	sort.Slice(cryptoCurrencyLatestListingResponse.Data, func(i, j int) bool {
		return cryptoCurrencyLatestListingResponse.Data[i].Name < cryptoCurrencyLatestListingResponse.Data[j].Name
	})

	// Prepare Message
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Singapore")
	telegramMessage := fmt.Sprintf("Crypto Currency Prices As Of: %s (SGT)\n\n", now.In(loc))
	for i := 1; i < len(cryptoCurrencyLatestListingResponse.Data); i++ {
		name := cryptoCurrencyLatestListingResponse.Data[i].Name
		price := cryptoCurrencyLatestListingResponse.Data[i].Quote.USD.Price
		telegramMessage += (name + " " + fmt.Sprintf("%f", price) + "\n")
	}

	// Send To Telegram
	fmt.Println(telegramMessage)
	SendToTelegramBot(telegramMessage, "bot"+GetBotToken(), GetChatId())
}

type LambdaEvent struct {
}

type LambdaResponse struct {
	Message string `json:"message"`
}

func HandleLambdaEvent(event *LambdaEvent) (*LambdaResponse, error) {
  CryptoCurrencyPriceUpdates()
	if event == nil {
		return &LambdaResponse{Message: "error"}, nil
	} else {
		return &LambdaResponse{Message: "done"}, nil
	}
}

func main() {
	if IS_PROD_ENVIRONMENT {
		lambda.Start(HandleLambdaEvent)
	} else {
		// HelloWorldMessage()
		// CryptoCurrencyPriceUpdates()
		return
	}
}
