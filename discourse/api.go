package discourse

import (
	"encoding/json"
	"fmt"
	"github.com/blang/e12bot/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
		User:    c.ApiUser,
		BaseURL: c.ApiURL,
	}
}

func (api *API) Get(path string, values url.Values) ([]byte, error) {
	values.Set("api_key", api.Key)
	values.Set("api_username", api.User)

	url := api.BaseURL + path + "?" + values.Encode()
	resp, err := http.Get(url)
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
	PostStream        *DiscoursePostStream `json:"post_stream"`
	TopicID           int                  `json:"id"`
	CategoryID        int                  `json:"category_id"`
	HighestPostNumber int                  `json:"highest_post_number"`
}

type DiscoursePostStream struct {
	Posts []*DiscoursePost `json:"posts"`
}

type DiscoursePost struct {
	Id         int                  `json:"id"`
	PostNumber int                  `json:"post_number"`
	Username   string               `json:"username"`
	UserID     int                  `json:"user_id"`
	Links      []*DiscoursePostLink `json:"link_counts"`
}

type DiscourseCreatePost struct {
	TopicID    int
	CategoryID int
	Archetype  string
	Raw        string
}
type DiscoursePostLink struct {
	Url string `json:"url"`
}

func (api *API) CategoryFeed(categoryName string) (*DiscourseCategoryFeed, error) {
	b, err := api.Get("/c/"+categoryName+".json", url.Values{})
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

func (api *API) PostFeed(topicId int, page int) (*DiscoursePostFeed, error) {
	values := url.Values{}
	values.Set("page", strconv.Itoa(page))
	b, err := api.Get("/t/"+strconv.Itoa(topicId)+".json", values)
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

func (api *API) CreatePost(createPost *DiscourseCreatePost) {
	values := url.Values{
		"topic_id":    {strconv.Itoa(createPost.TopicID)},
		"category":    {strconv.Itoa(createPost.CategoryID)},
		"archetype":   {createPost.Archetype},
		"raw":         {createPost.Raw},
		"post_number": {"2"},
	}
	resp, err := http.PostForm(api.BaseURL+"/posts"+"?api_key="+api.Key+"&api_username="+api.User,
		values)
	if err != nil {
		log.Printf("Error while creating post: %s", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Status code while creating post was %d", resp.StatusCode)
	}
	log.Printf("Successfully created post for topic: %d", createPost.TopicID)
}

func (api *API) UpdatePost(postID int, editReason string, content string) {
	values := url.Values{
		"post[raw]":         {content},
		"post[edit_reason]": {editReason},
	}

	url := api.BaseURL + "/posts/" + strconv.Itoa(postID) + ".json?api_key=" + api.Key + "&api_username=" + api.User
	data := values.Encode()
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		log.Printf("Error while encoding request: %s", err)
		return
	}
	req.ContentLength = int64(len(data))
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error while updating post: %s", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Status code while updating post was %d", resp.StatusCode)
		log.Printf("Response: %s", resp)
		contents, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Response Body: %s", contents)
	} else {
		log.Printf("Successfully updated post for postid: %d", postID)
	}
}
