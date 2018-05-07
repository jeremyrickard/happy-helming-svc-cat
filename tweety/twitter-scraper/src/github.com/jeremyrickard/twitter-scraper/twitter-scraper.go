package main

import (
	"crypto/tls"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/mgo.v2"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	//mongoURI := strings.TrimSuffix(os.Getenv("MONGO_URI"), "\n")
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
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Errorf("error getting client: %s", err)
	}
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB("kubecon").C("tweets")

	if collection == nil {
		fmt.Errorf("couldn't get collection")
	}
	accessToken := strings.TrimSuffix(os.Getenv("ACCESS_TOKEN"), "\n")
	accessTokenSecret := strings.TrimSuffix(os.Getenv("ACCESS_TOKEN_SECRET"), "\n")
	consumerKey := strings.TrimSuffix(os.Getenv("CONSUMER_KEY"), "\n")
	consumerSecret := strings.TrimSuffix(os.Getenv("CONSUMER_SECRET"), "\n")

	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret)
	v := url.Values{}
	v.Set("track", "#Kubecon")
	api.EnableThrottling(10*time.Second, 60)
	stream := api.PublicStreamFilter(v)
	for {
		message := <-stream.C
		if message != nil {
			tweet, ok := message.(anaconda.Tweet)
			if ok {
				fmt.Printf(tweet.Text)
				err := collection.Insert(&tweet)
				if err != nil {
					fmt.Printf("error: %s", err)
				}
			} else {
				fmt.Println("got something else...")
			}

		} else {
			fmt.Println("got an empty message")
		}

		fmt.Println("looping")
	}

}
