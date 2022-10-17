package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
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
	// Make get request to each and use the responses
	var serverStatus = make([]ServerStatus, len(serverNames)) // Data that I want to save
	records := url.Values{}                                   // Data for loggly

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

		// Turn data into go struct
		json.Unmarshal([]byte(body), &serverStatus[index])

		// Print stuff for me to see
		fmt.Println("\nServer Status: ", response.Status)
		printServerStatus(serverStatus[index])

		// Values for loggly
		records.Add(serverAddress, response.Status)
		defer response.Body.Close()
	}

	// Send the record to loggly
	var logglyURL = os.Getenv("Loggly_Token")
	_, err := http.PostForm(logglyURL, records)
	if err != nil {
		panic(err)
	}
	fmt.Println("Values send to Loggly: ", records)

	// Send the record to AWS
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
}

func printServerStatus(server ServerStatus) {
	fmt.Println("Hostname: ", server.Hostname)
	fmt.Println("IP: ", server.IP)
	fmt.Println("Version: ", server.Version)
	fmt.Println("Online: ", server.Online)
	fmt.Println("Max players: ", server.Players.Max)
	fmt.Println("Online players: ", server.Players.Online)
}
