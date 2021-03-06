/*
   Copyright (c) 2016 VMware, Inc. All Rights Reserved.
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package dao

import (
	"database/sql"
	"errors"

	"github.com/vmware/harbor/models"
	"github.com/vmware/harbor/utils"

	"github.com/vmware/harbor/utils/log"
	"github.com/astaxie/beego/orm"
)

// GetUser ...
func GetUser(query models.User) (*models.User, error) {

	o := GetOrmer()

	sql := `select user_id, username, email, realname, comment, reset_uuid, salt,
		sysadmin_flag, creation_time, update_time
		from user u
		where deleted = 0 `
	queryParam := make([]interface{}, 1)
	if query.UserID != 0 {
		sql += ` and user_id = ? `
		queryParam = append(queryParam, query.UserID)
	}

	if query.Username != "" {
		sql += ` and username = ? `
		queryParam = append(queryParam, query.Username)
	}

	if query.ResetUUID != "" {
		sql += ` and reset_uuid = ? `
		queryParam = append(queryParam, query.ResetUUID)
	}

	var u []models.User
	n, err := o.Raw(sql, queryParam).QueryRows(&u)

	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	return &u[0], nil
}

// LoginByDb is used for user to login with database auth mode.
func LoginByDb(auth models.AuthModel) (*models.User, error) {
	o := GetOrmer()

	var users []models.User
	n, err := o.Raw(`select * from user where (username = ? or email = ?) and deleted = 0`,
		auth.Principal, auth.Principal).QueryRows(&users)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	user := users[0]

	if user.Password != utils.Encrypt(auth.Password, user.Salt) {
		return nil, nil
	}

	user.Password = "" //do not return the password

	return &user, nil
}

// ListUsers lists all users according to different conditions.
func ListUsers(query models.User) ([]models.User, error) {
	o := GetOrmer()
	u := []models.User{}
	sql := `select  user_id, username, email, realname, comment, reset_uuid, salt,
		sysadmin_flag, creation_time, update_time
		from user u
		where u.deleted = 0 and u.user_id != 1 `

	queryParam := make([]interface{}, 1)
	if query.Username != "" {
		sql += ` and username like ? `
		queryParam = append(queryParam, query.Username)
	}
	sql += ` order by user_id desc `

	_, err := o.Raw(sql, queryParam).QueryRows(&u)
	return u, err
}

// ToggleUserAdminRole gives a user admin role.
func ToggleUserAdminRole(userID, hasAdmin int) error {
	o := GetOrmer()
        queryParams := make([]interface{}, 1)
	sql := `update user set sysadmin_flag = ? where user_id = ?`
	queryParams = append(queryParams, hasAdmin)
	queryParams = append(queryParams, userID)
	r, err := o.Raw(sql, queryParams).Exec()
	if err != nil {
		return err
	}

	if _, err := r.RowsAffected(); err != nil {
		return err
	}

	return nil
}

// ChangeUserPassword ...
func ChangeUserPassword(u models.User, oldPassword ...string) (err error) {
	if len(oldPassword) > 1 {
		return errors.New("Wrong numbers of params.")
	}

	o := GetOrmer()

	var r sql.Result
	if len(oldPassword) == 0 {
		//In some cases, it may no need to check old password, just as Linux change password policies.
		r, err = o.Raw(`update user set password=?, salt=? where user_id=?`, utils.Encrypt(u.Password, u.Salt), u.Salt, u.UserID).Exec()
	} else {
		r, err = o.Raw(`update user set password=?, salt=? where user_id=? and password = ?`, utils.Encrypt(u.Password, u.Salt), u.Salt, u.UserID, utils.Encrypt(oldPassword[0], u.Salt)).Exec()
	}

	if err != nil {
		return err
	}
	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New("No record has been modified, change password failed.")
	}

	return nil
}

// ResetUserPassword ...
func ResetUserPassword(u models.User) error {
	o := GetOrmer()
	r, err := o.Raw(`update user set password=?, reset_uuid=? where reset_uuid=?`, utils.Encrypt(u.Password, u.Salt), "", u.ResetUUID).Exec()
	if err != nil {
		return err
	}
	count, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("No record be changed, reset password failed.")
	}
	return nil
}

// UpdateUserResetUUID ...
func UpdateUserResetUUID(u models.User) error {
	o := GetOrmer()
	_, err := o.Raw(`update user set reset_uuid=? where email=?`, u.ResetUUID, u.Email).Exec()
	return err
}

// CheckUserPassword checks whether the password is correct.
func CheckUserPassword(query models.User) (*models.User, error) {

	currentUser, err := GetUser(query)

	if err != nil {
		return nil, err
	}

	if currentUser == nil {
		return nil, nil
	}

	sql := `select user_id, username, salt from user where deleted = 0`

	queryParam := make([]interface{}, 1)

	if query.UserID != 0 {
		sql += ` and password = ? and user_id = ?`
		queryParam = append(queryParam, utils.Encrypt(query.Password, currentUser.Salt))
		queryParam = append(queryParam, query.UserID)
	} else {
		sql += ` and username = ? and password = ?`
		queryParam = append(queryParam, currentUser.Username)
		queryParam = append(queryParam, utils.Encrypt(query.Password, currentUser.Salt))
	}
	o := GetOrmer()
	var user []models.User

	n, err := o.Raw(sql, queryParam).QueryRows(&user)

	if err != nil {
		return nil, err
	}

	if n == 0 {
		log.Warning("User principal does not match password. Current:", currentUser)
		return nil, nil
	}

	return &user[0], nil
}

// DeleteUser ...
func DeleteUser(userID int) error {
	o := GetOrmer()
	_, err := o.Raw(`update user set deleted = 1 where user_id = ?`, userID).Exec()
	return err
}

// ChangeUserProfile ...
func ChangeUserProfile(user models.User) error {
	o := GetOrmer()
	if _, err := o.Update(&user, "Email", "Realname", "Comment"); err != nil {
		log.Errorf("update user failed, error: %v", err)
		return err
	}
	return nil
}

//Update user token
func ChangeUserToken(UserToken models.UserToken) (*models.User, error) {

	sql := `select u.* from user_token t inner join user u on u.user_id=t.user_id where t.md5_token = ?`

	log.Warning("select user:", sql)

	o := GetOrmer()
	var user models.User

	log.Warning("user token:", UserToken.Token)
	log.Warning("user md5 token:", UserToken.Md5Token)
	log.Warning("user id:", UserToken.UserId)

	if err := o.Raw(sql, UserToken.Md5Token).QueryRow(&user); err != nil {

		UserID := int64(UserToken.UserId)
		if err == orm.ErrNoRows {
			if UserToken.Token != "" {
				//query user info
				sql := `select * from user_token where user_id=?`
				var userToken models.UserToken
				//MD5
				//h := md5.New()
				//h.Write([]byte(UserToken.Token))
				//md5_token := hex.EncodeToString(h.Sum(nil))
				//
				//log.Warning("go md5:", md5_token)
				err = o.Raw(sql, UserToken.UserId).QueryRow(&userToken)

				if (err == orm.ErrNoRows && UserToken.UserId != 0) {
					p, err := o.Raw("insert into user_token (user_id, token,md5_token) values (?, ?, ?)").Prepare()
					if err != nil {
						return nil, err
					}
					defer p.Close()

					r,err := p.Exec(UserToken.UserId, UserToken.Token, UserToken.Md5Token)

					if err != nil {
						return nil, err
					}

					UserID,_ = r.LastInsertId()
					log.Warning("InsertID:", UserID)

				}else{
					_, err = o.Raw(`update user_token set token = ?, md5_token = ? where user_id = ?`, userToken.Token, UserToken.Md5Token, UserToken.UserId).Exec()
				}

				//query user info
				sql = `select * from user where user_id=?`

				o.Raw(sql, UserID).QueryRow(&user)

				log.Warning("user info:", sql)

				return &user,nil
			}
		}

	}

	log.Warning("userid:", user.UserID)
	log.Warning("username:", user.Username)

	return &user,nil
}

