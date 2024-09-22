package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func getEnv(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	return os.Getenv(key)
}

func main() {

	slackBotToken := getEnv("SLACK_BOT_TOKEN")
	slackAppToken := getEnv("SLACK_APP_TOKEN")

	bot := slacker.NewClient(slackBotToken, slackAppToken)

	go printCommandEvents(bot.CommandEvents())

	bot.Command("My year of birth is <year>", &slacker.CommandDefinition{
		Description: "Year of birth calculator",
		Examples: []string{
			"My year of birth is 1999",
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)

			if err != nil {
				log.Fatalf("Error handeling command: %s", err)
			}

			age := 2024 - yob
			r := fmt.Sprintf("Your age is %d", age)
			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatalf("Error listening to bot: %s", err)
	}
}
