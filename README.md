# MySQL Hook for [Logrus](https://github.com/sirupsen/logrus)

> A mysql-based logrus hook

[![Build][Build-Status-Image]][Build-Status-Url] [![Coverage][Coverage-Image]][Coverage-Url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Quick Start

### Download and install

```bash
$ go get -u -v github.com/LyricTian/logrus-mysql-hook
```

### Usage

```go
import "github.com/LyricTian/logrus-mysql-hook"

// ...

mysqlHook := mysqlhook.New(
	mysqlhook.SetExec(mysqlhook.NewExec(db, "log")),
)

defer mysqlHook.Flush()

log := logrus.New()
log.AddHook(mysqlHook)
```

### Examples

```go
package main

import (
	"database/sql"
	"fmt"

	"github.com/LyricTian/logrus-mysql-hook"
	"github.com/Sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	tableName := "e_log"
	mysqlHook := mysqlhook.New(
		mysqlhook.SetExec(mysqlhook.NewExec(db, tableName)),
	)
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
```

## MIT License

    Copyright (c) 2018 Lyric

[Build-Status-Url]: https://travis-ci.org/LyricTian/logrus-mysql-hook
[Build-Status-Image]: https://travis-ci.org/LyricTian/logrus-mysql-hook.svg?branch=master
[Coverage-Url]: https://coveralls.io/github/LyricTian/logrus-mysql-hook?branch=master
[Coverage-Image]: https://coveralls.io/repos/github/LyricTian/logrus-mysql-hook/badge.svg?branch=master
[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/logrus-mysql-hook
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/logrus-mysql-hook
[godoc-url]: https://godoc.org/github.com/LyricTian/logrus-mysql-hook
[godoc-image]: https://godoc.org/github.com/LyricTian/logrus-mysql-hook?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg