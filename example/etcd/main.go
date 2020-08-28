package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"sapi/pkg/client/etcdv3"
)

type User struct {
	gorm.Model
	Name string
	Age  sql.NullInt64  // sql.NullInt64 实现了 Scanner/Valuer 接口
}

func main() {
	//ctx := context.Background()

	//logger.Init(logger.NewOptions())
	//
	//tracer.Init(tracer.NewOptions())
	//
	//redis.Init(ctx, redis.NewOptions(redis.Password("mylove1688")))
	//
	//option := model.NewOptions(
	//	model.Database("test"),
	//	model.Username("root"),
	//	model.Password(""),
	//	model.Prefix("t_"),
	//	)
	//err := model.Init(option)
	//
	//tr := global.Tracer("component-main")
	//ctx, span := tr.Start(ctx, "main")
	//span.SetAttributes(
	//	kv.String("test1", "redis"),
	//	kv.String("test2", "aa"),
	//)
	//span.End()
	//
	//tracer.AddHookSpanCtx(ctx)
	//
	//if err == nil {
	//	var user User
	//	user = User{Name: "Jinzhu", Age: sql.NullInt64{Int64:18,Valid:true}}
	//	model.DB.Create(&user)
	//	model.DB.Where("name = ?", "Jinzhu").First(&user)
	//	model.DB.Model(&user).Update("name", "hello")
	//	model.DB.Where("id = ?", 20).Unscoped().Delete(&user)
	//}
	//
	//redis.Set(ctx, "key", "asdasdad", 0)
	//
	//res := redis.Get(ctx, "key")
	//
	//logger.Log.Info("[redis]", zap.Any("key", res))
	//
	//logger.Log.Info("[initCfg] 配置", zap.Any("cfg", conf.Conf))
	//
	//watch, err := file.Watch("config/config.toml", &conf.Conf)
	//
	//go func() {
	//	for {
	//		_, err = watch.Next()
	//		if err == nil {
	//			logger.Log.Info("[initCfg] 配置", zap.Any("cfg", conf.Conf))
	//		}
	//	}
	//}()

	//fh, err := os.Open("config/config.toml")
	//if err != nil {
	//	return
	//}
	//defer fh.Close()
	//b, err := ioutil.ReadAll(fh)
	//if err != nil {
	//	return
	//}
	//
	//etcdv3 := etcdv3.NewEtcd()
	//etcdv3.Put(string(b))
	//etcdv3.Read(&conf.Conf)
	//
	//logger.Log.Info("[initCfg] 配置", zap.Any("cfg", conf.Conf))

	//etcdv3 := etcdv3.NewEtcd()
	//etcdv3.Read(&conf.Conf)
	//
	//var data interface{}
	//watch, err := etcdv3.Watch()
	//for {
	//	data, err = watch.Next()
	//	if err == nil {
	//		logger.Log.Info("[initCfg] 配置", zap.Any("cfg", data))
	//	}
	//}

	//bootstrap.ConfigFlag(
	//	&flag.StringFlag{
	//		Name:    "config_etcd_format",
	//		Usage:   "--config_etcd_format",
	//		Default: "yml",
	//	},
	//)
	//
	//eng := bootstrap.NewEngine()
	//err := eng.Startup(
	//	bootstrap.InitFlag,
	//	bootstrap.InitConfig,
	//	bootstrap.InitLog,
	//	bootstrap.InitRedis,
	//	bootstrap.InitModel,
	//	bootstrap.InitTracer,
	//	bootstrap.InitGin,
	//	)
	//
	//if err != nil {
	//	log.Panic(err)
	//	return
	//}
	//
	//rsy, err := etcd.NewRegistry(&registry.Options{
	//	Timeout: 3,
	//	TTL: 5,
	//})
	//
	//rsy.Register(bootstrap.GinServer.GetOption())
	//watcher, _ := rsy.Watch()
	//go func() {
	//	for {
	//		data, err := watcher.Next()
	//		if err == nil {
	//			fmt.Println(data.Action)
	//			fmt.Println(data.Service)
	//		}
	//	}
	//}()
	//
	//rsy.ListServices()
	////rsy.Close()
	//
	//eng.AddServe(bootstrap.GinServer)
	//eng.Serve()

	cli := etcdv3.NewOptions().Build()
	reg := cli.NewRegistry()
	err := reg.Register("/test/test", "asdasd", 5000)
	err = reg.Deregister("/test/test")
	fmt.Println(err)

	b, err := cli.GetValue(context.Background(), "/test/test")
	fmt.Println(err)
	fmt.Println(string(b))

}
