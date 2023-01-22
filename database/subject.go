package database

// SubjectCreateNew 创建一个新的subject，这里的subject指交换物的分类
func SubjectCreateNew(name string, description string) {
	_, _ = DatabaseEngine.Table(new(SubjectModel)).Insert(&SubjectModel{
		Name:        name,
		Description: description,
	})
}

// SubjectGetNameByID 取对应subject的name
func SubjectGetNameByID(id int) string {
	var i SubjectModel
	_, _ = DatabaseEngine.Table(new(SubjectModel)).Where("id = ?", id).Get(&i)
	return i.Name
}
