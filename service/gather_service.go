package service

import (
	"exchange-rate-gather/config"
	"exchange-rate-gather/model"
	"exchange-rate-gather/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

// GatherExchangeRates Get the exchange rates
func GatherExchangeRates() (info []model.ExchangeRate, err error) {
	htmlStr, err := utils.DownloadHtml(config.ExchangeRateGatherUrl)
	if err != nil {
		return info, err

	}

	// 解析HTML源代码
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return info, err
	}

	section := doc.Find(config.ExchangeRateHtmlTagSelector)
	if section.Length() > 0 {
		section.Each(func(i int, s *goquery.Selection) {
			currencyDst := strings.Trim(s.Find("td:nth-child(1)").Text(), " \n\r\t")
			currencyDstDes := strings.Trim(s.Find("td:nth-child(2)").Text(), " \n\r\t")

			rateStr := strings.Trim(s.Find("td:nth-child(3)").Text(), " \n\r\t")
			rate := utils.StrToFloat64(rateStr)

			// For each item found, get the band and title
			info = append(info, model.ExchangeRate{
				CurrencySrc:            "EUR",
				CurrencyDst:            currencyDst,
				CurrencyDstDescription: currencyDstDes,
				Rate:                   rate,
				ValidMonth:             time.Now().Format("2006-01-02"),
			})
		})
	}

	return info, err
}

// GatherExchangeRatesAndSave Get the exchange rates and save to database
func GatherExchangeRatesAndSave() {
	rates, err := GatherExchangeRates()
	if err != nil {
		fmt.Println("Gather exchange rates error: ", err)
		return
	}
	dsts := viper.GetString("currency-dst")

	for _, rate := range rates {
		if strings.Contains(dsts, rate.CurrencyDst) {
			fmt.Println("Save exchange rate: ", rate)

			var count int
			err = config.DB.Get(&count, model.QueryExchangeRateExists, rate.CurrencySrc, rate.CurrencyDst, rate.ValidMonth)
			if err != nil {
				fmt.Println("Query exchange rate error: ", err)
			}
			if count > 0 {
				fmt.Println("Exchange rate has exists: ", rate)
			} else {
				_, err = config.DB.NamedExec(model.InsertExchangeRate, rate)
				if err != nil {
					fmt.Println("Save exchange rate error: ", err)
				}
			}
		}

	}
}

// GatherExchangeRatesFromNlAndSave Get the exchange rates from nl
func GatherExchangeRatesFromNlAndSave(year, month string) {
	service := ExchangeRateForNlService{
		Year:  year,
		Month: month,
	}
	rates, err := service.GetExchangeRates()
	if err != nil {
		log.Println("Gather exchange rates error: ", err)
		return
	}
	// save to database
	for _, rate := range rates {
		log.Println("Save exchange rate: ", rate)

		var count int
		err = config.DB.Get(&count, model.QueryExchangeRateExists, rate.CurrencySrc, rate.CurrencyDst, rate.ValidMonth)
		if err != nil {
			log.Println("Query exchange rate error: ", err)
		}
		if count > 0 {
			log.Println("Exchange rate has exists: ", rate)
		} else {
			_, err = config.DB.NamedExec(model.InsertExchangeRate, rate)
			if err != nil {
				log.Println("Save exchange rate error: ", err)
			}
		}
	}
}
