package mysqlhook_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

func TestHook(t *testing.T) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/myapp_test?charset=utf8")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	var filter = func(entry *logrus.Entry) *logrus.Entry {
		if _, ok := entry.Data["foo2"]; ok {
			delete(entry.Data, "foo2")
		}

		return entry
	}

	tableName := "t_log"
	hook := mysqlhook.Default(db, tableName,
		mysqlhook.SetExtra(map[string]interface{}{"foo": "bar"}),
		mysqlhook.SetFilter(filter),
	)

	defer db.Exec(fmt.Sprintf("drop table `%s`", tableName))

	log := logrus.New()
	log.AddHook(hook)

	log.WithField("foo2", "bar").Infof("test foo")
	hook.Flush()

	row := db.QueryRow(fmt.Sprintf("select level,message,data,created from %s", tableName))

	var (
		level   int
		message string
		data    string
		created int64
	)

	err = row.Scan(&level, &message, &data, &created)
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

	if created == 0 || time.Unix(created, 0).IsZero() {
		t.Errorf("Not expected value:%v", created)
		return
	}
}

func ExampleHook() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/myapp_test?charset=utf8")
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
