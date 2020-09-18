package application

import (
	"context"
	"encoding/json"
	"firstGo/application/start"
	"firstGo/model"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func prepareInit()  {
	// 解析配置类
	// 加载日志logger
	start.InitLogger(nil)
	// 加载redis
	start.InitRedis()
	// 加载mysql
	start.InitMysql([]*start.MysqlConfig{}...)
}

func Run()  {

	prepareInit()
	//r := router.InitRouter()
	//r.Run()

	/*
		test代码
	*/
	TestMysql()
}

func Close()  {
	// 释放资源
}

func TestMysql()  {
	db, _ := start.GetMysqlConn("db_user", 0)
	var order model.Order
	// order_id,city_id,order_status,latlong,reroute_time
	row := db.QueryRowx("select * from * where order_id=?",1)
	err := row.StructScan(&order)
	if err != nil {
		start.Logger.Error("scan err:" +  err.Error())
	} else {
		//start.Logger.Info(fmt.Sprintf("%#v", order))
	}
	if bytes,err := json.Marshal(order); err ==nil {
		fmt.Println(string(bytes))
	} else {
		fmt.Println(err)
	}
}

func TestRedis()  {
	client,_ := start.MRedisClient.GetRedisClient()
	ctx := context.Background()
	err := client.SetNX(ctx, "goTest", "hello go-redis", 100 * time.Second).Err()
	if err != nil {
		start.Logger.Info("set goTest err:" +  err.Error())
	}

	tm, err := client.Get(ctx, "goTest").Result()
	if err == nil {
		start.Logger.Info("goTest:" +  tm)
	}
	client.Set(ctx, "getTest1", "hello getTest1", 0).Err()
	if err != nil {
		start.Logger.Info("set goTest err:" +  err.Error())
	}

	tm1, err := client.Exists(ctx, []string{"getTest1"}...).Result()
	if err == nil {
		start.Logger.Info("goTest:" + strconv.Itoa(int(tm1)))
	}

	if tm1 > 0 {
		tm1, err := client.Append(ctx, "getTest1", " append").Result()
		if err == nil {
			start.Logger.Info("getTest1 append:" + strconv.Itoa(int(tm1)))
		}
	}

	client.Set(ctx, "mykey", "10",0)

	tm2, err := client.DecrBy(ctx, "mykey", 5).Result()
	if err == nil {
		start.Logger.Info("mykey decrby 5:" + strconv.Itoa(int(tm2)))
	}

	client.Set(ctx, "mykey", "This is a string",0)
	tm3,err := client.GetRange(ctx, "mykey", 0, 3).Result()
	if err == nil {
		start.Logger.Info("mykey getrange 0,3:" + tm3)
	}

	fmt.Println(client.StrLen(ctx, "mykey").Val())

	client.Del(ctx, "mykey")

	client.SAdd(ctx, "myset", "hello")
	client.SAdd(ctx, "myset", "world")
	client.SAdd(ctx, "myset", "world")
	tm4, err := client.SMembers(ctx, "myset").Result()

	if (err == nil) {
		for key, value := range tm4 {
			fmt.Println(key, value)
		}
	}

	tm5, err := client.SCard(ctx, "myset").Result()
	if (err == nil) {
		start.Logger.Info("myset SCard:" + strconv.Itoa(int(tm5)))
	}
	mebers := make([]*redis.Z, 0)
	mebers = append(mebers, &redis.Z{Score:1, Member:"one"})
	mebers = append(mebers, &redis.Z{Score:1, Member:"uno"})
	mebers = append(mebers, &redis.Z{Score:2, Member:"two"})
	mebers = append(mebers, &redis.Z{Score:3, Member:"three"})
	client.ZAdd(ctx, "myzset", mebers...)

	vals, err := client.ZRangeByScore(ctx, "myzset", &redis.ZRangeBy{
		Min: "(1",
		Max: "4",
		Offset: 0,
		Count: 100,
	}).Result()

	fmt.Printf("%#v\n", vals)

	// 管道 批量设置
	pipe := client.Pipeline()
	pipe.SetNX(ctx, "pipeKey1", 1, 0)
	pipe.SetNX(ctx, "pipeKey2", 2, 0)
	pipe.SetNX(ctx, "pipeKey3", 3, 0)
	pipe.SetNX(ctx, "pipeKey4", 4, 0)
	pipe.Exec(ctx)
}
