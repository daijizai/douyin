package controller

import (
	"BytesDanceProject/pkg/jwt"
	"BytesDanceProject/service"
	"BytesDanceProject/tool"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish 视频发布
func Publish(c *gin.Context) {
	//用户鉴权
	token := c.PostForm("token")

	claim, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	} else if claim.Valid() != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: claim.Valid().Error()})
		return
	}

	//获取标题
	title := c.PostForm("title")

	//获取文件
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		fmt.Println(err.Error())
		return
	}

	//上传文件到七牛云空间
	err = service.UploadVideo(file, title, claim.UserId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		fmt.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		//StatusMsg:  finalName + " uploaded successfully",
		StatusMsg: "uploaded successfully",
	})
}

// PublishList 获取发布列表
func PublishList(c *gin.Context) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		return
	}
	activeUserId := userIdInterface.(int)

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		fmt.Println(err.Error())
		return
	}

	//获取用户发布的所有视频
	originalVideoList, err := service.ListVideosByUser(int(userId))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		fmt.Println(err.Error())
		return
	}

	//获取到的originalVideoList（model.Video）需要进行处理，使其变成满足前端接口的要求的videoList（controller.Video）
	var videoList = make([]Video, len(*originalVideoList))
	point := 0 //videoList的指针
	for _, originalVideo := range *originalVideoList {

		user, err := service.GetUser(originalVideo.AuthorId) //获取视频的作者
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println(err.Error())
			return
		}

		followerCount, err := service.CountFollower(int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println("获取点赞列表失败" + err.Error())
			return
		}

		followCount, err := service.CountFollowee(int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println("获取点赞列表失败" + err.Error())
			return
		}

		isFollow, err := service.CheckFollowStatus(activeUserId, int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println("获取点赞列表失败" + err.Error())
			return
		}

		author := User{
			Id:            user.Id,
			Name:          user.UserName,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      isFollow,
		}

		likeCount, err := service.CountLike(originalVideo.Id) //获取视频的喜欢数
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println("获取发布列表失败" + err.Error())
			return
		}

		commentCount, err := service.CountCommentByVideoId(originalVideo.Id) //获取视频的评论数
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println(err.Error())
			return
		}

		likeStatus, err := service.GetLikeStatus(originalVideo.Id, activeUserId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			fmt.Println(err.Error())
			return
		}

		playUrl := "http://" + viper.GetString("qiniuyun.domain") + "/" + originalVideo.PlayUrl
		coverUrl := "http://" + viper.GetString("qiniuyun.domain") + "/" + originalVideo.CoverUrl

		video := Video{ //注意video中omitempty！！！
			Id:            int64(originalVideo.Id),          //若为0则生成json时不包含该字段
			Author:        author,                           //待处理
			PlayUrl:       playUrl,                          //若为空则生成json时不包含该字段
			CoverUrl:      coverUrl,                         //若为空则生成json时不包含该字段
			FavoriteCount: likeCount,                        //若为0则生成json时不包含该字段
			CommentCount:  commentCount,                     //若为0则生成json时不包含该字段
			IsFavorite:    likeStatus,                       ////若为false则生成json时不包含该字段
			Title:         tool.Filter(originalVideo.Title), //若为空则生成json时不包含该字段
		}
		videoList[point] = video
		point++
	}

	//返回响应
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功获取当前登录用户所有投稿过的视频",
		},
		VideoList: videoList,
	})
}
