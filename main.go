package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Meta struct {
	OgTitle            string `json:"og_title"`
	OgDescription      string `json:"og_description"`
	OgType             string `json:"og_type"`
	OgUrl              string `json:"og_url"`
	OgImage            string `json:"og_image"`
	TwitterTitle       string `json:"twitter_title"`
	TwitterDescription string `json:"twitter_description"`
	TwitterCard        string `json:"twitter_card"`
	TwitterImage       string `json:"twitter_image"`
}

func metaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	meta := Meta{}
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "twitter:title" {
			meta.TwitterTitle, _ = s.Attr("content")
		} else if name == "twitter:description" {
			meta.TwitterDescription, _ = s.Attr("content")
		} else if name == "twitter:card" {
			meta.TwitterCard, _ = s.Attr("content")
		} else if name == "twitter:image" {
			meta.TwitterImage, _ = s.Attr("content")
		}

		property, _ := s.Attr("property")
		if property == "og:title" {
			meta.OgTitle, _ = s.Attr("content")
		} else if property == "og:description" {
			meta.OgDescription, _ = s.Attr("content")
		} else if property == "og:image" {
			meta.OgImage, _ = s.Attr("content")
		} else if property == "og:type" {
			meta.OgType, _ = s.Attr("content")
		} else if property == "og:url" {
			meta.OgUrl, _ = s.Attr("content")
		}
	})

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(meta)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprint(w, buf.String())
	if err != nil {
		return
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	http.HandleFunc("/meta", metaHandler)

	log.Println("Web サーバを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
