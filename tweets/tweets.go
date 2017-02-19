package tweets

import (
	"github.com/ChimeraCoder/anaconda"
	log "github.com/Sirupsen/logrus"

	"github.com/hirondelle-app/listener/config"
	"github.com/hirondelle-app/listener/firebase"
)

// Keywords ...
var Keywords []string

// StartStream function
func StartStream(s *anaconda.Stream) {

	log.Info("Streaming tweets")

	for {
		item := <-s.C
		switch status := item.(type) {
		case anaconda.Tweet:
			if retweet := status.RetweetedStatus; retweet != nil {
				if retweet.RetweetCount >= config.Cfg.Retweets && retweet.FavoriteCount >= config.Cfg.Likes {
					if isTweetInserted := stringInSlice(retweet.IdStr, config.TweetsInserted); isTweetInserted == false {

						keyword := getKeywordFromTweet(retweet)

						if keyword == "" {
							log.Warningln("Empty keyword for tweet ", retweet.IdStr)
							break
						}

						t := map[string]interface{}{
							"tweetId":   retweet.IdStr,
							"likes":     retweet.FavoriteCount,
							"retweets":  retweet.RetweetCount,
							"createdAt": retweet.CreatedAt,
							"keyword":   keyword,
						}

						log.Infof("%+v\n", t)

						if err := firebase.InsertTweets(t); err != nil {
							log.Fatal(err)
						}

						config.TweetsInserted = append(config.TweetsInserted, retweet.IdStr)
					}
				}
			}
		}
	}

}
