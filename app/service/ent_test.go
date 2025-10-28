package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/product"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productmanager"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/user"
	"context"
	"entgo.io/ent/dialect"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL驱动
	"testing"
)

func TestSelect(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3307)/activate_server?charset=utf8mb4&parseTime=True&loc=Local"
	client, err := ent.Open(dialect.MySQL, dsn)
	if err != nil {
		fmt.Println(err)
	}
	products, err := client.Product.Query().
		WithManagers(func(pmq *ent.ProductManagerQuery) {
			pmq.Select(
				productmanager.FieldProductID,
				productmanager.FieldUserID,
				productmanager.FieldPermissions,
				productmanager.FieldRole,
			).WithUser(func(uq *ent.UserQuery) {
				uq.Select(
					user.FieldID,
					user.FieldEmail,
				)
			})
		}).
		Select(
			product.FieldID,
			product.FieldCode,
			product.FieldProductType,
			product.FieldProductName,
			product.FieldCreatedAt,
			product.FieldUpdatedAt,
		).
		All(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v", products)
}
