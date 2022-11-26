package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env")
	}

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(ActivityStatusUpdate)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildPresences

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening session", err)
		return
	}

	fmt.Println("Bot is running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = dg.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ActivityStatusUpdate(s *discordgo.Session, v *discordgo.PresenceUpdate) {
	now := time.Now()
	userId := v.User.ID

	if len(v.Activities) == 0 {
		// TODO: ゲーム終了したときの処理
	} else {
		// TODO: ゲーム開始したときの処理
		gameName := v.Activities[0].Name
	}
}
