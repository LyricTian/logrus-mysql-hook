package mysqlhook

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
)

// Execer write the logrus entry to the database
type Execer interface {
	Exec(entry *logrus.Entry) error
}

// NewExec create an exec instance
func NewExec(db *sql.DB, tableName string) Execer {
	query := fmt.Sprintf("create table if not exists `%s` (`id` bigint not null primary key auto_increment, `level` int, `message` text, `data` text, `created` bigint)  engine=MyISAM charset=UTF8;", tableName)
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}

	return &defaultExec{db, tableName}
}

type defaultExec struct {
	db        *sql.DB
	tableName string
}

func (e *defaultExec) Exec(entry *logrus.Entry) error {
	jsonData, err := json.Marshal(entry.Data)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("insert into `%s` (`id`,`level`,`message`,`data`,`created`) values (null,?,?,?,?);", e.tableName)
	_, err = e.db.Exec(query, entry.Level, entry.Message, string(jsonData), entry.Time.Unix())
	if err != nil {
		return err
	}

	return nil
}
