package proxy

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/prometheus/common/log"
	"golang.org/x/net/proxy"
)

const (
	X_USER_AGENT      = "<>"
	USER_AGENT        = "<>"
	ORIGIN            = "<>"
	ACCEPT_LANG       = "<>"
	JSON_CONTENT_TYPE = "<>"

	
)


//You need to start Tor instances and then you can pass there port numbers in the map
//I have started 5 tor instances with ports (9060,9062,9064,9068)
func getSocksPort() string {

	socketMap := make(map[int]string)

	socketMap[0] = "9060"
	socketMap[1] = "9062"
	socketMap[2] = "9064"
	socketMap[3] = "9066"
	socketMap[4] = "9068"

	return socketMap[rand.Intn(5)]

}

func CallURL(lUrl string) (*goquery.Document, error) {

	if !strings.HasPrefix(lUrl, "https") {
		if strings.HasPrefix(lUrl, "/") {
			lUrl = ORIGIN + lUrl
		} else {
			lUrl = ORIGIN + "/" + lUrl
		}

	}

	tbProxyURL, err := url.Parse("socks5://localhost:" + getSocksPort())

	if err != nil {

		fmt.Println(err)
		log.Error(err)

	}

	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		fmt.Println(err)
		log.Error(err)

	}

	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}
	if err != nil {
		fmt.Println(err)
		log.Error(err)
	}

	req, err := http.NewRequest("GET", strings.TrimSpace(lUrl), nil)
	if err != nil {

		fmt.Println(time.Now().String() + ":Error Creating Request :" + err.Error())
		log.Error("[AMZ]:", time.Now().String()+":Error Creating Request :"+err.Error())
		return nil, err

	}

	req.Header.Set("X-User-Agent", X_USER_AGENT)
	req.Header.Set("Origin", ORIGIN)
	req.Header.Set("Accept-Language", ACCEPT_LANG)
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", JSON_CONTENT_TYPE)
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)

	if err != nil {

		if err == io.EOF {

			return nil, err

		}

		fmt.Println(time.Now().String() + ":Error While Requesting :" + err.Error())
		log.Error("[AMZ]:", time.Now().String()+":Error While Requesting :"+err.Error())

		return nil, err

	}

	if resp != nil {

		doc, err := goquery.NewDocumentFromResponse(resp)

		if err != nil {

			if err == io.EOF {

				return nil, err

			}

			fmt.Println(time.Now().String() + ":Error in Converting Request to Doc :" + err.Error())
			log.Error("[AMZ]:", time.Now().String()+":Error in Converting Request to Doc :"+err.Error())

			return nil, err

		}

		return doc, nil

	}

	return nil, err

}
