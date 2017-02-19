package cmd

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/cobra"
	"github.com/zabawaba99/firego"

	"github.com/hirondelle-app/listener/config"
	"github.com/hirondelle-app/listener/firebase"
	"github.com/hirondelle-app/listener/tweets"

	log "github.com/Sirupsen/logrus"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Listen tweets with specifics hashtags",
	Long:  `Listen tweets with specifics hashtags and insert them in Firebase`,
	RunE:  run,
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func run(cmd *cobra.Command, args []string) error {
	log.Info("Run cmd")

	// init Firebase
	firebase.InitFirebase()

	notifications := make(chan firego.Event)
	fbRealTime := firebase.FB.Child("keywords/byId")

	if err := fbRealTime.Watch(notifications); err != nil {
		log.Fatal(err)
	}

	defer fbRealTime.StopWatching()

	// init Anaconda
	anaconda.SetConsumerKey(config.Cfg.ConsumerKey)
	anaconda.SetConsumerSecret(config.Cfg.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.Cfg.AccessToken, config.Cfg.AccessTokenSecret)

	var err error

	for _ = range notifications {
		log.Info("Watching keywords")

		tweets.Keywords, err = firebase.GetKeywords()
		if err != nil {
			log.Fatal(err)
		}

		log.Info(tweets.Keywords)

		keywordsStr := tweets.KeywordsToStr()

		v := url.Values{}
		v.Set("track", keywordsStr)
		v.Set("languages", "en")
		s := api.PublicStreamFilter(v)

		defer s.Stop()

		go tweets.StartStream(s)
	}

	return nil
}
