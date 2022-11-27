package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
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
	apiURL := os.Getenv("API_URL")
	userId := v.User.ID
	client := &http.Client{}

	// TODO それぞれ認証情報追加する, Query->他の形式に変えるかも
	if len(v.Activities) == 0 && v.Status == discordgo.StatusOnline {
		params := "userId=" + url.QueryEscape(userId)
		endpoint := fmt.Sprintf("%s/endplaying?%s", apiURL, params)
		PostURL(endpoint, client)
	} else {
		gameName := v.Activities[0].Name
		params := "userId=" + url.QueryEscape(userId) + "&gameName=" + url.QueryEscape(gameName)
		endpoint := fmt.Sprintf("%s/startplaying?%s", apiURL, params)
		PostURL(endpoint, client)
	}
}

func PostURL(endpoint string, client *http.Client) {
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}
