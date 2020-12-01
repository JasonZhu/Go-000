package week02

import (
	"database/sql"
	"fmt"
	"pkg/errors"
)

type UserInput struct {
	ID   int
	Name string
}

type User struct {
	ID   int
	Name string
}

var db *sql.DB // 代码说明

// =================== Ctrl ============================
// main 相当于controller
// 场景：根据 UserId 修改 UserName
// 控制层: 处理异常
// ====================================================

func main() {

	var iUser UserInput
	// 这里是伪代码，提取参数
	// if err := json.Unmarshal([]byte(reqBody), &iUser); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"state": "error", "message": err.Error()})
	// 	return
	// }

	if err := BizUpdateUser(&iUser); err != nil {
		// 这里是伪代码，处理异常
		// 	c.JSON(http.StatusInternalServerError, gin.H{"state": "error", "message": err.Error()})
		// 	return
	}

	// 这里是伪代码，返回成功
	// c.JSON(http.StatusOK, gin.H{"state": "success"})

}

// =================== Biz ============================
// 业务逻辑层: 透传异常
// ====================================================

// BizUpdateUser 修改用户属性
func BizUpdateUser(in *UserInput) error {

	//检查逻辑
	user, err := DaoFindUser(in.ID)
	if err != nil {
		return err
	}

	//修改逻辑
	user.Name = in.Name
	return DaoUpdateUser(user)
}

// =================== DAO ============================
// 数据访问层: 包装异常
// ====================================================

// DaoFindUser 数据访问层代码，用于根据ID查找用户
func DaoFindUser(id int) (*User, error) {
	user := User{ID: id}
	if err := db.QueryRow(`
		SELECT name 
		FROM users 
		WHERE id = ?
	`, id).Scan(&user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, fmt.Sprintf("not found user by id %d", id))
		}
		return nil, errors.Wrap(err, "database error")
	}
	return &user, nil
}

// DaoUpdateUser 数据访问层代码，用于根据ID修改用户
func DaoUpdateUser(user *User) error {
	if _, err := db.Exec(`
		UPDATE users
		SET name = ? 
		WHERE id = ?
	`, user.Name, user.ID); err != nil {
		return errors.Wrap(err, fmt.Sprintf("update user error by id %d", user.ID))
	}
	return nil
}
