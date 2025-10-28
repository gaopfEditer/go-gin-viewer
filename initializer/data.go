package initializer

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/user"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"

	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"entgo.io/ent/dialect"

	"cambridge-hit.com/gin-base/activateserver/pkg/util/cache"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/objstorage"
	_ "github.com/go-sql-driver/mysql" // MySQL驱动
)

// freeCacheInit FreeCache初始化
func cacheInit() {
	//cache.MyFreeCache = cache.NewFreeCache()
	//cache.MyFileCache = cache.NewFileCache("./cache")
	cache.MyRedis = cache.NewRedisCache()
}

func ossInit(product string) {
	var err error
	objstorage.InitS3Clients(
		product,
		resource.Conf.ObjStorageConfig.Region, //服务器上传下载可用内网
		resource.Conf.ObjStorageConfig.BucketName,
		resource.Conf.ObjStorageConfig.AccessID,
		resource.Conf.ObjStorageConfig.AccessSecret,
	)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func createDefaultAdminUser(client *ent.Client) {
	// 初始化匿名用户(记录日志时的默认值，禁止登录)
	createUser(client, dto.AnonymousID, "anonymous@anonymous.com", "123456", false)
	// 初始化超级管理员用户
	createUser(client, dto.SuperAdminID, "system@cambridge-hit.com", "123456", true)
}

func createUser(client *ent.Client, id int, email, password string, is_enabled bool) {
	exists, err := client.User.Query().Where(user.ID(id)).Exist(context.Background())
	if err != nil {
		log.Fatalf("检查用户是否存在失败: %v", err)
		return
	}

	// 如果不存在 id 为 指定id 的用户，则创建该用户
	if !exists {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("密码加密失败: %v", err)
			return
		}

		_, err = client.User.Create().
			SetID(id).
			SetEmail(email).
			SetPassword(string(hashedPassword)).
			SetIsEnabled(is_enabled).
			Save(context.Background())

		if err != nil {
			log.Fatalf("创建默认用户失败: %v", err)
			return
		}

		log.Println("已创建默认用户,id:", id)
	}

}

func dbInit() {
	cfg := resource.Conf.MysqlConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	client, err := ent.Open(dialect.MySQL, dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
		return
	}

	// 运行数据库迁移（自动创建表）
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
		return
	}

	dto.SetClient(client)

	// 创建默认管理员用户
	createDefaultAdminUser(client)
}
