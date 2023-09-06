package svc

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"tikstart/rpc/contact/contactClient"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/videoClient"
	"tikstart/task/server/internal/config"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	RDS        *redis.Client
	UserRpc    userClient.User
	VideoRpc   videoClient.Video
	ContactRpc contactClient.Contact
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(getDSN(&c)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.MySQL.TablePrefix, // 表明前缀，可不设置
			SingularTable: true,                // 使用单数表名，即不会在表名后添加复数s
		},
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("cannot connect to mysql: %v", err)
	}

	rds := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Pass,
		DB:       c.Redis.DB,
	})
	_, err = rds.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("cannot connect to redis: %v", err)
	}

	return &ServiceContext{
		Config:     c,
		DB:         db,
		RDS:        rds,
		UserRpc:    userClient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc:   videoClient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		ContactRpc: contactClient.NewContact(zrpc.MustNewClient(c.ContactRpc)),
	}
}

func getDSN(c *config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQL.User,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.Database,
	)
}
