package database

// HaveSetting 获取指定设置是否存在
func HaveSetting(name string) bool {
	cnt, _ := DatabaseEngine.Table(new(SettingModel)).Where("name = ?", name).Count()
	return cnt != 0
}

// WriteOrUpdateSetting 新建或更新设置
func WriteOrUpdateSetting(name string, value string) {
	if HaveSetting(name) {
		_, _ = DatabaseEngine.Table(new(SettingModel)).Where("name = ?").Update(map[string]interface{}{"value": value})
		return
	}
	_, _ = DatabaseEngine.Table(new(SettingModel)).Insert(SettingModel{
		Name:  name,
		Value: value,
	})
}

// GetDefaultSetting 获取设置的默认值
func GetDefaultSetting(name string) string {
	switch name {
	case "user":
		return "0"
	case "":
		return ""
	}
	return "none"
}

// GetSetting 取设置值
func GetSetting(name string) string {
	if !HaveSetting(name) {
		return GetDefaultSetting(name)
	}
	var res SettingModel
	_, _ = DatabaseEngine.Table(new(SettingModel)).Where("name = ?", name).Get(&res)
	return res.Value
}
