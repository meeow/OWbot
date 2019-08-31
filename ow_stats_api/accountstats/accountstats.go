package accountstats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"../apiconfig"

	"github.com/bwmarrin/discordgo"
)

var (
	apiPrefix        = apiconfig.Cfg.OWAPIPrefix
	apiDefaultSuffix = apiconfig.Cfg.OWAPISuffix
	apiHeroesSuffix  = apiconfig.Cfg.OWAPIHeroesSuffix
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

func urlFormatBtag(btag string) string {
	btag = strings.TrimSpace(btag)
	return strings.Replace(btag, "#", "-", -1)
}

func getRawAccountStats(btag string, heroes []string) string {
	if isValidBtag(btag) == false {
		fmt.Println("Invalid btag: ", btag)
		return ""
	}
	btag = urlFormatBtag(btag)
	url := ""

	// Form url to hit the external ow-api.com API
	if len(heroes) < 1 {
		url = apiPrefix + btag + apiDefaultSuffix
	} else {
		url = apiPrefix + btag + apiHeroesSuffix + strings.Join(heroes, ",")
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

	return string(html)
}

func substituteHero(inp string, hero string) string {
	return fmt.Sprintf(inp, hero)
}

// GetStats gets only select stats specified in keys
func GetStats(btag string, keys []string, heroes ...string) string {
	stats := getRawAccountStats(btag, heroes)
	fmt.Println("Fetched: ", stats)

	var statsMap map[string]string
	err := json.Unmarshal([]byte(stats), &statsMap)

	for _, k := range keys {

	}
}

// GetEmbeddedStats takes a string representing stats in flattened JSON
// and returns an embed struct
func GetEmbeddedStats(stats string) *discordgo.MessageEmbed {
	playerInfo := flattenStats(stats)
	btag := fmt.Sprint(playerInfo["name"])
	thirdPartyStatsPath := apiconfig.Cfg.ThirdPartyStatsPrefix + urlFormatBtag(btag) + apiconfig.Cfg.ThirdPartyStatsSuffix
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

	for _, k := range append(apiconfig.Cfg.StatsKeys, apiconfig.Cfg.HeroKeys...) {
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
