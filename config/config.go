package config

import "github.com/zabawaba99/firego"

var (
	Cfg            Config
	keywords       []string
	TweetsInserted []string
	firebase       *firego.Firebase
)

type Config struct {
	AccessToken       string `mapstructure:"access-token"`
	AccessTokenSecret string `mapstructure:"access-token-secret"`
	ConsumerKey       string `mapstructure:"consumer-key"`
	ConsumerSecret    string `mapstructure:"consumer-secret"`
	FBToken           string `mapstructure:"firebase-token"`
	Likes             int    `mapstructure:"likes"`
	Retweets          int    `mapstructure:"retweets"`
}
