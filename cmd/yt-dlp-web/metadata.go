package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bogem/id3v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type DeezerResponse struct {
	Data []struct {
		Title  string `json:"title"`
		Artist struct {
			Name          string `json:"name"`
			PictureMedium string `json:"picture_medium"`
		} `json:"artist"`
		Album struct {
			Title       string `json:"title"`
			CoverMedium string `json:"cover_medium"`
		} `json:"album"`
	} `json:"data"`
}

func UpdateFileMetadata(filePath string) string {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		log.Println("Cant open file", filePath, err)
		return "Cant open file"
	}
	defer tag.Close()

	artist := tag.Artist()
	title := tag.Title()
	title = sanitizeTitle(artist, title)

	deezerResponse, err := getMetadata(artist, title)
	if err != nil {
		if strings.Contains(title, " - ") { // t:Scooter - Fire -> t:Scooter a:Fire
			kv := strings.Split(title, " - ")
			artist = kv[0]
			title = kv[1]
			if len(title) > 0 && len(artist) > 0 {
				deezerResponse, err = getMetadata(artist, title)
			}
		}
		if err != nil {
			tag.SetAlbum(title)
			tag.Save()
			// poor mans fallback
			return "Error getting metadata: " + err.Error()
		}
	}

	trackData := deezerResponse.Data[0]
	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	tag.SetArtist(trackData.Artist.Name)
	tag.SetTitle(trackData.Title)
	tag.SetAlbum(trackData.Album.Title)

	imageURL := trackData.Album.CoverMedium
	if imageURL == "" {
		imageURL = trackData.Artist.PictureMedium
	}
	updateImage(imageURL, tag)

	metaText := fmt.Sprintf("%s / %s - %s cover: %s", trackData.Album.Title, trackData.Artist.Name, trackData.Title, imageURL)
	if err := tag.Save(); err != nil {
		log.Println("Cant save tag "+metaText, err)
		return "Cant save tag " + metaText + " : " + err.Error()
	} else {
		log.Println("Updated metadata: " + metaText)
		return "Updated metadata: " + metaText
	}
}

func sanitizeTitle(artist, title string) string {
	sanitizedTitle := strings.Replace(title, artist+" - ", "", -1)
	sanitizedTitle = strings.Replace(sanitizedTitle, artist, "", -1)
	re := regexp.MustCompile(`[(\[].*?[)\]]`)
	sanitizedTitle = re.ReplaceAllString(sanitizedTitle, "")
	sanitizedTitle = strings.TrimSpace(sanitizedTitle)
	return sanitizedTitle
}

func getMetadata(artist string, title string) (*DeezerResponse, error) {
	if title == "" || artist == "" {
		return nil, errors.New("empty title or artist")
	}
	log.Printf("Searching metadata for: %s - %s\n", artist, title)
	body, err := httpGET("https://api.deezer.com/search?q=" + url.QueryEscape(artist+" - "+title))
	if err != nil {
		return nil, errors.New("could not call deezer: " + err.Error())
	}
	var deezerResponse DeezerResponse
	if err := json.Unmarshal(body, &deezerResponse); err != nil {
		return nil, errors.New("could not unmarshal deezer response")
	}
	if len(deezerResponse.Data) == 0 {
		return nil, errors.New("empty response (nothing found)")
	}
	return &deezerResponse, nil
}

func updateImage(imageURL string, tag *id3v2.Tag) {
	if imageURL == "" {
		log.Println("Empty image url")
		return
	}
	imageData, err := httpGET(imageURL)
	if err != nil {
		log.Println("Could not download image", err)
		return
	}
	if err == nil {
		pic := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/jpeg",
			PictureType: id3v2.PTFrontCover,
			Description: "Front cover",
			Picture:     imageData,
		}
		tag.AddAttachedPicture(pic)
		log.Println("Updated cover with " + imageURL)
	} else {
		log.Println("Could not download image", err)
	}
}

func httpGET(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received %d response code", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
