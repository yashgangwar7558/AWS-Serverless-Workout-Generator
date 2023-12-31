package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	// "github.com/aws/aws-lambda-go/lambda"

	openai "github.com/sashabaranov/go-openai"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func handler() error {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	client := openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: `give me 45min workout plan using weightlifting exercises with barbel only. 
					Follow this format: Name of exercise (Minutes): Number of reps * Number of sets.`,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return err
	}

	fmt.Println(resp.Choices[0].Message.Content)

	clientTwilio := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_USERNAME"),
		Password: os.Getenv("TWILIO_PASSWORD"),
	})

	phoneNumbers := []string{"+918474907558", "+917451024845"}

	for _, phoneNumber := range phoneNumbers {
		fmt.Println(phoneNumber)
		params := &twilioApi.CreateMessageParams{}
		params.SetTo(fmt.Sprintf("whatsapp:%s", phoneNumber)) // fmt.Sprintf("whatsapp:%s", os.Getenv("TO_NUMBER"))
		params.SetFrom(fmt.Sprintf("whatsapp:%s", os.Getenv("FROM_NUMBER")))
		params.SetBody(resp.Choices[0].Message.Content)

		_, err = clientTwilio.Api.CreateMessage(params)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	log.Println("Workout plan has been sent")

	return nil
}

func main() {
	// lambda.Start(handler)
	handler()
}
