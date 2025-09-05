package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
    Channel struct{
        Title string `xml:"title"`
        Link string `xml:"link"`
        Description string `xml:"description"`
        Item []RSSItem `xml:"item"`
    } `xml:"channel"`
}

type RSSItem struct {
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    PubDate string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context,feedURL string) (*RSSFeed,error){
    cli := http.DefaultClient
    req,err := http.NewRequestWithContext(ctx,"GET",feedURL,nil)
    if err != nil {
        return nil,err
    }
    req.Header.Set("User-Agent","gator")
    res, err := cli.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    body,err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    var feed RSSFeed
    if err := xml.Unmarshal(body,&feed); err != nil {
        return nil, err
    }

    return &feed,nil


}

func NewFeed(ctx context.Context,feedURL string) (*RSSFeed,error) {
    feed,err := fetchFeed(ctx,feedURL)
    if err != nil {
        return nil, err
    }
    feed.UnescapeFeed()
    return feed,nil
}

func (f *RSSFeed) UnescapeFeed(){
    f.Channel.Title = html.UnescapeString(f.Channel.Title)
    f.Channel.Description = html.UnescapeString(f.Channel.Description)
    for i := range f.Channel.Item {
        f.Channel.Item[i].Title = html.UnescapeString(f.Channel.Item[i].Title)
        f.Channel.Item[i].Description = html.UnescapeString(f.Channel.Item[i].Description)
    }
}

