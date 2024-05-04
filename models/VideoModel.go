package models

type VideoModel struct {
	VideoID     string `json:"v_id" bson:"v_id"`
	Orientation int    `json:"orientation" bson:"orientation"`
	Preview     string `json:"preview"  bson:"preview"`
	Hsl         string `json:"hls" bson:"hsl"`
	Thumbnail   string `json:"thumbnail" bson:"thumbanail"`
}
