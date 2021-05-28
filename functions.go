package main

import (
	"autosellManager/Structs"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	httpClient = http.Client{}
)

func GetUserId(cookie string) int64 {
	req, err := http.NewRequest("GET", "https://users.roblox.com/v1/users/authenticated", nil)
	if err != nil {
		return -1
	}

	req.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: cookie})

	res, err := httpClient.Do(req)
	if err != nil {
		return -1
	}

	if res.StatusCode == 401 {
		return -1
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1
	}

	var uid Structs.UserId
	err = json.Unmarshal(body, &uid)

	return uid.ID
}

func GetCollections(userid int64, cursor string, collections []Structs.Collections) []Structs.Collections {

	req, err := http.NewRequest("GET", fmt.Sprintf("https://inventory.roblox.com/v1/users/%d/assets/collectibles?assetType=null&cursor=%s&limit=100&sortOrder=Desc", userid, cursor), nil)
	if err != nil {
		return []Structs.Collections{}
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return []Structs.Collections{}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Structs.Collections{}
	}

	var tradeCollections Structs.TradeCollections
	err = json.Unmarshal(body, &tradeCollections)

	for i := 0; i < len(tradeCollections.Data); i++ {
		collections = append(collections, Structs.Collections{
			Name:               tradeCollections.Data[i].Name,
			UserAssetId:        tradeCollections.Data[i].UserAssetID,
			AssetId:            tradeCollections.Data[i].AssetID,
			RecentAveragePrice: tradeCollections.Data[i].RecentAveragePrice,
		})
	}

	if tradeCollections.NextPageCursor != "" {
		return GetCollections(userid, tradeCollections.NextPageCursor, collections)
	}

	return collections
}

func GetBestPrice(cookie string, assetid int64) int { //just putting it out here, it's a big meme how it needs a cookie | -1 = error, -2 = toomanyrequests
	req, err := http.NewRequest("GET", fmt.Sprintf("https://economy.roblox.com/v1/assets/%d/resellers?cursor=&limit=10", assetid), nil)
	if err != nil {
		return -1
	}

	req.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: cookie})

	res, err := httpClient.Do(req)
	if err != nil {
		return -1
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1
	}

	if strings.Contains(string(body), "TooManyRequests") { //probably make this recursive? idk
		return -2
	}

	if res.StatusCode != 200 {
		return -1 //would only happen if cookie is invalid... but cookie is backed by userId func so idk if someone is that special to even be able to do that
	}

	var resellers Structs.Resellers
	err = json.Unmarshal(body, &resellers)
	if err != nil {
		return -1
	}

	if len(resellers.Data) == 0 {
		return -1 //no one is selling this limited or some shit and i'm too lazy to deal with multiple error codes so -1 here bitch
	}

	return resellers.Data[0].Price
}

func GetXsrf(cookie string) string {
	req, err := http.NewRequest("GET", "https://www.roblox.com/transactions", nil)
	if err != nil {
		return ""
	}

	req.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: cookie})

	res, err := httpClient.Do(req)
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	xsrf, _ := regexp.Compile("data-token=\"(.*?)\" />")

	if !xsrf.Match(body) {
		return ""
	}

	return xsrf.FindStringSubmatch(string(body))[1]
}

func SellItem(cookie string, xsrf string, assetid int64, uaid int64, price int) bool { //does not account for token invalidation
	req, err := http.NewRequest("POST", "https://www.roblox.com/asset/toggle-sale", bytes.NewBufferString(fmt.Sprintf("assetId=%d&userAssetId=%d&price=%d&sell=true", assetid, uaid, price)))
	if err != nil {
		return false
	}

	req.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: cookie})
	req.Header.Add("x-csrf-token", xsrf)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)
	if err != nil {
		return false
	}

	//return strings.Contains(string(body), "true") //i know i know, i should parse it against a struct but i'm too lazy

	if res.StatusCode != 200 {
		return false
	}

	return true //better solution to check statuscode so it can "handle" shit if people want to add that
}
