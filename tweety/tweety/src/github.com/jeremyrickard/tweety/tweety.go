package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func getTweetHandler(dialInfo *mgo.DialInfo) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}
		defer session.Close()
		collection := session.DB("kubecon").C("tweets")
		fmt.Println("handling request")
		getTweet(w, collection)
	}
}

func getTweet(w http.ResponseWriter, collection *mgo.Collection) {
	tweet := &anaconda.Tweet{}
	count, err := collection.Count()
	if count > 100 {
		count = 100
	}
	if err != nil {
		fmt.Printf("got an error %s", err)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	err = collection.Find(nil).Skip(rand.Intn(count)).Limit(10).One(tweet)
	if err != nil {
		fmt.Printf("got an error %s", err)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	fmt.Printf("found %s", tweet.Text)
	fmt.Fprintf(w, "%s", generateMessage(tweet.Text))
}

func generateMessage(message string) string {

	var lines bytes.Buffer
	words := strings.Split(message, " ")
	//word, words := words[0], words[1:]
	for len(words) > 0 {
		var word string
		var line string
		for len(words) > 0 && len(line)+len(words[0])+1 < 39 {
			word, words = words[0], words[1:]
			line = line + " " + word

		}
		outputLine := fmt.Sprintf("   | %s", line)
		filler := 43 - len(outputLine)
		for i := 0; i < filler; i++ {
			outputLine = outputLine + " "
		}
		outputLine = outputLine + "|"
		lines.WriteString(outputLine)
		lines.WriteString("\n")
		line = ""
	}

	header := `
   /---------------------------------------\
`

	footer := `   \---------------------------------------/
    \
     \
      \ __
       /  \
       |  |
       @  @
       |  |
       || |/
       || ||
       |\_/|
       \___/
	`

	return header + lines.String() + footer
}

func main() {
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPort := os.Getenv("MONGO_PORT")
	dialInfo := &mgo.DialInfo{
		Addrs: []string{
			fmt.Sprintf(
				"%s:%s",
				mongoHost,
				mongoPort,
			),
		},
		Timeout:  60 * time.Second,
		Database: "kubecon",
		Username: mongoUser,
		Password: mongoPassword,
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
	}
	http.HandleFunc("/", getTweetHandler(dialInfo))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
