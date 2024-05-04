package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"video-streaming/db"
	"video-streaming/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MASTER_API_KEY = "97ff04cd-c56c-4153-9ccf-3607d40f81dcfd18a642-26a5-4283-b977-21cd4f821dd8"

func CreatePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var post models.Post
	POSTID := primitive.NewObjectID()
	post.ID = POSTID
	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	userID, _ := primitive.ObjectIDFromHex(r.FormValue("userID"))
	post.UserID = userID

	// Create library
	vl, err := CreateNewStramLibrary(post.ID.Hex())
	if err != nil {
		fmt.Println("here")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post.LibraryID = vl.Id
	post.LibraryKey = vl.ApiKey

	// Upload videos
	content, err := UploadVideosStoStream(r, vl, post.ID.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// save it to database
	post.Content = content

	id, err := db.InitDatabase(DATABASE, URI).InsertNewPost(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		PostID  string `json:"_id"`
	}

	res.Error = false
	res.Message = "Post successfully created"
	res.PostID = id.Hex()

	json.NewEncoder(w).Encode(res)

}

func GetPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	post, err := db.InitDatabase(DATABASE, URI).GetPost(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res struct {
		Error   bool         `json:"error"`
		Message string       `json:"message"`
		Post    *models.Post `json:"post"`
	}

	res.Error = false
	res.Message = "Post successfully founded"
	res.Post = post

	json.NewEncoder(w).Encode(res)

}

func GetFeed(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	posts, err := db.InitDatabase(DATABASE, URI).GetArrayOfPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res struct {
		Error   bool           `json:"error"`
		Message string         `json:"message"`
		Post    []*models.Post `json:"post"`
	}

	res.Error = false
	res.Message = "Posts founded"
	res.Post = posts

	json.NewEncoder(w).Encode(res)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {

	post, err := db.InitDatabase(DATABASE, URI).GetPost(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, video := range post.Content {
		err := DeleteVideo(post.LibraryID, video.VideoID, post.LibraryKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = DeleteLibrary(post.LibraryID, MASTER_API_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.InitDatabase(DATABASE, URI).DeletePost(post.ID.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	res.Error = false
	res.Message = "Post successfuly deleted"

	json.NewEncoder(w).Encode(res)

}
