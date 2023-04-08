package comment

import (
	"change/database"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strconv"
	"xorm.io/xorm"
)

func getCommentDB(change int) string {
	rootPath, _ := os.Getwd()
	return filepath.Join(rootPath, "data", "comment", strconv.Itoa(change)+".db")
}

// CreateNewComment 创建一个新的评论
func CreateNewComment(i CommentModel) {
	dbConnect, _ := xorm.NewEngine("sqlite3", getCommentDB(i.Change))
	_ = dbConnect.Sync2(&CommentModel{})
	_, _ = dbConnect.Table(new(CommentModel)).Insert(i)
	_ = dbConnect.Close()
}

func GetAllComments(change int, user int) []map[string]interface{} {
	dbConnect, _ := xorm.NewEngine("sqlite3", getCommentDB(change))
	var allComments []CommentModel
	_ = dbConnect.Sync2(&CommentModel{})
	_ = dbConnect.Table(new(CommentModel)).Find(&allComments)
	_ = dbConnect.Close()
	var res []map[string]interface{}
	for _, i := range allComments {
		cnt := make(map[string]interface{})
		cnt["id"] = i.ID
		if i.Parent != "" {
			cnt["parent"] = i.Parent
		}
		cnt["created"] = i.Created.Format("2006-01-02")
		cnt["modified"] = i.Modified.Format("2006-01-02")
		cnt["content"] = i.Content
		cnt["fullname"] = i.Fullname
		cnt["profile_picture_url"] = "/avatar/" + strconv.Itoa(i.User)
		cnt["upvote_count"] = i.UpvoteCount
		if database.LikeHave(change, i.ID, user) {
			cnt["user_has_upvoted"] = true
		} else {
			cnt["user_has_upvoted"] = false
		}
		res = append(res, cnt)
	}
	return res
}
