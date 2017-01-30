package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

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
	TweetID  string `json:"twitter_id"`
	Likes    int    `json:"likes"`
	Retweets int    `json:"retweets"`
	Keyword  string `json:"keyword"`
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
	v.Set("track", customKeywords(config.Keywords))
	v.Set("languages", "en")
	s := api.PublicStreamFilter(v)

	tweetIDList := make([]string, 0)

	for {
		item := <-s.C
		switch status := item.(type) {
		case anaconda.Tweet:
			if retweet := status.RetweetedStatus; retweet != nil {
				if retweet.RetweetCount >= config.Retweets && retweet.FavoriteCount >= config.Likes {
					if tweetInDB := stringInSlice(retweet.IdStr, tweetIDList); tweetInDB == false {

						tweet := &tweet{
							TweetID:  retweet.IdStr,
							Likes:    retweet.FavoriteCount,
							Retweets: retweet.RetweetCount,
							Keyword:  getHashtag(config.Keywords, retweet),
						}
						tweetJson, _ := json.Marshal(tweet)
						fmt.Println(string(tweetJson))
						tweetIDList = append(tweetIDList, retweet.IdStr)
					}
				}
			}
		}
	}
}

func getHashtag(keywordsStr string, tweet *anaconda.Tweet) string {

	var keyword string

	// Get hashtag
	for _, k := range strings.Split(keywordsStr, ",") {
		for _, h := range tweet.Entities.Hashtags {
			hashtagLower := strings.ToLower(h.Text)
			if k == hashtagLower {
				keyword = hashtagLower
				break
			}
		}
	}

	// We still don't have a keyword, so we take user @mention
	for _, k := range strings.Split(keywordsStr, ",") {
		for _, u := range tweet.Entities.User_mentions {
			userLower := strings.ToLower(u.Screen_name)
			if k == userLower {
				keyword = userLower
				break
			}
		}
	}

	return keyword
}

func customKeywords(keywordsStr string) string {
	var customKeywords string

	for _, keyword := range strings.Split(keywordsStr, ",") {
		//customKeywords += fmt.Sprintf("@%[1]s,#%[1]s,", keyword)
		customKeywords += fmt.Sprintf("@%[1]s,", keyword)
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
