package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func Get(ctx context.Context, key string) string {
	var mes string
	var cmd *redis.StringCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		cmd = Client.Get(ctx, key)
	} else {
		cmd = ClusterClient.Get(ctx, key)
	}

	if err := cmd.Err(); err != nil {
		mes = ""
	} else {
		mes = cmd.Val()
	}
	return mes
}

func GetRaw(ctx context.Context, key string) (bts []byte, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		bts, err = Client.Get(ctx, key).Bytes()
	} else {
		bts, err = ClusterClient.Get(ctx, key).Bytes()
	}

	if err != nil && err != redis.Nil {
		return []byte{}, err
	}
	return bts, nil
}

func MGet(ctx context.Context, keys ...string) ([]string, error) {
	var sliceCmd *redis.SliceCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		sliceCmd = Client.MGet(ctx, keys...)
	} else {
		sliceCmd = ClusterClient.MGet(ctx, keys...)
	}

	if err := sliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	tmp := sliceCmd.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice, nil
}

func MGets(ctx context.Context, keys []string) (ret []interface{}, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		ret, err = Client.MGet(ctx, keys...).Result()
	} else {
		ret, err = ClusterClient.MGet(ctx, keys...).Result()
	}

	if err != nil && err != redis.Nil {
		return []interface{}{}, err
	}
	return ret, nil
}

func Set(ctx context.Context, key string, value interface{}, expire time.Duration) bool {
	var err error
	if options.Mode == "single" || options.Mode == "sentinel" {
		err = Client.Set(ctx, key, value, expire).Err()
	} else {
		err = ClusterClient.Set(ctx, key, value, expire).Err()
	}

	if err != nil {
		return false
	}
	return true
}

// HGetAll 从redis获取hash的所有键值对
func HGetAll(ctx context.Context, key string) map[string]string {
	var hash map[string]string
	var stringMapCmd *redis.StringStringMapCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		stringMapCmd = Client.HGetAll(ctx, key)
	} else {
		stringMapCmd = ClusterClient.HGetAll(ctx, key)
	}

	if err := stringMapCmd.Err(); err != nil && err != redis.Nil {
		hash = make(map[string]string)
	} else {
		hash = stringMapCmd.Val()
	}

	return hash
}

// HGet 从redis获取hash单个值
func HGet(ctx context.Context, key string, fields string) (string, error) {
	var stringCmd *redis.StringCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		stringCmd = Client.HGet(ctx, key, fields)
	} else {
		stringCmd = ClusterClient.HGet(ctx, key, fields)
	}

	err := stringCmd.Err()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		return "", nil
	}
	return stringCmd.Val(), nil
}

// HMGet 批量获取hash值
func HMGet(ctx context.Context, key string, fileds []string) []string {
	var sliceCmd *redis.SliceCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		sliceCmd = Client.HMGet(ctx, key, fileds...)
	} else {
		sliceCmd = ClusterClient.HMGet(ctx, key, fileds...)
	}

	if err := sliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	tmp := sliceCmd.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice
}

// HMGetMap 批量获取hash值，返回map
func HMGetMap(ctx context.Context, key string, fields []string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}

	var sliceCmd *redis.SliceCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		sliceCmd = Client.HMGet(ctx, key, fields...)
	} else {
		sliceCmd = ClusterClient.HMGet(ctx, key, fields...)
	}

	if err := sliceCmd.Err(); err != nil && err != redis.Nil {
		return make(map[string]string)
	}

	tmp := sliceCmd.Val()
	hashRet := make(map[string]string, len(tmp))

	var tmpTagID string

	for k, v := range tmp {
		tmpTagID = fields[k]
		if v != nil {
			hashRet[tmpTagID] = v.(string)
		} else {
			hashRet[tmpTagID] = ""
		}
	}
	return hashRet
}

// HMSet 设置redis的hash
func HMSet(ctx context.Context, key string, hash map[string]interface{}, expire time.Duration) bool {
	if len(hash) > 0 {
		var err error
		if options.Mode == "single" || options.Mode == "sentinel" {
			err = Client.HMSet(ctx, key, hash).Err()
		} else {
			err = ClusterClient.HMSet(ctx, key, hash).Err()
		}

		if err != nil {
			return false
		}
		Client.Expire(ctx, key, expire)
		return true
	}
	return false
}

// HSet hset
func HSet(ctx context.Context, key string, field string, value interface{}) bool {
	var err error
	if options.Mode == "single" || options.Mode == "sentinel" {
		err = Client.HSet(ctx, key, field, value).Err()
	} else {
		err = ClusterClient.HSet(ctx, key, field, value).Err()
	}

	if err != nil {
		return false
	}
	return true
}

// HDel ...
func HDel(ctx context.Context, key string, field ...string) bool {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.HDel(ctx, key, field...)
	} else {
		intCmd = ClusterClient.HDel(ctx, key, field...)
	}

	if err := intCmd.Err(); err != nil {
		return false
	}

	return true
}

// SetWithErr ...
func SetWithErr(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	if options.Mode == "single" || options.Mode == "sentinel" {
		return Client.Set(ctx, key, value, expire).Err()
	}

	return ClusterClient.Set(ctx, key, value, expire).Err()
}

// SetNx 设置redis的string 如果键已存在
func SetNx(ctx context.Context, key string, value interface{}, expiration time.Duration) bool {
	var res bool
	var err error
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.SetNX(ctx, key, value, expiration).Result()
	} else {
		res, err = ClusterClient.SetNX(ctx, key, value, expiration).Result()
	}

	if err != nil {
		return false
	}

	return res
}

// SetNxWithErr 设置redis的string 如果键已存在
func SetNxWithErr(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		return Client.SetNX(ctx, key, value, expiration).Result()
	}
	return ClusterClient.SetNX(ctx, key, value, expiration).Result()
}

// Incr redis自增
func Incr(ctx context.Context, key string) bool {
	var err error
	if options.Mode == "single" || options.Mode == "sentinel" {
		err = Client.Incr(ctx, key).Err()
	} else {
		err = ClusterClient.Incr(ctx, key).Err()
	}

	if err != nil {
		return false
	}
	return true
}

// IncrWithErr ...
func IncrWithErr(ctx context.Context, key string) (int64, error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		return Client.Incr(ctx, key).Result()
	}
	return ClusterClient.Incr(ctx, key).Result()
}

// IncrBy 将 key 所储存的值加上增量 increment 。
func IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.IncrBy(ctx, key, increment)
	} else {
		intCmd = ClusterClient.IncrBy(ctx, key, increment)
	}

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// Decr redis自减
func Decr(ctx context.Context, key string) bool {
	var err error
	if options.Mode == "single" || options.Mode == "sentinel" {
		err = Client.Decr(ctx, key).Err()
	} else {
		err = ClusterClient.Decr(ctx, key).Err()
	}

	if err != nil {
		return false
	}
	return true
}

// Type ...
func Type(ctx context.Context, key string) (string, error) {
	var statusCmd *redis.StatusCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		statusCmd = Client.Type(ctx, key)
	} else {
		statusCmd = ClusterClient.Type(ctx, key)
	}

	if err := statusCmd.Err(); err != nil {
		return "", err
	}
	return statusCmd.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
func ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	var stringSliceCmd *redis.StringSliceCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		stringSliceCmd = Client.ZRevRange(ctx, key, start, stop)
	} else {
		stringSliceCmd = ClusterClient.ZRevRange(ctx, key, start, stop)
	}

	if err := stringSliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// ZRevRangeWithScores ...
func ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	var zSliceCmd *redis.ZSliceCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		zSliceCmd = Client.ZRevRangeWithScores(ctx, key, start, stop)
	} else {
		zSliceCmd = ClusterClient.ZRevRangeWithScores(ctx, key, start, stop)
	}

	if err := zSliceCmd.Err(); err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}
	return zSliceCmd.Val(), nil
}

// ZRange ...
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	var stringSliceCmd *redis.StringSliceCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		stringSliceCmd = Client.ZRange(ctx, key, start, stop)
	} else {
		stringSliceCmd = ClusterClient.ZRange(ctx, key, start, stop)
	}

	if err := stringSliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// ZRevRank ...
func ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.ZRevRank(ctx, key, member)
	} else {
		intCmd = ClusterClient.ZRevRank(ctx, key, member)
	}

	if err := intCmd.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZRevRangeByScore ...
func ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) (res []string, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.ZRevRangeByScore(ctx, key, opt).Result()
	} else {
		res, err = ClusterClient.ZRevRangeByScore(ctx, key, opt).Result()
	}

	if err != nil && err != redis.Nil {
		return []string{}, err
	}

	return res, nil
}

// ZRevRangeByScoreWithScores ...
func ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) (res []redis.Z, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.ZRevRangeByScoreWithScores(ctx, key, opt).Result()
	} else {
		res, err = ClusterClient.ZRevRangeByScoreWithScores(ctx, key, opt).Result()
	}

	if err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}

	return res, nil
}

// ZCard 获取有序集合的基数
func nZCard(ctx context.Context, key string) (int64, error) {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.ZCard(ctx, key)
	} else {
		intCmd = ClusterClient.ZCard(ctx, key)
	}

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZScore 获取有序集合成员 member 的 score 值
func ZScore(ctx context.Context, key string, member string) (float64, error) {
	var floatCmd *redis.FloatCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		floatCmd = Client.ZScore(ctx, key, member)
	} else {
		floatCmd = ClusterClient.ZScore(ctx, key, member)
	}

	err := floatCmd.Err()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return floatCmd.Val(), err
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
func ZAdd(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.ZAdd(ctx, key, members...)
	} else {
		intCmd = ClusterClient.ZAdd(ctx, key, members...)
	}

	if err := intCmd.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func ZCount(ctx context.Context, key string, min, max string) (int64, error) {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.ZCount(ctx, key, min, max)
	} else {
		intCmd = ClusterClient.ZCount(ctx, key, min, max)
	}

	if err := intCmd.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return intCmd.Val(), nil
}

// Del redis删除
func Del(ctx context.Context, key string) int64 {
	var res int64
	var err error
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.Del(ctx, key).Result()
	} else {
		res, err = ClusterClient.Del(ctx, key).Result()
	}

	if err != nil {
		return 0
	}
	return res
}

// DelWithErr ...
func DelWithErr(ctx context.Context, key string) (int64, error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		return Client.Del(ctx, key).Result()
	}
	return ClusterClient.Del(ctx, key).Result()
}

// HIncrBy 哈希field自增
func HIncrBy(ctx context.Context, key string, field string, incr int) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		Client.HIncrBy(ctx, key, field, int64(incr))
	} else {
		ClusterClient.HIncrBy(ctx, key, field, int64(incr))
	}
}

// Exists 键是否存在
func Exists(ctx context.Context, key string) bool {
	var res int64
	var err error
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.Exists(ctx, key).Result()
	} else {
		res, err = ClusterClient.Exists(ctx, key).Result()
	}

	if err != nil {
		return false
	}
	return res == 1
}

// ExistsWithErr ...
func ExistsWithErr(ctx context.Context, key string) (bool, error) {
	var res int64
	var err error
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.Exists(ctx, key).Result()
	} else {
		res, err = ClusterClient.Exists(ctx, key).Result()
	}

	if err != nil {
		return false, nil
	}
	return res == 1, nil
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
func LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.LPush(ctx, key, values...)
	} else {
		intCmd = ClusterClient.LPush(ctx, key, values...)
	}

	if err := intCmd.Err(); err != nil {
		return 0, err
	}

	return intCmd.Val(), nil
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.RPush(ctx, key, values...)
	} else {
		intCmd = ClusterClient.RPush(ctx, key, values...)
	}

	if err := intCmd.Err(); err != nil {
		return 0, err
	}

	return intCmd.Val(), nil
}

// RPop 移除并返回列表 key 的尾元素。
func RPop(ctx context.Context, key string) (string, error) {
	var stringCmd *redis.StringCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		stringCmd = Client.RPop(ctx, key)
	} else {
		stringCmd = ClusterClient.RPop(ctx, key)
	}

	if err := stringCmd.Err(); err != nil {
		return "", err
	}

	return stringCmd.Val(), nil
}

// LRange 获取列表指定范围内的元素
func LRange(ctx context.Context, key string, start, stop int64) (res []string, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.LRange(ctx, key, start, stop).Result()
	} else {
		res, err = ClusterClient.LRange(ctx, key, start, stop).Result()
	}

	if err != nil {
		return []string{}, err
	}

	return res, nil
}

// LLen ...
func LLen(ctx context.Context, key string) int64 {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.LLen(ctx, key)
	} else {
		intCmd = ClusterClient.LLen(ctx, key)
	}

	if err := intCmd.Err(); err != nil {
		return 0
	}

	return intCmd.Val()
}

// LLenWithErr ...
func LLenWithErr(ctx context.Context, key string) (int64, error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		return Client.LLen(ctx, key).Result()
	}
	return ClusterClient.LLen(ctx, key).Result()
}

// LRem ...
func LRem(ctx context.Context, key string, count int64, value interface{}) int64 {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.LRem(ctx, key, count, value)
	} else {
		intCmd = ClusterClient.LRem(ctx, key, count, value)
	}

	if err := intCmd.Err(); err != nil {
		return 0
	}

	return intCmd.Val()
}

// LIndex ...
func LIndex(ctx context.Context, key string, idx int64) (string, error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		return Client.LIndex(ctx, key, idx).Result()
	}
	return ClusterClient.LIndex(ctx, key, idx).Result()
}

// LTrim ...
func LTrim(ctx context.Context, key string, start, stop int64) (string, error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		return Client.LTrim(ctx, key, start, stop).Result()
	}
	return ClusterClient.LTrim(ctx, key, start, stop).Result()
}

// ZRemRangeByRank 移除有序集合中给定的排名区间的所有成员
func ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (res int64, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.ZRemRangeByRank(ctx, key, start, stop).Result()
	} else {
		res, err = ClusterClient.ZRemRangeByRank(ctx, key, start, stop).Result()
	}

	if err != nil {
		return 0, err
	}

	return res, nil
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) (res bool, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.Expire(ctx, key, expiration).Result()
	} else {
		res, err = ClusterClient.Expire(ctx, key, expiration).Result()
	}

	if err != nil {
		return false, err
	}

	return res, err
}

// ZRem 从zset中移除变量
func ZRem(ctx context.Context, key string, members ...interface{}) (res int64, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.ZRem(ctx, key, members...).Result()
	} else {
		res, err = ClusterClient.ZRem(ctx, key, members...).Result()
	}

	if err != nil {
		return 0, err
	}
	return res, nil
}

// SAdd 向set中添加成员
func SAdd(ctx context.Context, key string, member ...interface{}) (int64, error) {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.SAdd(ctx, key, member...)
	} else {
		intCmd = ClusterClient.SAdd(ctx, key, member...)
	}

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// SMembers 返回set的全部成员
func SMembers(ctx context.Context, key string) ([]string, error) {
	var stringSliceCmd *redis.StringSliceCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		stringSliceCmd = Client.SMembers(ctx, key)
	} else {
		stringSliceCmd = ClusterClient.SMembers(ctx, key)
	}

	if err := stringSliceCmd.Err(); err != nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// SIsMember ...
func SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	var boolCmd *redis.BoolCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		boolCmd = Client.SIsMember(ctx, key, member)
	} else {
		boolCmd = ClusterClient.SIsMember(ctx, key, member)
	}

	if err := boolCmd.Err(); err != nil {
		return false, err
	}
	return boolCmd.Val(), nil
}

// HKeys 获取hash的所有域
func HKeys(ctx context.Context, key string) []string {
	var stringSliceCmd *redis.StringSliceCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		stringSliceCmd = Client.HKeys(ctx, key)
	} else {
		stringSliceCmd = ClusterClient.HKeys(ctx, key)
	}

	if err := stringSliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	return stringSliceCmd.Val()
}

// HLen 获取hash的长度
func HLen(ctx context.Context, key string) int64 {
	var intCmd *redis.IntCmd
	if options.Mode == "single" || options.Mode == "sentinel" {
		intCmd = Client.HLen(ctx, key)
	} else {
		intCmd = ClusterClient.HLen(ctx, key)
	}

	if err := intCmd.Err(); err != nil && err != redis.Nil {
		return 0
	}
	return intCmd.Val()
}

// GeoAdd 写入地理位置
func GeoAdd(ctx context.Context, key string, location *redis.GeoLocation) (res int64, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.GeoAdd(ctx, key, location).Result()
	} else {
		res, err = ClusterClient.GeoAdd(ctx, key, location).Result()
	}

	if err != nil {
		return 0, err
	}

	return res, nil
}

// GeoRadius 根据经纬度查询列表
func GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (res []redis.GeoLocation, err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		res, err = Client.GeoRadius(ctx, key, longitude, latitude, query).Result()
	} else {
		res, err = ClusterClient.GeoRadius(ctx, key, longitude, latitude, query).Result()
	}
	
	if err != nil {
		return []redis.GeoLocation{}, err
	}

	return res, nil
}

// Close closes the client, releasing any open resources.
//
// It is rare to Close a Client, as the Client is meant to be
// long-lived and shared between many goroutines.
func Close() (err error) {
	if options.Mode == "single" || options.Mode == "sentinel" {
		if Client != nil {
			err = Client.Close()
		}
	} else {
		if ClusterClient != nil {
			err = ClusterClient.Close()
		}
	}
	return
}