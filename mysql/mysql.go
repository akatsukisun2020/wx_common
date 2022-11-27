package mysql

import (
	"database/sql"
	"time"

	"github.com/akatsukisun2020/wx_common/logger"

	_ "github.com/go-sql-driver/mysql"
)

type MYSQLCli struct {
	dsn          string
	maxIdleConns int
	maxOpenConns int
	maxLifetime  time.Duration
}

func NewMYSQLCli(dsn string, maxIdleConns, maxOpenConns int, maxLifetime time.Duration) *MYSQLCli {
	return &MYSQLCli{
		dsn:          dsn,
		maxIdleConns: maxIdleConns,
		maxOpenConns: maxOpenConns,
		maxLifetime:  maxLifetime,
	}
}

// Open 打开数据库:  注意，这里实际上就是打开的是一个dns对应的db的链接池子 ==> 不同的dns对应的池子，使用方可以自己使用map等管理
func (cli *MYSQLCli) OpenDB() (*sql.DB, error) {
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, err := sql.Open("mysql", cli.dsn)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(cli.maxLifetime)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(cli.maxIdleConns)
	// 最大连接数量
	DB.SetMaxOpenConns(cli.maxOpenConns)

	//验证连接
	if err = DB.Ping(); err != nil {
		logger.Errorf("OpenDB fail, err:%v", err)
		return nil, err
	}
	logger.Debugf("dns:%s, OpenDB success", cli.dsn)
	return DB, nil
}

func (cli *MYSQLCli) CloseDB(db *sql.DB) error { // 如果没有必要，这个一般在整个应用进程退出的时候，自动退出就好了，不用显示调用
	return db.Close()
}
