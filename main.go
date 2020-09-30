package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/montanaflynn/stats"
)

func main() {
	//Loads DISCORD_TOKEN variable from .env
	godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	//Creates session
	dg, _ := discordgo.New("Bot " + token)
	//Sets up listeners
	dg.AddHandler(sendMessage)
	//Opens up a connection
	dg.Open()
	defer dg.Close()

	//Blocks from quitting unless any of the signals below are received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
func getStats(s string) string {

	byteValue := []byte(s)
	arr := []float64{}
	err := json.Unmarshal(byteValue, &arr)
	if err != nil {
		return "Please enter valid input"
	}
	sort.Float64s(arr)

	v, _ := stats.Variance(arr)
	sd, _ := stats.StandardDeviation(arr)

	mean, _ := stats.Mean(arr)
	median, _ := stats.Median(arr)
	mode, _ := stats.Mode(arr)
	retVal := fmt.Sprintf("Variance: %0.4f\n", v)
	retVal += fmt.Sprintf("SD: %0.4f\n", sd)
	retVal += fmt.Sprintf("Range: %0.4f\n", arr[len(arr)-1]-arr[0])

	retVal += fmt.Sprintf("Mean: %0.4f\nMedian: %0.4f\nMode: ", mean, median)

	for _, item := range mode {
		retVal += fmt.Sprintf("%0.4f, ", item)
	}
	retVal = strings.TrimSuffix(retVal, ", ")
	retVal += "\n"
	quartiles, _ := stats.Quartile(arr)
	retVal += fmt.Sprintf("Q1: %0.4f, Q2: %0.4f, Q3: %0.4f\n", quartiles.Q1, quartiles.Q2, quartiles.Q3)
	iqr, _ := stats.InterQuartileRange(arr)
	retVal += fmt.Sprintf("IQR: %0.4f\n", iqr)

	return retVal
}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "vstats") {
		// s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, getStats(m.Content[7:]))
	}
}
