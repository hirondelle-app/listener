package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"encoding/gob"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bitly/go-nsq"
	"github.com/jessevdk/go-flags"
)

var keywords []keyword

var config struct {
	Likes        int    `short:"l" long:"likes" description:"Minimun likes for a tweet" required:"true"`
	Retweets     int    `short:"r" long:"retweets" description:"Minimun retweets for a tweet" required:"true"`
	QueueAddress string `short:"q" long:"queue-address" description:"Queue address" required:"true"`
	ApiAddress   string `short:"a" long:"api-address" description:"Api address" required:"true"`
	Twitter      struct {
		AccessToken       string `long:"access-token" description:"Access token for Twitter Api" required:"true"`
		AccessTokenSecret string `long:"access-token-secret" description:"Secret access token for Twitter Api" required:"true"`
		ConsumerKey       string `long:"consumer-key" description:"Consumer key for Twitter Api" required:"true"`
		ConsumerSecret    string `long:"consumer-secret" description:"Consumer secret for Twitter Api" required:"true"`
	}
}

type tweet struct {
	TweetID   string `json:"twitter_id"`
	Likes     int    `json:"likes"`
	Retweets  int    `json:"retweets"`
	KeywordID int64  `json:"keywordID"`
}

func main() {
	_, err := flags.Parse(&config)

	if err != nil {
		log.Fatal(err)
	}

	anaconda.SetConsumerKey(config.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(config.Twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.Twitter.AccessToken, config.Twitter.AccessTokenSecret)

	keywords = getKeywords()

	v := url.Values{}
	v.Set("track", keywordsToStr())
	v.Set("languages", "en")
	s := api.PublicStreamFilter(v)

	tweetIDList := make([]string, 0)

	configNSQ := nsq.NewConfig()
	producer, err := nsq.NewProducer(config.QueueAddress, configNSQ)
	if err != nil {
		log.Fatal(err)
	}

	for {
		item := <-s.C
		switch status := item.(type) {
		case anaconda.Tweet:
			if retweet := status.RetweetedStatus; retweet != nil {
				if retweet.RetweetCount >= config.Retweets && retweet.FavoriteCount >= config.Likes {
					if tweetInDB := stringInSlice(retweet.IdStr, tweetIDList); tweetInDB == false {

						keywordID := getKeywordIDFromTweet(retweet)

						if keywordID == 0 {
							break
						}

						tweet := &tweet{
							TweetID:   retweet.IdStr,
							Likes:     retweet.FavoriteCount,
							Retweets:  retweet.RetweetCount,
							KeywordID: keywordID,
						}

						buf := new(bytes.Buffer)
						enc := gob.NewEncoder(buf)
						enc.Encode(tweet)

						err = producer.PublishAsync("tweets", buf.Bytes(), nil)
						if err != nil {
							log.Panic("Could not connect")
						}
						tweetIDList = append(tweetIDList, retweet.IdStr)
					}
				}
			}
		}
	}
}

func getKeywordIDFromTweet(tweet *anaconda.Tweet) int64 {

	var keywordID int64

	// Get hashtag
	for _, k := range keywords {
		for _, h := range tweet.Entities.Hashtags {
			hashtagLower := strings.ToLower(h.Text)
			if k.Label == hashtagLower {
				keywordID = k.ID
				return keywordID
			}
		}
	}

	// We still don't have a keyword, so we take user @mention
	for _, k := range keywords {
		for _, u := range tweet.Entities.User_mentions {
			userLower := strings.ToLower(u.Screen_name)
			if k.Label == userLower {
				keywordID = k.ID
				return keywordID
			}
		}
	}

	return keywordID
}

type keyword struct {
	ID    int64
	Label string `json:"label"`
}

func getKeywords() []keyword {
	var keywords []keyword

	resp, err := http.Get(fmt.Sprintf("%s/keywords", config.ApiAddress))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&keywords)

	return keywords
}

func keywordsToStr() string {
	var customKeywords string

	for _, k := range keywords {
		customKeywords += fmt.Sprintf("@%[1]s,#%[1]s,", k.Label)
	}

	return customKeywords

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
