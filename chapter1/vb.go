package chapter1

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type WeiboInterface interface {
	// WriteWeibo 发布微博
	WriteWeibo(string, string, string) string
	// WeiboDetail 查看微博详情
	WeiboDetail(string) map[string]string
	// Feeds 微博的feed流
	Feeds(int64) []map[string]string
	// CommentsWeibo 给微博写评论
	CommentsWeibo()
	// CommentsList 某条微博的评论列表
	CommentsList(string) []map[string]string
	// ThumbWeibo 给微博点赞
	ThumbWeibo(string, string)
	// 取消点赞
	UnThumbWeibo()
	// ThumbList 某条微博的点赞列表
	ThumbList(string) ([]string, error)
	// 获取某条微博的总点赞数
	ThumbCount(string) int
}

type WeiboStruct struct {
	Conn *redis.Client
}

var (
	ctx                = context.Background()
	OneWeek            = 7 * 24 * 60 * 60 * time.Second
	WeiboPerPage int64 = 5 // 1页5条数据

	WeiboCountKey  = "weibo-count"     // 总微博数str
	WeiboDetailKey = "weibo-detail:"   // 微博详情hash
	SendTimeKey    = "weibo-send-time" // 微博发送时间zset
	ThumbsKey      = "thumb-weibo"     // 某条微博的总点赞zset，每条微博一条记录
	ThumbUsersKey  = "thumb-users:"    // 某条微博的点赞用户set
)

func NewWeibo(conn *redis.Client) *WeiboStruct {
	return &WeiboStruct{Conn: conn}
}

func (this WeiboStruct) WriteWeibo(userId, title, link string) string {
	weiboId := strconv.Itoa(int(this.Conn.Incr(ctx, WeiboCountKey).Val()))
	weiboThumbsKey := fmt.Sprintf("%s%s", ThumbUsersKey, weiboId)
	weiboDetailKey := fmt.Sprintf("%s%s", WeiboDetailKey, weiboId)

	// 七天之后无法投票
	// 默认发送微博即点赞
	this.Conn.SAdd(ctx, weiboThumbsKey, userId)
	this.Conn.Expire(ctx, weiboThumbsKey, OneWeek*time.Second)

	this.Conn.HSet(ctx, weiboDetailKey, map[string]interface{}{
		"title":  title,
		"link":   link,
		"poster": userId,
		"time":   time.Now().Unix(),
		"thumbs": 1,
	})

	// 其实不用加这条记录，查的时候如果查不到默认给0就可以了，95%的微博是没有人点赞的，没必要浪费redis空间
	this.Conn.ZAdd(ctx, ThumbsKey, &redis.Z{Member: weiboDetailKey, Score: 1})
	this.Conn.ZAdd(ctx, SendTimeKey, &redis.Z{Member: weiboDetailKey, Score: float64(time.Now().Unix())})

	return weiboId
}

func (this WeiboStruct) WeiboDetail(weiboId string) map[string]string {

	detail := this.Conn.HGetAll(ctx, WeiboDetailKey+weiboId).Val()
	//var vb Weibo
	//err := mapstructure.Decode(detail, &vb)
	//if err != nil {
	//	return Weibo{}
	//}
	//return vb
	return detail
}

func (this WeiboStruct) Feeds(page int64) []map[string]string {
	start := (page - 1) * WeiboPerPage
	end := start + WeiboPerPage - 1

	// 点赞数从多到少
	ids := this.Conn.ZRevRange(ctx, ThumbsKey, start, end).Val()
	vbs := []map[string]string{}
	for _, id := range ids {
		vbData := this.Conn.HGetAll(ctx, id).Val()
		vbData["id"] = id
		vbs = append(vbs, vbData)
	}

	return vbs
}

func (this WeiboStruct) CommentsWeibo() {

}

func (this WeiboStruct) CommentsList() {

}

// todo 某个用户发布的所有微博

func (this WeiboStruct) ThumbWeibo(weiboId, userId string) {
	// 只有有人点赞的微博才创建对应的key，没有就不创建，节省空间
	weiboThumbsKey := fmt.Sprintf("%s%s", ThumbUsersKey, weiboId)
	weiboMember := fmt.Sprintf("weibo:%s", weiboId)
	weiboThumbsMember := fmt.Sprintf("user:%s", userId)

	// 超过一周的微博不能投票

	// 查看thumbs里有没有该member，没有则添加数据，有则+1
	rank := this.Conn.ZRank(ctx, ThumbsKey, weiboMember)
	if rank != nil {
		this.Conn.ZIncrBy(ctx, ThumbsKey, 1, weiboMember)
	} else {
		this.Conn.ZAdd(ctx, ThumbsKey, &redis.Z{Member: weiboMember, Score: 1})
	}

	// 一个用户对一条微博只能投一票
	members := this.Conn.SIsMember(ctx, weiboThumbsKey, weiboThumbsMember)
	if members != nil {
		this.Conn.SAdd(ctx, weiboThumbsKey, weiboThumbsMember)
	} else {
		fmt.Printf("已投票，无法再投: %s", userId)
	}
}

func (this WeiboStruct) ThumbList(weiboId string) ([]string, error) {
	weiboThumbsKey := fmt.Sprintf("%s%s", ThumbUsersKey, weiboId)

	thumbList := this.Conn.SMembers(ctx, weiboThumbsKey)

	if thumbList.Err() != nil {
		panic(thumbList.Err())
		return nil, thumbList.Err()
	}

	// todo 从redis里的users:xxx直接取出用户信息

	return thumbList.Val(), nil
}
