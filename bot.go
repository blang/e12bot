package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/blang/e12bot/config"
	"github.com/blang/e12bot/discourse"
	"github.com/blang/e12bot/parsers/wiki"
	"github.com/blang/e12bot/parsing"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var cfg *config.Config
var api *discourse.API

var (
	listen     = flag.String("listen", ":8081", "addr to listen on")
	configFile = flag.String("config", "./config.json", "config file")
	logFile    = flag.String("log", "./log", "log file")
)

func init() {
	flag.Parse()
	readCfg, err := config.Parse(*configFile)
	if err != nil {
		panic("Can't read config: " + err.Error())
	}
	cfg = readCfg
	api = discourse.APIFromConfig(cfg)
}

var wg sync.WaitGroup
var parsers = &parsing.ParserCollection{}

type HTMLPrintable interface {
	HTML() string
}

func main() {
	bootstrapParsers()
	go startServer()
	for {
		fmt.Printf("Go on\n")
		wg = sync.WaitGroup{}
		wg.Add(1)
		go processTopics()
		wg.Wait()
		break
		select {
		case <-time.After(10 * time.Second):

		}
	}
}

func startServer() {
	http.HandleFunc("/slotlist", handleSlotlist)
	http.ListenAndServe(*listen, nil)
}

func bootstrapParsers() {
	parsers.Handle(&wiki.WikiParser{})
}

func processTopics() {
	defer wg.Done()
	feed, err := api.CategoryFeed(cfg.Category)
	if err != nil {
		log.Printf("Can't fetch category feed: %s", err)
		return
	}
	if feed.TopicList == nil {
		log.Printf("Can't get topiclist")
	}
	for _, t := range feed.TopicList.Topics {
		if !t.Closed || t.Closed { //TODO: Fix
			wg.Add(1)
			go processTopic(t)
		}
	}

}

func processTopic(t *discourse.DiscourseTopic) {
	defer wg.Done()
	feed, err := api.PostFeed(t.Id)
	if err != nil {
		log.Printf("Can't fetch post feed of id %d: %s", t.Id, err)
	}
	if feed.PostStream == nil {
		log.Printf("Can't get poststream")
		return
	}
	var slotlist *parsing.SlotList
	for i, p := range feed.PostStream.Posts {
		if i > 0 {
			break
		}
		for _, l := range p.Links {
			if parsers.Accept(l.Url) {

				log.Printf("Found interesting link %s on Topic %d on Post %d", l.Url, t.Id, p.Id)
				b, err := HTTPGet(ParserUrl(l.Url))
				if err != nil {
					log.Printf("Error while fetching url %s, error: %s", l.Url, err)
					continue
				}
				slotlist = parsers.Parse(string(b[:]), l.Url)

				if slotlist == nil {
					log.Printf("Error while parsing wiki url %s: %s", l.Url, err)
					continue
				}
				json, err := json.Marshal(slotlist)
				log.Printf("Slotlist found for: %s", l.Url)
				log.Printf("Slotlist: %s", json)

			} else {
				log.Printf("No Praser for link %s", l.Url)
			}
		}
	}

	wg.Add(1)
	go handleMissionTopic(feed, slotlist)
}

func handleMissionTopic(feed *discourse.DiscoursePostFeed, slotlist *parsing.SlotList) {
	log.Printf("Handle Mission Topic ID %d", feed.TopicID)
	defer wg.Done()
	for i, post := range feed.PostStream.Posts {
		if i == 0 {
			continue
		}
		log.Printf("Scan post id %d", post.Id)
		if post.Username == api.User {
			//post found
			updatePost(post, slotlist)
			return
		}
	}
	createPost(feed, slotlist)
}

func updatePost(post *discourse.DiscoursePost, slotlist *parsing.SlotList) {
	log.Printf("Update post id %d", post.Id)
	if post == nil {
		log.Printf("Post is nil")
		return
	}
	if slotlist == nil {
		log.Printf("Slostlist is nil")
		return
	}
	b, err := json.MarshalIndent(slotlist, "", " ")
	if err != nil {
		log.Printf("Marshal slotlist failed: %s", err)
		return
	}
	log.Printf("Marshall slotlist to post: %s", string(b[:]))
	//post.Id
	api.UpdatePost(post.Id, "Update slotlist", string(b[:]))
}

func createPost(feed *discourse.DiscoursePostFeed, slotlist *parsing.SlotList) {
	if slotlist == nil {
		log.Printf("Slotlist is nil")
		return
	}
	log.Printf("Create post")
	b, err := json.MarshalIndent(slotlist, "", " ")
	if err != nil {
		log.Printf("Marshal slotlist failed: %s", err)
		return
	}
	log.Printf("Marshall slotlist to post: %s", string(b[:]))
	//return
	createPost := &discourse.DiscourseCreatePost{
		TopicID:    feed.TopicID,
		CategoryID: feed.CategoryID,
		Archetype:  "regular",
		Raw:        string(b[:]),
	}
	api.CreatePost(createPost)
}

func ParserUrl(url string) string {
	if strings.Contains(url, "http://wiki.echo12.de/wiki/") {
		return url + "?action=raw"
	}
	return url
}

func HTTPGet(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %s", resp.Header["Status"][0])
	}
	// if !strings.Contains(resp.Header["Content-Type"][0], "application/json") {
	// 	return nil, fmt.Errorf("Invalid content type %s", resp.Header["Content-Type"][0])
	// }

	return body, nil
}
