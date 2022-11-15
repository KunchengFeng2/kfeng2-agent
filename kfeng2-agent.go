package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var loggly_Token string
var database *dynamodb.DynamoDB
var tableName string
var apiURL string
var serverNames []string

func init() {
	// loggly_Token and AWS login credentials should be set as environment variables
	loggly_Token = os.Getenv("Loggly_Token")
	if len(loggly_Token) == 0 {
		fmt.Println("Problem occured when reading loggly token from environment.")
	}

	// Establish a new connection
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Problem occured when forming new session!")
		os.Exit(1)
	}
	database = dynamodb.New(sess)
	tableName = "Kfeng2_MC_Servers"

	// API information + the servers that I want to check up.
	apiURL = "https://api.mcsrvstat.us/2/"
	serverNames = []string{
		"www.grmpixelmon.com",
		"mcsl.applecraft.org",
		"play.ml-mc.com",
		"mcsl.dogecraft.net",
		"mc.advancius.net",
		"mcsl.wildprison.net",
		"play.civmc.net",
		"mcsl.cosmosmc.org",
		"msl.mc-complex.com",
		"play.pixelmonrealms.com",
		"play.anubismc.com",
		"mcslf.pika.host",
		"play.pokesaga.org",
		"play.applemc.fun:19132",
		"Herobrine.org",
		"mcslf.jartex.fun",
		"play.fruitysmp.com",
		"mcsl.zedarmc.com",
		"mcsl.oneblockmc.com",
		"play.jackpotmc.com",
		"mcsl.lemoncloud.net",
		"play.vulengate.com",
		"sl.minecadia.com",
		"msl.mc-blaze.com",
	}
}

// The other functions are below main()
func main() {
	for {
		fmt.Println("Starting a new pulling cycle...")

		var dataSize = 0

		for _, server := range serverNames {
			var address = apiURL + server
			var body = apiCall(address)
			var bodyGo = toGoStruct(body)
			sendToDB(bodyGo)
			dataSize += len(body)
		}
		sendToLoggly(dataSize)

		fmt.Println("Waiting for next cycle... ")
		time.Sleep(30 * time.Minute)
	}
}

func apiCall(address string) []byte {
	response, err := http.Get(address)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	response.Body.Close()
	return body
}

func toGoStruct(data []byte) ServerStatus {
	var serverStatus = ServerStatus{}
	json.Unmarshal(data, &serverStatus)
	serverStatus.Time = time.Now().String()
	return serverStatus
}

// I'm only putting items into DB one at a time
// Although this will generate many API calls, AWS only charges base on the amount of
// data been send in.
// Plus batchWriteItem only allow at most 25 items at a time, so it complicates things.
func sendToDB(item ServerStatus) {
	awsItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		panic("Error marshalling into aws attribute map")
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      awsItem,
	}

	resp, err := database.PutItem(params)
	if err != nil {
		fmt.Println("Problem with putting item into DB.")
		fmt.Println("DynamoDB response: ", resp)
		panic(err)
	}
}

func sendToLoggly(size int) {
	var data = url.Values{
		"dataSize": {strconv.Itoa(size)},
	}

	response, err := http.PostForm(loggly_Token, data)
	if err != nil {
		fmt.Println("Error when sending data to loggly")
		fmt.Println("Loggly response: ", response)
		fmt.Println(err)
	}
}

// The actual response contains tons more information, I'm only using what I need.
// Golang assign values by same variable names, no need to mimic the complete structure.
type ServerStatus struct {
	IP       string
	Version  string
	Online   bool
	Hostname string
	Players  struct {
		Online int
		Max    int
	}
	Time string
}

func printServerStatus(server ServerStatus) {
	fmt.Println("Hostname: ", server.Hostname)
	fmt.Println("IP: ", server.IP)
	fmt.Println("Version: ", server.Version)
	fmt.Println("Online: ", server.Online)
	fmt.Println("Max players: ", server.Players.Max)
	fmt.Println("Online players: ", server.Players.Online)
	fmt.Println("Time: ", server.Time)
}
