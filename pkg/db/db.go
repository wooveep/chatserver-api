/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 14:38:52
 * @LastEditTime: 2023-04-05 15:54:17
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/db/db.go
 */
package db

import (
	"chatserver-api/pkg/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ IDataSource = (*defaultPostGresDataSource)(nil)

// IDataSource 定义数据库数据源接口，按照业务需求可以返回主库链接Master和从库链接Slave
type IDataSource interface {
	Master() *gorm.DB
	Close()
}

type defaultPostGresDataSource struct {
	master *gorm.DB // 定义私有属性，用来持有主库链接，防止每次创建，创建后直接返回该变量。
}

func (d *defaultPostGresDataSource) Master() *gorm.DB {
	if d.master == nil {
		panic("The [master] connection is nil, Please initialize it first.")
	}
	return d.master
}

func (d *defaultPostGresDataSource) Close() {
	// 关闭主库链接
	if d.master != nil {
		m, err := d.master.DB()
		if err != nil {
			_ = m.Close()
		}
	}
}

func NewDefaultPostGre(c config.DBConfig) *defaultPostGresDataSource {
	return &defaultPostGresDataSource{
		master: connect(
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.Dbname,
			c.MaximumPoolSize,
			c.MaximumIdleSize),
	}
}

func connect(user, password, host, port, dbname string, maxPoolSize, maxIdle int) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s  dbname=%s  port=%s  sslmode=disable TimeZone=Asia/Shanghai",
		host,
		user,
		password,
		dbname,
		port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // 缓存每一条sql语句，提高执行速度
	})
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetConnMaxLifetime(time.Hour)
	// 设置连接池大小
	sqlDb.SetMaxOpenConns(maxPoolSize)
	sqlDb.SetMaxIdleConns(maxIdle)
	return db
}
