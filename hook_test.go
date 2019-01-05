package mysqlhook_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dsn = "root:@tcp(127.0.0.1:3306)/myapp_test?charset=utf8&parseTime=true"
)

func TestHook(t *testing.T) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	tableName := "t_log"
	extraItems := []*mysqlhook.ExecExtraItem{
		mysqlhook.NewExecExtraItem("type", "varchar(50)"),
	}
	hook := mysqlhook.DefaultWithExtra(db, tableName,
		extraItems,
		mysqlhook.SetExtra(map[string]interface{}{"foo": "bar"}),
		mysqlhook.SetFilter(func(entry *logrus.Entry) *logrus.Entry {
			if _, ok := entry.Data["foo2"]; ok {
				delete(entry.Data, "foo2")
			}
			return entry
		}),
	)

	defer db.Exec(fmt.Sprintf("drop table `%s`", tableName))

	log := logrus.New()
	log.AddHook(hook)

	log.WithField("foo2", "bar").WithField("type", "test").Infof("test foo")
	hook.Flush()

	row := db.QueryRow(fmt.Sprintf("select level,message,data,time,type from %s", tableName))

	var (
		level   int
		message string
		data    string
		tt      time.Time
		typ     string
	)

	err = row.Scan(&level, &message, &data, &tt, &typ)
	if err != nil {
		t.Error(err)
		return
	}

	if logrus.Level(level) != logrus.InfoLevel {
		t.Errorf("Not expected value:%v", level)
		return
	}

	if message != "test foo" {
		t.Errorf("Not expected value:%v", message)
		return
	}

	if typ != "test" {
		t.Errorf("Not expected value:%v", typ)
		return
	}

	var m map[string]string
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		t.Error(err)
		return
	}

	if m == nil || m["foo2"] != "" || m["foo"] != "bar" {
		t.Errorf("Not expected value:%v", m)
		return
	}

	if tt.IsZero() {
		t.Errorf("Not expected value:%v", tt)
		return
	}
}

func ExampleHook() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	tableName := "e_log"
	mysqlHook := mysqlhook.Default(db, tableName)
	defer db.Exec(fmt.Sprintf("drop table %s", tableName))

	log := logrus.New()
	log.AddHook(mysqlHook)
	log.WithField("foo", "bar").Info("foo test")
	mysqlHook.Flush()

	var message string
	row := db.QueryRow(fmt.Sprintf("select message from %s", tableName))
	err = row.Scan(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(message)

	// Output: foo test
}
