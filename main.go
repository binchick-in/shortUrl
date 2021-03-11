package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
	"encoding/json"
)

const (
	portLowerLimit    = 0  // The minimum port value
	portUpperLimit    = 65353  // The maximum port value
	defaultPortNumber = "8080"  // The default
)

var jsonFileData map[string]string

// Get an open file object used for writing log lines.
func getLogFile() *os.File {
	f, err := os.OpenFile("webserver.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	return f
}

func getJsonFileData() map[string]string {
	jsonLookupFile, err := ioutil.ReadFile("urlMap.json")
	if err != nil {
		log.Fatal("Could not open url map file. Check that it exists...")
		panic(err)
	}
	var data map[string]string
	_ = json.Unmarshal([]byte(jsonLookupFile), &data)
	return data
}

// Check to ensure the inputted port number is a valid port number
// return false if not...
func isValidPortNumber(rawPort string) bool {
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		fmt.Printf("%s is not a valid port\n", rawPort)
		log.Fatalf("%s is not a valid port\n", rawPort)
		return false
	}
	if port >= portLowerLimit && port <= portUpperLimit {
		return true
	}
	return false
}

/*
Pass a slice of strings, return only the non-empty
strings in a new array
*/
func cleanSlice(i []string) []string {
	var retVal []string

	for _, item := range i {
		i := strings.TrimSpace(item)
		if len(i) > 0 {
			retVal = append(retVal, item)
		}
	}
	return retVal
}

// The one and only route that exists in this webapp.
func soleRoute(w http.ResponseWriter, req *http.Request) {
	cleanPath := cleanSlice(strings.Split(req.URL.Path, "/"))
	fmt.Println("Cleaned path => ", cleanPath, "Length of cleaned path => ", len(cleanPath))

	if len(cleanPath) != 1 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid Input\n")
		return
	}

	lookupKey := cleanPath[0]
	result, ok := jsonFileData[lookupKey]
	if !ok {
		log.Printf("%s was not found in the map\n", lookupKey)
		w.WriteHeader(404)
		fmt.Fprintf(w, "Key not found\n")
		return
	}

	w.Header().Set("Location", result)
	w.WriteHeader(302)
	fmt.Fprintf(w, "")
}

func main() {
	jsonFileData = getJsonFileData()
	logFile := getLogFile()
	log.SetOutput(logFile)

	// jsonLookupData := getJsonFileData()
	var port string
	if len(os.Args) >= 2 {
		rawPortInput := os.Args[1]
		if isValidPortNumber(rawPortInput) {
			port = fmt.Sprintf(":%s", rawPortInput)
		} else {
			os.Exit(1)
		}
	} else {
		fmt.Printf("Using default port: %s\n", defaultPortNumber)
		port = ":8080"
	}

	fmt.Printf("Server listening on port: %s\n", port)
	http.HandleFunc("/", soleRoute)
	http.ListenAndServe(port, nil)
}
