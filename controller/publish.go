package controller

import (
	"BytesDanceProject/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//token := c.Query("token")

	//用户鉴权
	//if _, exist := usersLoginInfo[token]; !exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//	return
	//}

	//获取文件
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//上传文件到七牛云空间
	err = service.UploadVideo(file)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//此处注释掉的为官方demo
	//filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]
	//finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	//saveFile := filepath.Join("./public/", finalName)
	//if err := c.SaveUploadedFile(data, saveFile); err != nil {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		//StatusMsg:  finalName + " uploaded successfully",
		StatusMsg: "uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	//token := c.Query("token")

	//用户鉴权
	//if _, exist := usersLoginInfo[token]; !exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//	return
	//}

	//获取当前登录用户发布的所有视频
	originalVideoList, err := service.ListVideosByUser() //【！！！！！此处应该传入当前登录用户的对象，因为还没有创建user对象，故不进行此操作】
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//获取到的originalVideoList（model.Video）需要进行处理，使其变成满足前端接口的要求的videoList（controller.Video）
	var videoList = make([]Video, len(*originalVideoList))
	point := 0 //videoList的指针
	for _, originalVideo := range *originalVideoList {

		//根据authorId获取author对象
		//authorId := originalVideo.AuthorId

		var favoriteCount int64 = 666 //！！！！假数据
		//查询当前视频的点赞数

		var commentCount int64 = 777 //！！！！假数据
		//查询当前视频的评论数

		isFavorite := false //！！！！！！假数据
		//查询当前登录用户是否喜欢该视频。如果当前用户没有登录，则为false

		video := Video{ //注意video中omitempty！！！
			Id:            int64(originalVideo.Id), //若为0则生成json时不包含该字段
			Author:        User{},                  //待处理
			PlayUrl:       originalVideo.PlayUrl,   //若为空则生成json时不包含该字段
			CoverUrl:      originalVideo.CoverUrl,  //若为空则生成json时不包含该字段
			FavoriteCount: favoriteCount,           //若为0则生成json时不包含该字段
			CommentCount:  commentCount,            //若为0则生成json时不包含该字段
			IsFavorite:    isFavorite,              ////若为false则生成json时不包含该字段
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