package main

import (
	"autosellManager/Structs"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	settings         Structs.Config
	blacklistedItems []int64
	userId           int64
	collections      []Structs.Collections
)

func main() { //two methods, BEST_PRICE or RAP
	//self-explanatory
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(configFile, &settings)

	blFile, err := os.Open("blacklisted.txt")
	blScanner := bufio.NewScanner(blFile)

	for blScanner.Scan() {
		iVal, err := strconv.ParseInt(blScanner.Text(), 10, 64)
		if err != nil {
			continue
		}
		blacklistedItems = append(blacklistedItems, iVal)
	}

	userId = GetUserId(settings.Cookie)
	if userId == -1 {
		log.Fatal("cookie is invalid")
	}

	colTemp := GetCollections(userId, "", []Structs.Collections{})

	for i := 0; i < len(colTemp); i++ {
		blFlag := false
		for j := 0; j < len(blacklistedItems); j++ {
			if colTemp[i].AssetId == blacklistedItems[j] {
				blFlag = true
				break
			}
		}

		if blFlag {
			continue
		}

		collections = append(collections, colTemp[i])
	}

	if settings.Method == "RESET" {
		xsrf := GetXsrf(settings.Cookie)

		if xsrf == "" {
			log.Fatal("could not get xsrf")
		}

		for i := 0; i < len(collections); i++ {
			sellFlags := TakeOffSale(settings.Cookie, xsrf, collections[i].AssetId, collections[i].UserAssetId)
			if sellFlags {
				fmt.Printf("Taken %s offsale\n", collections[i].Name)
			}
		}

		return
	}

	if settings.Method == "BEST_PRICE" {
		for i := 0; i < len(collections); i++ {
			priceFlags := GetBestPrice(settings.Cookie, collections[i].AssetId)
			if priceFlags == -1 {
				log.Fatal("an error occured when grabbing price") //i can't be bothered to remove the asset that errors out soo...
			}

			if priceFlags == -2 {
				fmt.Println("a rate limit has occured, sending another request in 10 seconds...")
				time.Sleep(10 * time.Second)
				i--
				continue
			}

			collections[i].SellPrice = priceFlags - 1
		}
	}

	if settings.Method == "RAP" {
		for i := 0; i < len(collections); i++ {
			collections[i].SellPrice = int(float64(collections[i].RecentAveragePrice) * settings.RapMultiplier)
		}
	}

	for i := 0; i < len(collections); i++ {
		fmt.Println(collections[i].Name, collections[i].SellPrice)
	}

	fmt.Print("Are you ready to sell (y/n): ")

	stdReader := bufio.NewReader(os.Stdin)
	iFlag, _, _ := stdReader.ReadRune()

	if iFlag == 'y' {
		xsrf := GetXsrf(settings.Cookie)

		if xsrf == "" {
			log.Fatal("could not get xsrf")
		}

		for i := 0; i < len(collections); i++ {
			sellFlag := SellItem(settings.Cookie, xsrf, collections[i].AssetId, collections[i].UserAssetId, collections[i].SellPrice) //there's probably a rate limit to this api, haven't checked it out

			if sellFlag {
				fmt.Printf("Successfully sold %s at %d robux\n", collections[i].Name, collections[i].SellPrice)
				continue
			}

			fmt.Printf("Could not sell %s\n", collections[i].Name)
		}
	}
}
