package discourse

import (
	"github.com/blang/e12bot/config"
	"testing"
)

func freshAPI() *API {
	cfg, err := config.Parse("../config.json")
	if err != nil {
		panic("Can't get fresh api from config")
	}
	return APIFromConfig(cfg)
}

func TestCategoryFeed(t *testing.T) {
	api := freshAPI()
	feed, err := api.CategoryFeed("missionen")
	if err != nil {
		t.Fatalf("Error while fetching category feed", err)
	}
	if len(feed.TopicList.Topics) == 0 {
		t.Errorf("Topics list is empty of feed %s", feed)
	}
}

func TestPostFeed(t *testing.T) {
	api := freshAPI()
	feed, err := api.PostFeed(93, 1)
	if err != nil {
		t.Fatalf("Error while fetching category feed", err)
	}
	if len(feed.PostStream.Posts) == 0 {
		t.Errorf("Posts list is empty of feed %s", feed)
	}
	if len(feed.PostStream.Posts[0].Links) == 0 {
		t.Errorf("Link list is empty of post 0 of feed %s", feed)
	}
}
