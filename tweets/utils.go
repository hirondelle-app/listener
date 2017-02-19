package tweets

import (
	"fmt"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

func getKeywordFromTweet(tweet *anaconda.Tweet) string {

	var keyword string

	// Get hashtag
	for _, k := range Keywords {
		for _, h := range tweet.Entities.Hashtags {
			hashtagLower := strings.ToLower(h.Text)
			if k == hashtagLower {
				keyword = k
				return keyword
			}
		}
	}

	// We still don't have a keyword, so we take user @mention
	for _, k := range Keywords {
		for _, u := range tweet.Entities.User_mentions {
			userLower := strings.ToLower(u.Screen_name)
			if k == userLower {
				keyword = k
				return keyword
			}
		}
	}

	return keyword
}

// KeywordsToStr ...
func KeywordsToStr() string {
	var customKeywords string

	for _, k := range Keywords {
		customKeywords += fmt.Sprintf("@%[1]s,#%[1]s,", k)
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
