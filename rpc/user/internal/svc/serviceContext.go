package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"tikstart/common/model"
	"tikstart/rpc/user/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	RDS    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(getDSN(&c)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.MySQL.TablePrefix, // 表明前缀，可不设置
			SingularTable: true,                // 使用单数表名，即不会在表名后添加复数s
		},
		//Logger: logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
	})

	rds := redis.MustNewRedis(c.Redis.RedisConf)

	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.User{}, &model.Follow{}, &model.Friend{})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		RDS:    rds,
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
