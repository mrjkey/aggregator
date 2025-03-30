package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mrjkey/aggregator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(bodyContent, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Items {
		feed.Channel.Items[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Items[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

// fetchFeed
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("bad bad args")
	}
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every: %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

	// feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(feed)

	// return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("not enough arguments")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	args := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
	}

	feed, err := s.db.AddFeed(context.Background(), args)
	if err != nil {
		return err
	}
	err = addFeedFollow(s, user, feed)
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

// func handlerGetFeeds(s *state, _ command) error {
// 	feeds, err := s.db.GetFeeds(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	for _, values := range feeds {
// 		fmt.Println(values.Name, values.Url, values.Name_2)
// 	}

// 	return nil
// }

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("bad, giv args")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	err = addFeedFollow(s, user, feed)
	if err != nil {
		return err
	}
	return nil
}

func addFeedFollow(s *state, user database.User, feed database.Feed) error {
	feedFollowParams := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Println(feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowersForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, follow := range follows {
		fmt.Println(follow)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("more args")
	}

	url := cmd.args[0]

	deleteParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	}

	err := s.db.DeleteFeedFollow(context.Background(), deleteParams)
	return err
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	// mark  feed as fetched
	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
		ID:            nextFeed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("For feed: %v\n", nextFeed.Name)
	// iterate over the items
	for _, item := range rssFeed.Channel.Items {
		savePost(s, item, nextFeed)
		// fmt.Println(item.Title)
		// break
	}

	fmt.Println()

	return nil
}

func savePost(s *state, item RSSItem, feed database.Feed) error {
	publishedAt, err := time.Parse(time.RFC1123, item.PubDate)
	if err != nil {
		publishedAt, err = time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			publishedAt = time.Now()
		}
	}

	params := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       item.Title,
		Url:         item.Link,
		Description: sql.NullString{String: item.Description, Valid: true},
		PublishedAt: publishedAt,
		FeedID:      feed.ID,
	}

	_, err = s.db.CreatePost(context.Background(), params)
	if err != nil {
		return err
	}

	return nil
}

func handlerBrowser(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		conv, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = conv
	}

	params := database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("%v:\n\n%v\n\n\n", post.Title, post.Description)
	}

	return nil
}
