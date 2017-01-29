package main

import (
	"fmt"
	"log"
	"net/url"

	"encoding/json"

	"github.com/ChimeraCoder/anaconda"
	"github.com/jessevdk/go-flags"
)

var config struct {
	Likes    int    `short:"l" long:"likes" description:"Minimun likes for a tweet" required:"true"`
	Retweets int    `short:"r" long:"retweets" description:"Minimun retweets for a tweet" required:"true"`
	Keywords string `short:"k" long:"keywords" description:"Keywords for the streaming api" required:"true"`
	Twitter  struct {
		AccessToken       string `long:"access-token" description:"Access token for Twitter Api" required:"true"`
		AccessTokenSecret string `long:"access-token-secret" description:"Secret access token for Twitter Api" required:"true"`
		ConsumerKey       string `long:"consumer-key" description:"Consumer key for Twitter Api" required:"true"`
		ConsumerSecret    string `long:"consumer-secret" description:"Consumer secret for Twitter Api" required:"true"`
	}
}

type tweet struct {
	ID       string `json:"twitter_id"`
	Likes    int    `json:"likes"`
	Retweets int    `json:"retweets"`
}

func main() {
	_, err := flags.Parse(&config)

	if err != nil {
		log.Fatal(err)
	}

	anaconda.SetConsumerKey(config.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(config.Twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.Twitter.AccessToken, config.Twitter.AccessTokenSecret)

	v := url.Values{}
	v.Set("track", "python,golang,javascript")
	s := api.PublicStreamFilter(v)

	for {
		item := <-s.C
		switch status := item.(type) {
		case anaconda.Tweet:
			if retweet := status.RetweetedStatus; retweet != nil {
				if retweet.RetweetCount >= 20 && retweet.FavoriteCount >= 20 {
					tweet := &tweet{}
					tweet.ID = retweet.IdStr
					tweet.Likes = retweet.FavoriteCount
					tweet.Retweets = retweet.RetweetCount
					tweetJson, _ := json.Marshal(tweet)
					fmt.Println(string(tweetJson))
				}
			}
		}
	}
}
