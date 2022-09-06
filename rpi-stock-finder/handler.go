package function

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	units "github.com/docker/go-units"
)

type RSS struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Category    []string `xml:"category"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	GUID        string   `xml:"guid"`
}

func (r *RSSItem) GetPubDate() (time.Time, error) {
	return time.Parse("Mon, 02 Jan 2006 15:04:05 MST", r.PubDate)
}

func Handle(w http.ResponseWriter, r *http.Request) {

	rssBytes, err := getBytesFromURL("https://rpilocator.com/feed/")

	if err != nil {
		log.Printf("error loading RSS feed: %s", err)
		http.Error(w, "error loading RSS feed", http.StatusInternalServerError)
		return
	}

	items, err := getStock(rssBytes, "US", time.Hour*24)

	if err != nil {
		http.Error(w, "error parsing RSS feed", http.StatusInternalServerError)
		return
	}

	discordURL := os.Getenv("discord_url")

	for _, item := range items {
		tmp := os.TempDir()
		path := filepath.Join(tmp, item.GUID)

		if _, err := os.Stat(path); os.IsNotExist(err) {

			// Write a lock file for the GUID to prevent the message from being processed again.
			if err := ioutil.WriteFile(path, []byte{}, os.ModePerm); err != nil {
				log.Printf("Unable to write lock file: %s %s", path, err)
			}

			if err := sendAlert(discordURL, item); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Printf("Alert sent to Discord server")
		}
	}
}

type DiscordMsg struct {
	Content string `json:"content"`
}

func sendAlert(discordURL string, item RSSItem) error {
	pubDate, _ := item.GetPubDate()

	// "Stock US RPi Zero 2 PiHut.com - 30 minutes ago"
	msg := fmt.Sprintf("%s - %s ago", item.Title, units.HumanDuration(time.Since(pubDate)))

	msgBytes, err := json.Marshal(DiscordMsg{Content: msg})
	if err != nil {
		return err
	}

	log.Println(msg)
	log.Println(msgBytes)

	// req, err := http.NewRequest(http.MethodPost, discordURL, bytes.NewBuffer(msgBytes))
	// if err != nil {
	// 	return fmt.Errorf("error with Discord URL, check it's a valid format %w", err)
	// }

	// req.Header.Add("Content-Type", "application/json")

	// res, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return fmt.Errorf("error sending alert to Discord %w", err)
	// }

	// var resBody []byte
	// if res.Body != nil {
	// 	defer res.Body.Close()
	// 	resBody, _ = ioutil.ReadAll(res.Body)
	// }

	// if res.StatusCode != http.StatusNoContent {
	// 	return fmt.Errorf("received non-204 status code from Discord %s, status: %d", resBody, res.StatusCode)
	// }

	return nil
}

// getBytesFromURL downloads bytes from a HTTP URL
func getBytesFromURL(url string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var body []byte
	if res.Body != nil {
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, string(body))
	}

	return body, nil
}

// getStock parses an RSS feed, then filters by category and
// an age duration field
func getStock(body []byte, category string, age time.Duration) ([]RSSItem, error) {

	var rssItems []RSSItem
	var p RSS

	if err := xml.Unmarshal(body, &p); err != nil {
		return rssItems, err
	}

	var found []RSSItem
	for _, item := range p.Channel.Item {
		d, err := item.GetPubDate()
		if err != nil {
			return rssItems, err
		}

		if time.Since(d) < age {

			for _, c := range item.Category {
				if c == category {
					found = append(found, item)
				}
			}
		}
	}

	return found, nil
}
