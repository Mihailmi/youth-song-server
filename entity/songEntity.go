package entity

type Song struct {
	song_Number string   `bson:"number"`
	song_Lyrics []string `bson:"lyrics"`
	song_Title  string   `bson:"title"`
	song_pptx   string   `bson:"pptx"`
	file_type   string   `bson:"type"`
	file_name   string   `bson:"name"`
}
