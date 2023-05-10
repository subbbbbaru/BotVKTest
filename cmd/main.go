package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/subbbbbaru/botvktest/api"
	"github.com/subbbbbaru/botvktest/longpoll"
)

const groupId = "220402638"
const apiVersion = "5.131"

func main() {
	log.Println("Hello! I am BotVKTest!")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading ENV: %s", err.Error())
	}
	TOKEN := os.Getenv("TOKEN")

	vk := api.NewVK(TOKEN, apiVersion)

	lp, err := longpoll.NewLongpoll(vk, groupId)
	if err != nil {
		log.Fatal(err)
	}

	lp.LongpollHandler() // run bot
}
