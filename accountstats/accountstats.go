package accountstats

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"../crunchjson"
	"../embed"

	"github.com/bwmarrin/discordgo"
)

var (
	owAPIprefix = "https://ow-api.com/v1/stats/pc/us/"
	owAPIsuffix = "/complete"
)

const (
	minBtagLength = 7
)

func isValidBtag(btag string) bool {
	//fmt.Println("Testing ", btag)

	switch {
	case !(strings.Contains(btag, "#") || strings.Contains(btag, "-")):
		return false
	case len(btag) < minBtagLength:
		return false
	}

	//fmt.Println(btag, "is valid!")
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
	btag = urlFormatBtag(btag)
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
	btags = filterInvalidBtags(btags)

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

// Convert JSON response
func pruneStats(stats string) map[string]interface{} {

	statsMap := crunchjson.JSONtoMap(stats)
	fmt.Println(statsMap["name"])

	return statsMap
}

// GetEmbeddedStats takes a string arr of RawAccountStats strings and
// returns them in a formatted embed message
func GetEmbeddedStats(btags []string) *discordgo.MessageEmbed {

	//debug
	stats := ConcurrentGetRawAccountStats(btags)
	prunedStats := pruneStats(stats[0])
	fmt.Println(prunedStats)
	firstPlayer := prunedStats

	tempEmb := embed.NewEmbed().
		SetTitle(fmt.Sprint(firstPlayer["name"]))

	emb := tempEmb.MessageEmbed

	return emb
}
