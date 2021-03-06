package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// PlayHistory contains a user's play history.
type PlayHistory struct {
	Items    []HistoryItem `json:"items"`
	Next     string        `json:"next"`
	Limit    int           `json:"limit"`
	Endpoint string        `json:"href"`
}

// TrackContext contains metadata on the context in which the track was listened to.
type TrackContext struct {
	Type         string            `json:"type"`
	Endpoint     string            `json:"href"`
	ExternalURLS map[string]string `json:"external_urls"`
	URI          URI               `json:"uri"`
}

// HistoryItem contains the track and its metadata.
type HistoryItem struct {
	Track    SimpleTrack  `json:"track"`
	PlayedAt string       `json:"played_at"`
	Context  TrackContext `json:"context"`
}

// TopTracks contains both a list of tracks and paging information.
type TopTracks struct {
	Items    []TrackItem `json:"items"`
	Total    int         `json:"total"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Endpoint string      `json:"href"`
	Previous string      `json:"previous"`
	Next     string      `json:"next"`
}

// TrackItem contains basic info about a track.
type TrackItem struct {
	Album        AlbumInfo         `json:"album"`
	Artists      []ArtistInfo      `json:"artists"`
	DiscNumber   int               `json:"disc_number"`
	DurationMS   int               `json:"duration_ms"`
	Explicit     bool              `json:"explicit"`
	ExternalIDs  map[string]string `json:"external_ids"`
	ExternalURLs map[string]string `json:"external_urls"`
	Endpoint     URI               `json:"href"`
	ID           ID                `json:"id"`
	IsPlayable   bool              `json:"is_playable"`
	Name         string            `json:"name"`
	Popularity   int               `json:"popularity"`
	PreviewURL   string            `json:"preview_url"`
	TrackNumber  int               `json:"track_number"`
	Type         string            `json:"track"`
	URI          URI               `json:"uri"`
}

// AlbumInfo contains album information for a particular album.
type AlbumInfo struct {
	AlbumType    string            `json:"album_type"`
	ExternalURLs map[string]string `json:"external_urls"`
	Endpoint     string            `json:"href"`
	ID           ID                `json:"id"`
	Images       []Image           `json:"images"`
	Name         string            `json:"name"`
	ItemType     string            `json:"type"`
	URI          URI               `json:"uri"`
}

// TopArtists contains both a list of artists and paging information.
type TopArtists struct {
	Items    []ArtistItem `json:"items"`
	Total    int          `json:"total"`
	Limit    int          `json:"limit"`
	Offset   int          `json:"offset"`
	Endpoint string       `json:"href"`
	Previous string       `json:"previous"`
	Next     string       `json:"next"`
}

// ArtistItem contains extensive info about an artist.
type ArtistItem struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Followers    Followers         `json:"followers"`
	Genres       []string          `json:"genres"`
	Endpoint     string            `json:"href"`
	ID           ID                `json:"id"`
	Images       []Image           `json:"images"`
	Name         string            `json:"name"`
	Popularity   int               `json:"popularity"`
	Type         string            `json:"type"`
	URI          URI               `json:"uri"`
}

// ArtistInfo contains basic artist information for a particular artist.
type ArtistInfo struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Endpoint     string            `json:"href"`
	ID           ID                `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	URI          URI               `json:"uri"`
}

// CurrentUserRecentTracks returns the user's most recently played tracks in a single PlayHistory
// object. It supports up to 50 tracks in a single call with only the 50 most recent tracks available
// for each user. Requires authorization under user-read-recently-played scope.
func (c *Client) CurrentUserRecentTracks(total int) (*PlayHistory, error) {
	if total <= 0 || total > 50 {
		return nil, errors.New("CurrentUserRecentTracks supports up to 50 tracks per call")
	}
	spotifyURL := baseAddress + "me/player/recently-played?limit=" + strconv.Itoa(total)
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		fmt.Println("resp err")
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}

	var h PlayHistory
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		return nil, err
	}

	return &h, nil
}

// CurrentUserTopTracks returns the user's top tracks in a single TopTracks object.
// It supports up to 50 tracks in a single call with only the top 50 tracks available
// for each user. It also supports three different time ranges from where to fetch the
// tracks. Valid ranges include "short_term" (4 weeks), "medium_term" (6 months), and
// "long_term" (years). Requires authorization under user-top-read scope.
func (c *Client) CurrentUserTopTracks(opt *Options) (*TopTracks, error) {
	v := url.Values{}

	if opt != nil {
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Timerange != nil {
			v.Set("time_range", *opt.Timerange)
		}
	}

	spotifyURL := baseAddress + "me/top/tracks?" + v.Encode()
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}

	var t TopTracks
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// CurrentUserTopArtists returns the user's top artists in a single TopArtists object.
// It supports up to 50 artists in a single call with only the top 50 artists available
// for each user. It also supports three different time ranges from where to fetch the
// artists. Valid ranges include "short_term" (4 weeks), "medium_term" (6 months), and
// "long_term" (years). Requires authorization under user-top-read scope.
func (c *Client) CurrentUserTopArtists(opt *Options) (*TopArtists, error) {
	v := url.Values{}

	if opt != nil {
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Timerange != nil {
			v.Set("time_range", *opt.Timerange)
		}
	}

	spotifyURL := baseAddress + "me/top/artists?" + v.Encode()
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}

	var t TopArtists
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
