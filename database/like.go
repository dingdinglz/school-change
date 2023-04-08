package database

// LikeCreateNew 创建新的点赞
func LikeCreateNew(change int, id string, user int) {
	_, _ = DatabaseEngine.Table(new(LikeModel)).Insert(LikeModel{
		User:   user,
		Change: change,
		Id:     id,
	})
}

// LikeHave 根据change和id以及user判断点赞是否存在
func LikeHave(change int, id string, user int) bool {
	has, _ := DatabaseEngine.Table(new(LikeModel)).Where("change = ? AND id = ? AND user = ?", change, id, user).Exist()
	return has
}

// LikeDelete 根据change和id以及user删除点赞
func LikeDelete(change int, id string, user int) {
	_, _ = DatabaseEngine.Table(new(LikeModel)).Where("change = ? AND id = ? AND user = ?", change, id, user).Delete(new(LikeModel))
}
