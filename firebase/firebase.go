package firebase

import (
	"fmt"

	"github.com/hirondelle-app/listener/config"
	"github.com/zabawaba99/firego"
)

const firebaseURL = "https://hirondelle-e44d5.firebaseio.com/"

// FB *Firego var
var FB *firego.Firebase

// InitFirebase Init Firebase
func InitFirebase() {
	FB = firego.New(firebaseURL, nil)
	FB.Auth(config.Cfg.FBToken)
}

// GetKeywords Get keywords from Firebase
func GetKeywords() ([]string, error) {
	var keywords []string

	fbKeywords := FB.Child("keywords/byId")

	var results map[string]interface{}
	if err := fbKeywords.Value(&results); err != nil {
		return keywords, err
	}

	for k := range results {
		keywords = append(keywords, k)
	}

	return keywords, nil
}

func InsertTweets(t map[string]interface{}) error {

	fbTweets := FB.Child(fmt.Sprintf("tweets/byId/%s", t["tweetId"]))
	if err := fbTweets.Set(t); err != nil {
		return err
	}

	tweetKeyword := map[string]string{t["tweetId"].(string): "true"}

	fbTweetKeywords := FB.Child(fmt.Sprintf("tweets/byKeyword/%s", t["keyword"]))
	if err := fbTweetKeywords.Update(tweetKeyword); err != nil {
		return err
	}
	return nil
}
