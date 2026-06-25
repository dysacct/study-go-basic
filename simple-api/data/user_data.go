package data

import "simple-api/models"

// users 是硬编码的假数据，模拟数据库
var users = []models.User{
	{ID: 1, Name: "John", Email: "john@gmail.com"},
	{ID: 2, Name: "Jane", Email: "jane@gmail.com"},
	{ID: 3, Name: "Jim", Email: "jim@gmail.com"},
}

// GetAllUsers 返回所有用户
func GetAllUsers() []models.User {
	return users
}

// GetUserByID 根据ID 返回单个用户, 如果不存在返回nil
func GetUserByID(id int) *models.User {
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}
