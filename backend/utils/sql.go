package utils

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Query(queryStatement string) {
	fmt.Println("1")
	// MySQL 데이터베이스 정보
	// db, err := sql.Open("mysql", "root:myyonseinft@tcp(52.204.136.223:3306)/myyonseinft")
	db, err := sql.Open("mysql", "root:myyonseinft@tcp(localhost:3306)/myyonseinft")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("2")

	// 쿼리 실행
	rows, err := db.Query(queryStatement)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	fmt.Println("3")

	// 결과 반복문으로 처리
	for rows.Next() {
		var address string
		var department string
		var tokenuri string
		var tx string
		var time time.Time
		if err := rows.Scan(&address, &department, &tokenuri, &tx, &time); err != nil {
			panic(err.Error())
		}
		fmt.Println(address, department, tokenuri, tx, time)
	}
	// 오류 체크
	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
}
