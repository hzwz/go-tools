package main

import (
	"database/sql"
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	username string
	password string
	host     string
	database string
	coloum   string
	num      int64
	times    int64
	sql_exec string
	filename string
	from     int64
	to       int64
	step     int64
)

func Writer(data string, file *os.File) {
	_, err := file.WriteString(data)
	if err != nil {
		fmt.Println(err.Error)
	}

}

func main() {
	flag.StringVar(&username, "u", "", "DB username")
	flag.StringVar(&password, "p", "", "DB password")
	flag.StringVar(&host, "h", "127.0.0.1", "DB host")
	flag.StringVar(&database, "d", "", "Database Name")
	flag.StringVar(&coloum, "c", "", "DB coloum")
	flag.StringVar(&filename, "o", "", "Export file name")
	flag.Int64Var(&step, "s", 100, "The number of row  got each times")
	flag.Int64Var(&from, "f", 1, "Start id")
	flag.Int64Var(&to, "t", 1000, "End ID")

	flag.Parse()

	times = int64(math.Ceil((float64(to) - float64(from)) / float64(step)))
	if filename == "" {
		filename = coloum + ".csv"
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, host, database)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println("Failed to open file")
	}
	defer file.Close()

	for i := 0; i < int(times); i++ {
	loop0:

		sql_exec = fmt.Sprintf("SELECT * FROM %s where id>=%d and id<%d", coloum, from+int64(int64(i)*step), from+(int64(i+1)*step))
		rows, err := db.Query(sql_exec)
		if err != nil {
			//panic(err.Error())

			goto loop0
		}

		cols, _ := rows.Columns()
		vals := make([][]byte, len(cols))
		scans := make([]interface{}, len(cols))
		for k, _ := range vals {

			scans[k] = &vals[k]

		}
		for rows.Next() {
		loop1:
			row := make(map[string]string)
			err = rows.Scan(scans...)
			if err != nil {
				//panic(err.Error())
				time.Sleep(200 * time.Millisecond)
				goto loop1
			}

			line := ""

			for k, v := range vals {

				key := cols[k]
				//这里把[]byte数据转成string

				row[key] = string(v)
				line = line + "####" + string(v)

			}
			Writer(line+"\n", file)

		}
		rows.Close()
		time.Sleep(50 * time.Microsecond)

	}

}
