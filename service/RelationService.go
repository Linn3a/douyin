package service

import (
	"douyin/models"
	// "errors"
	// "gorm.io/gorm"
	"fmt"
	"strings"
	"strconv"
	"douyin/middleware/rabbitmq"

)

const USER_TABLE_NAME = "users"
const RELATION_TABLE_NAME = "user_follows"


// redis 查询优化
func GetUserFollowCount(u *models.UserInfo) error {
	followCount, err := models.RedisClient.SCard(RedisCtx, SOCIAL_FOLLOWING_KEY+strconv.Itoa(int(u.ID))).Result()
	if err != nil {
		return fmt.Errorf("social following set check error: %v", err)

	}
	u.FollowCount = followCount
	return nil
}

func GetUserFollowerCount(u *models.UserInfo) error {
	followerCount, err := models.RedisClient.SCard(RedisCtx, SOCIAL_FOLLOWER_KEY+strconv.Itoa(int(u.ID))).Result()
	if err != nil {
		return fmt.Errorf("social following set check error: %v", err)
	}
	u.FollowerCount = followerCount
	return nil
}

func GetUserIsFollow(u *models.UserInfo, from_id uint) error {
	isFollowed, err := models.RedisClient.SIsMember(RedisCtx, SOCIAL_FOLLOWING_KEY+strconv.Itoa(int(from_id)), u.ID).Result()
	if err != nil {
		return fmt.Errorf("social following set check error: %v", err)
	}
	u.IsFollow = isFollowed
	return nil
}

//---------------------------------



// FollowAction 关注操作
func FollowAction(fromId uint, toId uint) error {
	// 关注消息加入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(fromId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(toId)))
	rabbitmq.RmqFollowAdd.Publish(sb.String())
	fmt.Println("关注消息入队成功")
	return nil
}


func CancleAction(fromId uint,toId uint)error{
	// 取关消息加入消息列表
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(fromId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(toId)))
	rabbitmq.RmqFollowDel.Publish(sb.String())
	// 记录日志
	fmt.Println("取关消息入队成功")
	return nil
		
}

// HasRelation fromId(follower_id关注者) 是否关注 toId（followed_id被关注者）
func HasRelation(fromId uint, toId uint) bool {
	tmp := models.Relation{}
	var cnt int64
	if err := models.DB.Table("user_follows").
		Where("follower_id = ? AND followed_id = ?", fromId, toId).Find(&tmp).Count(&cnt).
		Error; err != nil{ //没有该条记录
		return false
	}
	if cnt == 0{
		return false
	}
	return true
}

// CreateRelation 创建关注
func CreateRelation(fromId uint, toId uint) error {

	relation := models.Relation{
		FollowedId:   toId,
		FollowerId: fromId,
	}
	
	if err := models.DB.Table("user_follows").Create(&relation).Error; err != nil { //创建记录
			return err
	}
	return nil
}

// DeleteRelation 取消关注
func DeleteRelation(fromId uint, toId uint) error {

	relation := models.Relation{
		FollowedId:   toId,
		FollowerId: fromId,
	}
	if HasRelation(fromId,toId){
		err := models.DB.Table("user_follows").Delete(models.Relation{},&relation).Error
		return err
	}
	return nil
}

// FollowList 获取关注表
func FollowList(Id uint) ([]models.UserInfo, error) {
	var userList []models.User
	var userInfoList []models.UserInfo
	// RELATION_TABLE_NAME：FollowerId（关注者），FollowedId（被关注者）
	//该User在 RELATION表中作为FollowerId，需获取对应的所有FollowedId的Users列表
	if err := models.DB.Model(&models.User{}).
		Joins("left join "+RELATION_TABLE_NAME+" on "+USER_TABLE_NAME+".id = "+RELATION_TABLE_NAME+".followed_id").
		Where(RELATION_TABLE_NAME+".follower_id=?", Id).
		Scan(&userList).Error; err == nil {
			// TODO: add info from other service
			userInfoList := make([]models.UserInfo,len(userList))
			for i,u :=range userList{
				userInfoList[i] = GenerateUserInfo(&u)
				userInfoList[i].IsFollow = true
			}
		return userInfoList, nil
	}else{
		return userInfoList,err
	}

}

// FollowerList  获取粉丝表
func FollowerList(Id uint) ([]models.UserInfo, error) {
	var userList []models.User
	var userInfoList []models.UserInfo

	if err := models.DB.Model(&models.User{}).
		Joins("left join "+RELATION_TABLE_NAME+" on "+USER_TABLE_NAME+".id = "+RELATION_TABLE_NAME+".follower_id").
		Where(RELATION_TABLE_NAME+".followed_id=?", Id).
		Scan(&userList).Error; err == nil {
			// TODO: add info from other service
			userInfoList := make([]models.UserInfo,len(userList))
			for i,u :=range userList{
				userInfoList[i] = GenerateUserInfo(&u)
				userInfoList[i].IsFollow = HasRelation(Id,u.ID)
			}
		return userInfoList, nil
	}else{
		return userInfoList, err
	}
	
}

// FriendList 获取朋友列表（互相关注）
func FriendList(Id uint) ([]models.UserInfo, error) {
	var friendList []models.UserInfo
	// 查询 Id 的粉丝列表
	// 检查 粉丝列表中的用户是否Id也关注 
	followerList, err := FollowerList(Id)
	if err != nil {
		return friendList, err
	} else {
		for _, userInfo := range followerList {
			if userInfo.IsFollow {
				friendList = append(friendList, userInfo)
			}
		}
		return friendList, nil
	}
}

//获取User的关注总数
func GetFollowCnt(Id uint)(int64,error){
	var cnt int64
	if err := models.DB.
		Table("user_follows").
		Where("follower_id=?", Id).
		Count(&cnt).Error; err != nil {
			fmt.Println("cnt:",cnt)
		return 0,err
	}
	return cnt,nil

}
//获取User的粉丝总数
func GetFollowerCnt(Id uint)(int64,error){
	var cnt int64
	if err := models.DB.
		Table("user_follows").
		Where("followed_id=?", Id).
		Count(&cnt).Error; err != nil {
			fmt.Println("cnt:",cnt)
		return 0,err
	}
	return cnt,nil

}