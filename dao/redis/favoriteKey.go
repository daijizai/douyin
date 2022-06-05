package redis

const (
	KeyPrefix               = "douyin"
	KeyVideoFavoritedZSetPf = "video:favorited:" // zset : 记录用户及点赞操作类型；参数是 video_id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
