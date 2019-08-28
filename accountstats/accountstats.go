package accountstats

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"../embed"

	"github.com/bwmarrin/discordgo"
)

var (
	owAPIprefix = "https://ow-api.com/v1/stats/pc/us/"
	owAPIsuffix = "/profile"
)

const (
	minBtagLength = 7
)

func isValidBtag(btag string) bool {
	switch {
	case !strings.Contains(btag, "#") || !strings.Contains(btag, "-"):
		return false
	case len(btag) < minBtagLength:
		return false
	}

	return true
}

func filterInvalidBtags(btags []string) []string {
	var correctBtags []string
	for _, btag := range btags {
		if isValidBtag(btag) {
			correctBtags = append(correctBtags, btag)
		}
	}
	return correctBtags
}

func urlFormatBtag(btag string) string {
	btag = strings.TrimSpace(btag)
	return strings.Replace(btag, "#", "-", -1)
}

func getRawAccountStats(btag string, stats chan string) {
	url := owAPIprefix + btag + owAPIsuffix
	fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%s\n", html)
	stats <- string(html)
}

// ConcurrentGetRawAccountStats takes an array of btag strings and
// returns an array of JSON account info but without guaranteeing same order.
func ConcurrentGetRawAccountStats(btags []string) []string {

	stats := make(chan string)
	//numBtags := len(btags)
	var btagStats []string

	for _, btag := range btags {
		go getRawAccountStats(btag, stats)
	}

	// this may overflow if some provided btags are invalid
	for i := 0; i < len(btags); i++ {
		stat := <-stats
		btagStats = append(btagStats, stat)
	}

	fmt.Println(btagStats)
	return btagStats

}

// GetEmbeddedStats takes a string arr of RawAccountStats strings and
// returns them in a formatted embed message
func GetEmbeddedStats(btags []string) *discordgo.MessageEmbed {
	emb := embed.NewEmbed().
		SetTitle("Test").
		MessageEmbed

	return emb
}
