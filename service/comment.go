package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/dao/redis"
	"BytesDanceProject/model"
	"BytesDanceProject/tool"
	"fmt"
	"time"
)

// CreateComment 创建评论
func CreateComment(userId int, videoId int, commentText string) (*model.Comment, error) {

	now := time.Now()
	time := time.Unix(now.Unix(), 0) // 参数分别是：秒数,纳秒数

	NewComment := model.Comment{
		UserID:     userId,
		VideoID:    videoId,
		Content:    commentText,
		CreateDate: time,
		IsDeleted:  0,
		UpdateDate: time,
	}

	//将评论存入MySQL中
	err := mysql.InsertComment(&NewComment)
	if err != nil {
		return nil, err
	}

	//将评论存入Redis中
	key := tool.GetVideoCommentKey(videoId)
	err = redis.AddCommentToSortedSet(key, now.Unix(), &NewComment)
	if err != nil {
		return nil, err
	}

	return &NewComment, nil
}

// DeleteComment 删除评论
func DeleteComment(commentId int) error {

	comment, err := mysql.GetComment(commentId)
	if err != nil {
		return err
	}
	fmt.Println(comment)
	key := tool.GetVideoCommentKey(comment.VideoID)

	err = redis.RemoveComment(key, comment)
	if err != nil {
		return err
	}

	//修改mysql中评论的状态
	err = mysql.UpdateCommentStatus(commentId)
	if err != nil {
		return err
	}

	return nil
}

// ListComment 获取videoId的所有未被删除的评论
func ListComment(videoId int) (*[]model.Comment, error) {

	//commentList, err := mysql.ListCommentDESCByCreateDate(videoId)
	//if err != nil {
	//	return nil, err
	//}

	key := tool.GetVideoCommentKey(videoId)
	commentList, err := redis.ListComment(key)
	if err != nil {
		return nil, err
	}

	return commentList, err
}
