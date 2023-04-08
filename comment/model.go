package comment

import "time"

type CommentModel struct {
	ID                string    `json:"id"`
	Parent            string    `json:"parent"`
	Created           time.Time `json:"created"`
	Modified          time.Time `json:"modified"`
	Content           string    `json:"content"`
	Fullname          string    `json:"fullname"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	UpvoteCount       int       `json:"upvote_count"`
	Change            int       `json:"change"`
	User              int       `json:"user"`
}
