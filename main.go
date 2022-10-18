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

// The other functions are below main()
func main() {

	// Websites that I'm using for this project
	var apiBaseJava = "https://api.mcsrvstat.us/2/"
	var serverNames = []string{
		"grmpixelmon.com",
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
		"play.vulengate.com",
		"sl.minecadia.com",
		"msl.mc-blaze.com",
	}

	// The main loop, placed in a different function for better control over intervals
	for {
		doStuff(apiBaseJava, serverNames)
		time.Sleep(30 * time.Minute)
	}
}

func doStuff(apiBaseJava string, serverNames []string) {
	// Create a AWS session
	// Credentials are load from environment variables, remember to set that on local machine.
	// Except AWS_REGION have to be set here, it dosen't work from environment for some reason
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println("Problem occured when forming new session!")
	}

	db := dynamodb.New(sess)

	// Make get request to each and use the responses
	var serverStatus = make([]ServerStatus, len(serverNames)) // Data that I want to save
	var dataSize = 0                                          // Datasize for loggly

	for index, serverAddress := range serverNames {
		var address = apiBaseJava + serverAddress

		// Making the get request
		response, err := http.Get(address)
		if err != nil {
			panic(err)
		}

		// Read the body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		dataSize += len(body)

		// Turn data into go struct
		json.Unmarshal([]byte(body), &serverStatus[index])
		serverStatus[index].Time = time.Now().String()

		// Print stuff for me to see
		fmt.Println("\nServer Status: ", response.Status)
		printServerStatus(serverStatus[index])

		// Send to AWS DynamoDB
		serverAWSMap, err := dynamodbattribute.MarshalMap(serverStatus[index])
		if err != nil {
			panic("Cannot marshal server status into AttributeValue map")
		}

		params := &dynamodb.PutItemInput{
			TableName: aws.String("Kfeng2_MC_Servers"),
			Item:      serverAWSMap,
		}
		resp, err := db.PutItem(params)
		if err != nil {
			fmt.Println("Problem with putting item into DB.")
			fmt.Println(err)
		}
		fmt.Println("DynamoDB response:", resp)

		defer response.Body.Close()
	}

	// Send the record to loggly
	var logglyURL = os.Getenv("Loggly_Token")
	var data = url.Values{
		"dataSize": {strconv.Itoa(dataSize)},
	}
	_, err = http.PostForm(logglyURL, data)
	if err != nil {
		panic(err)
	}
	fmt.Println("Values send to Loggly: ", data)

	fmt.Println("Waiting for next cycle... ")
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
