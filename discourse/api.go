package discourse

import (
	"encoding/json"
	"fmt"
	"github.com/blang/e12bot/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type API struct {
	Key     string
	User    string
	BaseURL string
}

func APIFromConfig(c *config.Config) *API {
	if c == nil {
		return nil
	}
	return &API{
		Key:     c.ApiKey,
		User:    c.ApiURL,
		BaseURL: c.ApiURL,
	}
}

func (api *API) Get(path string) ([]byte, error) {
	resp, err := http.Get(api.BaseURL + path + "?api_key=" + api.Key + "&api_username=" + api.User)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %s", resp.Header["Status"][0])
	}
	if !strings.Contains(resp.Header["Content-Type"][0], "application/json") {
		return nil, fmt.Errorf("Invalid content type %s", resp.Header["Content-Type"][0])
	}

	return body, nil
}

type DiscourseCategoryFeed struct {
	TopicList *DiscourseTopicList `json:"topic_list"`
}

type DiscourseTopicList struct {
	Topics []*DiscourseTopic `json:"topics"`
}

type DiscourseTopic struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Closed bool   `json:"closed"`
}

type DiscoursePostFeed struct {
	PostStream *DiscoursePostStream `json:"post_stream"`
}

type DiscoursePostStream struct {
	Posts []*DiscoursePost
}

type DiscoursePost struct {
	Id         int `json:"id"`
	PostNumber int `json:"post_number"`
	// Cooked     string               `json:"cooked"`
	Links []*DiscoursePostLink `json:"link_counts"`
}
type DiscoursePostLink struct {
	Url string `json:"url"`
}

func (api *API) CategoryFeed(categoryName string) (*DiscourseCategoryFeed, error) {
	b, err := api.Get("/category/" + categoryName + ".json")
	if err != nil {
		return nil, err
	}
	var feed DiscourseCategoryFeed
	err = json.Unmarshal(b, &feed)
	if err != nil {
		return nil, err
	}

	if feed.TopicList == nil {
		return nil, fmt.Errorf("Can't get topic list from json: %s", b)
	}

	return &feed, nil
}

func (api *API) PostFeed(categoryId int) (*DiscoursePostFeed, error) {
	b, err := api.Get("/t/" + strconv.Itoa(categoryId) + ".json")
	if err != nil {
		return nil, err
	}
	var feed DiscoursePostFeed
	err = json.Unmarshal(b, &feed)
	if err != nil {
		return nil, err
	}

	if feed.PostStream == nil {
		return nil, fmt.Errorf("Can't get post streamfrom json: %s", b)
	}

	return &feed, nil

}
