package accountstats

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"../config"
	"../crunchjson"
	"../embed"
	"../flatten"

	"github.com/bwmarrin/discordgo"
)

var (
	apiPrefix        = config.Cfg.OWAPIPrefix
	apiDefaultSuffix = config.Cfg.OWAPISuffix
	apiHeroesSuffix  = config.Cfg.OWAPIHeroesSuffix
)

const (
	minBtagLength = 7
)

func isValidBtag(btag string) bool {
	switch {
	case !(strings.Contains(btag, "#") || strings.Contains(btag, "-")):
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

//heroes can be list of multiple comma seperated heroes or just single hero
func getRawAccountStats(btag string, stats chan string, heroes string) {
	btag = urlFormatBtag(btag)
	url := ""

	if heroes == "" {
		url = apiPrefix + btag + apiDefaultSuffix
	} else {
		url = apiPrefix + btag + apiHeroesSuffix + heroes
		fmt.Println("HeroSuffix:", apiHeroesSuffix)
		fmt.Println("URL:", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	stats <- string(html)
}

// ConcurrentGetRawAccountStats takes an array of btag strings and
// returns an array of JSON account info but without guaranteeing same order.
func ConcurrentGetRawAccountStats(btags []string) []string {
	btags = filterInvalidBtags(btags)

	stats := make(chan string)
	var btagStats []string

	for _, btag := range btags {
		arr := strings.Split(btag, ",")
		tag := arr[0]

		if len(arr) > 1 {
			go getRawAccountStats(tag, stats, strings.Join(arr[1:], ","))
		} else {
			go getRawAccountStats(tag, stats, "")
		}

	}

	// this may overflow if some provided btags are invalid
	for i := 0; i < len(btags); i++ {
		stat := <-stats
		btagStats = append(btagStats, stat)
	}

	return btagStats
}

// Convert JSON response
func flattenStats(stats string) map[string]interface{} {
	statsMap := crunchjson.JSONtoMap(stats)

	flat, _ := flatten.Flatten(statsMap, "", flatten.DotStyle)

	return flat
}

// GetEmbeddedStats takes a string representing stats in flattened JSON
// and returns an embed struct
func GetEmbeddedStats(stats string) *discordgo.MessageEmbed {
	playerInfo := flattenStats(stats)
	btag := fmt.Sprint(playerInfo["name"])
	thirdPartyStatsPath := config.Cfg.ThirdPartyStatsPrefix + urlFormatBtag(btag) + config.Cfg.ThirdPartyStatsSuffix
	iconPath := fmt.Sprint(playerInfo["icon"])
	privateProfile := fmt.Sprint(playerInfo["private"])
	thumbnailPath := fmt.Sprint(playerInfo["ratingIcon"])
	if thumbnailPath == "" {
		thumbnailPath = fmt.Sprint(playerInfo["icon"])
	}

	tempEmb := embed.NewEmbed().
		SetAuthor(btag, iconPath, thirdPartyStatsPath).
		SetColor(0x00ff00).
		SetThumbnail(thumbnailPath)

	if privateProfile == "true" {
		tempEmb = tempEmb.AddField("Private Profile", ":c")
		return tempEmb.Truncate().MessageEmbed
	}

	for _, k := range append(config.Cfg.StatsKeys, config.Cfg.HeroKeys...) {
		kTemp := strings.Split(k, ">")
		key := kTemp[0]
		formattedFieldName := kTemp[1]
		value := fmt.Sprint(playerInfo[key])
		fmt.Println(key)

		if len(value) > 0 && value != "<nil>" {
			tempEmb = tempEmb.AddField(formattedFieldName, value)
		}
	}

	emb := tempEmb.Truncate().MessageEmbed

	return emb
}

// GetAllEmbeddedStats takes a string arr of RawAccountStats strings and
// returns an array of embeds
func GetAllEmbeddedStats(btags []string) []*discordgo.MessageEmbed {
	var embeds []*discordgo.MessageEmbed
	stats := ConcurrentGetRawAccountStats(btags)

	for _, stat := range stats {
		embeds = append(embeds, GetEmbeddedStats(stat))
	}

	return embeds

}
