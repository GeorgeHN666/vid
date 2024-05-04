package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"video-streaming/models"
)

const (
	PORTRAIT  = 0
	LANDSCAPE = 1
	SQUARE    = 2
)

type VideoLibrary struct {
	Id     int64  `json:"Id"`
	Name   string `json:"Name"`
	ApiKey string `json:"ApiKey"`
}

type VideoPlaceholder struct {
	VideoID string `json:"guid"`
}

type Video struct {
	Width            int32  `json:"width"`
	Height           int32  `json:"height"`
	ThumbnailUrl     string `json:"thumbnailUrl"`
	VideoPlaylistUrl string `json:"videoPlaylistUrl"`
	PreviewUrl       string `json:"previewUrl"`
}

func CreateNewStramLibrary(name string) (*VideoLibrary, error) {

	url := "https://api.bunny.net/videolibrary"

	var payload struct {
		Name string `json:"Name"`
	}
	payload.Name = name

	pay, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(pay)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", MASTER_API_KEY)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var VL VideoLibrary

	err = json.NewDecoder(res.Body).Decode(&VL)
	if err != nil {
		return nil, err
	}

	return &VL, nil
}

// UploadVideosStoStream upload video to video streaming service
func UploadVideosStoStream(r *http.Request, VL *VideoLibrary, postID string) ([]*models.VideoModel, error) {

	var res []*models.VideoModel

	files := r.MultipartForm.File["content"]
	if len(files) < 1 {
		return nil, errors.New("no files were uploaded")
	}

	for _, file := range files {

		var video models.VideoModel

		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		// first create video placeholder
		VP, err := CreateVideoPlaceholder(VL, file.Filename)
		if err != nil {
			return nil, err
		}
		// store keys
		video.VideoID = VP.VideoID
		// UploadVideo
		err = UploadVideo(VL, VP.VideoID, f)
		if err != nil {
			return nil, err
		}

		vd, err := GetPlayData(VL.Id, VP.VideoID, VL.ApiKey)
		if err != nil {
			return nil, err
		}
		fmt.Println("vd: ", vd)
		video.Thumbnail = vd.ThumbnailUrl
		video.Hsl = vd.VideoPlaylistUrl
		video.Preview = vd.PreviewUrl

		res = append(res, &video)
	}

	return res, nil

}

func GetPlayData(VideoLibraryID int64, VideoID string, VLKey string) (*Video, error) {

	url := fmt.Sprintf("https://video.bunnycdn.com/library/%d/videos/%s/play", VideoLibraryID, VideoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", VLKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var videoData Video

	err = json.NewDecoder(res.Body).Decode(&videoData)
	if err != nil {
		return nil, err
	}

	return &videoData, nil
}

func DeleteVideo(libraryID int64, videoID, libraryKEY string) error {
	url := fmt.Sprintf("https://video.bunnycdn.com/library/%v/videos/%v", libraryID, videoID)

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("AccessKey", libraryKEY)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var r struct {
		Success    bool   `json:"success"`
		Message    string `json:"message"`
		StatusCode int32  `json:"statusCode"`
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return err
	}

	if !r.Success {
		return fmt.Errorf("video coulnt be deleted due to %v ::: STATUS - %v", r.Message, r.StatusCode)
	}

	return nil
}

func DeleteLibrary(libraryID int64, masterKey string) error {

	url := fmt.Sprintf("https://api.bunny.net/videolibrary/%v", libraryID)

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("AccessKey", masterKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return fmt.Errorf("LibraryCouldnt be deleted due to status code being %v", res.StatusCode)
	}

	return nil
}

func UploadVideo(VL *VideoLibrary, VideoID string, file multipart.File) error {
	url := fmt.Sprintf("https://video.bunnycdn.com/library/%v/videos/%v", VL.Id, VideoID)

	FileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	requestBody := bytes.NewReader(FileContent)

	req, err := http.NewRequest("PUT", url, requestBody)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", VL.ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func CreateVideoPlaceholder(VL *VideoLibrary, videoName string) (*VideoPlaceholder, error) {

	var payload struct {
		Title string `json:"title"`
	}

	payload.Title = videoName

	url := fmt.Sprintf("https://video.bunnycdn.com/library/%v/videos", VL.Id)

	pay, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(pay)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", VL.ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var VP VideoPlaceholder

	err = json.NewDecoder(res.Body).Decode(&VP)
	if err != nil {
		return nil, err
	}

	return &VP, nil

}

func GetOrientation(width, height int32) int {

	fmt.Println("width: ", width, " height: ", height)

	if width < height {
		return PORTRAIT
	} else if width == height {
		return SQUARE
	}

	return LANDSCAPE
}
