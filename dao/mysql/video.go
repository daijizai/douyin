package mysql

import (
	"BytesDanceProject/model"
	"fmt"
)

/**
 * @author  daijizai Congregalis
 * @date  2022/5/10 20:22
 * @version  1.0
 * @description
 */

// ListVideo 供feed流使用 获取视频列表  【！！！！应该限制时间】
func ListVideo(videoCount int) (*[]model.Video, error) {
	//sqlStr := `select id,author_id,play_url,cover_url,favorite_count,comment_count from video limit ?`
	sqlStr := `select id,author_id,play_url,cover_url from video limit ?`
	var videoList []model.Video
	err := db.Select(&videoList, sqlStr, videoCount)
	if err != nil {
		return nil, err
	}
	return &videoList, nil
}

// InsertVideo 插入一条video记录 id为主键自增
func InsertVideo(v model.Video) error {
	//sqlStr := `INSERT INTO video(author_id, play_url,cover_url,favorite_count,comment_count,is_deleted,create_time)
	//VALUES(?,?,?,0,0,0,?)`
	sqlStr := `INSERT INTO video(author_id, play_url,cover_url,is_deleted,create_time) 
	VALUES(?,?,?,0,?)`
	ret, err := db.Exec(sqlStr, v.AuthorId, v.PlayUrl, v.CoverUrl, v.CreateTime)

	if err != nil {
		return err
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		return err
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
	return nil
}

// ListVideoByAuthorId 根据作者id获取视频列表
func ListVideoByAuthorId(authorId int) (*[]model.Video, error) {
	//sqlStr := `select id,author_id,play_url,cover_url,favorite_count,comment_count from video where author_id =?`
	sqlStr := `select id,author_id,play_url,cover_url from video where author_id =?`
	var videoList []model.Video
	err := db.Select(&videoList, sqlStr, authorId)
	if err != nil {
		return nil, err
	}
	return &videoList, nil
}
