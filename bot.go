package main

import (
	_ "expvar"
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

type TopicNextUpdate struct {
	m map[int]time.Time
	sync.RWMutex
}

const (
	UpdateInterval    = time.Minute * 5
	ExtUpdateInterval = time.Minute * 20
	CheckInterval     = time.Minute * 1
)

var topicNextUpdate = &TopicNextUpdate{m: make(map[int]time.Time)}

var (
	listen     = flag.String("listen", ":8081", "addr to listen on")
	configFile = flag.String("config", "./config.json", "config file")
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
		// break
		select {
		case <-time.After(CheckInterval):

		}
	}
}

func startServer() {
	http.HandleFunc("/slotlist", handleSlotlist)
	http.ListenAndServe(*listen, nil)
}

func bootstrapParsers() {
	parsers.Handle(&wiki.WikiTableParser{})
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
		return
	}
	for _, t := range feed.TopicList.Topics {
		if !t.Closed {
			topicNextUpdate.RLock()
			nextUpdate, found := topicNextUpdate.m[t.Id]
			topicNextUpdate.RUnlock()
			if found && time.Now().Before(nextUpdate) {
				log.Printf("Topic %d does not need an update yet", t.Id)
				continue
			}
			wg.Add(1)
			go processTopic(t)
		} else {
			// Cleanup map
			topicNextUpdate.Lock()
			delete(topicNextUpdate.m, t.Id)
			topicNextUpdate.Unlock()
		}
	}

}

func processTopic(t *discourse.DiscourseTopic) {
	defer wg.Done()
	if t == nil {
		log.Printf("Topic nil")
		return
	}
	if !(t.Id > 0) {
		log.Printf("Topic ID wrong")
		return
	}
	feed, err := api.PostFeed(t.Id, 1)
	if err != nil {
		log.Printf("Can't fetch post feed of id %d: %s", t.Id, err)
		return
	}
	if feed == nil {
		log.Println("Can't get feed")
		return
	}
	if feed.PostStream == nil {
		log.Printf("Can't get poststream")
		return
	}
	var slotlist *parsing.SlotList

	if len(feed.PostStream.Posts) > 0 {

		url := ""
		//Search first post for links
		p := feed.PostStream.Posts[0]
		for _, l := range p.Links {
			if parsers.Accept(l.Url) {

				// log.Printf("Found interesting link %s on Topic %d on Post %d", l.Url, t.Id, p.Id)
				b, err := HTTPGet(ParserUrl(l.Url))
				if err != nil {
					log.Printf("Error while fetching url %s, error: %s", l.Url, err)
					continue
				}
				slotlist = parsers.Parse(string(b[:]), l.Url)

				if slotlist == nil {
					log.Printf("Nil slotlist while parsing url %s: %s", l.Url, err)
					continue
				}
				url = l.Url
				log.Printf("Slotlist found for %s, slotgroups: %d", l.Url, len(slotlist.SlotListGroups))
			} else {
				log.Printf("No Praser for link %s", l.Url)
			}
		}

		isExternal := IsExternal(url)

		topicID := feed.TopicID
		categoryID := feed.CategoryID

		//find bot post
		botpost := findBotPost(feed)
		if botpost > 0 {
			updatePost(topicID, botpost, slotlist, isExternal)
		} else {
			createPost(topicID, categoryID, slotlist, isExternal)
		}

	} else {
		log.Printf("No posts on stream for topic %d", feed.TopicID)
		return
	}
}

//Sideeffect: Changes feed
func findBotPost(feed *discourse.DiscoursePostFeed) int {
	page := 1
	var err error
	for feed != nil && feed.PostStream != nil && len(feed.PostStream.Posts) > 0 {
		for _, p := range feed.PostStream.Posts {
			if p.Username == api.User {
				return p.Id
			}
		}
		page += 1
		topicID := feed.TopicID
		feed, err = api.PostFeed(feed.TopicID, page)
		if err != nil {
			log.Printf("Can not get feed page %d for topic %d", page, topicID)
			return 0
		}
	}
	return 0
}

func updatePost(topicID int, postID int, slotlist *parsing.SlotList, isExternal bool) {
	log.Printf("Update post id %d on topic %d", postID, topicID)
	if postID == 0 {
		log.Printf("Post is not correct")
		return
	}

	slotListStr := EncodeSlotList(slotlist)
	api.UpdatePost(postID, "Update slotlist", slotListStr)
	topicNextUpdate.Lock()
	if isExternal {
		topicNextUpdate.m[topicID] = time.Now().Add(ExtUpdateInterval)
	} else {
		topicNextUpdate.m[topicID] = time.Now().Add(UpdateInterval)
	}
	topicNextUpdate.Unlock()
}

func createPost(topicID int, categoryID int, slotlist *parsing.SlotList, isExternal bool) {
	log.Printf("Create post for topic %d", topicID)
	if topicID == 0 || categoryID == 0 {
		log.Printf("Topicid %d or categoryid %d wrong", topicID, categoryID)
		return
	}

	slotListStr := EncodeSlotList(slotlist)
	createPost := &discourse.DiscourseCreatePost{
		TopicID:    topicID,
		CategoryID: categoryID,
		Archetype:  "regular",
		Raw:        slotListStr,
	}

	api.CreatePost(createPost)
	topicNextUpdate.Lock()
	if isExternal {
		topicNextUpdate.m[topicID] = time.Now().Add(ExtUpdateInterval)
	} else {
		topicNextUpdate.m[topicID] = time.Now().Add(UpdateInterval)
	}
	topicNextUpdate.Unlock()
}

func IsExternal(url string) bool {
	return strings.Contains(url, "heeresgruppe2012.de")
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

	return body, nil
}
