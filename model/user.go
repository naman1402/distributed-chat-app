package model

type User struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
}

type LoginReq struct {
	Id string `json:"id"`
}

// func CreateUser(userId string, username string) {
// 	query := `INSERT INTO users(id, username) VALUES(?, ?)`
// 	database.ExecuteQuery(query, userId, username)
// }

// func CheckIfUserExist(userid string) (string, string) {
// 	query := `SELECT id, username FROM users WHERE username = ?`
// 	ID, username := database.CheckIfExist(query, userid)
// 	return ID, username
// }
