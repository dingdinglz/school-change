package database

// SubjectCreateNew 创建一个新的subject，这里的subject指交换物的分类
func SubjectCreateNew(name string, description string) {
	_, _ = DatabaseEngine.Table(new(SubjectModel)).Insert(&SubjectModel{
		Name:        name,
		Description: description,
	})
}
