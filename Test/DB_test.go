package test

import (
	"database/sql"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	connectionString = "server=JIN_SURFACE\\SQLEXPRESS;database=BaseballData;integrated security=true"
)

func openDB() (*sql.DB, error) {
	// MSSQL에 연결하기 위한 DSN(Data Source Name)
	dsn := connectionString
	return sql.Open("mssql", dsn)
}

func TestDBConnection(t *testing.T) {
	// 데이터베이스 연결 시도
	db, err := openDB()
	if err != nil {
		t.Fatalf("DB 연결에 실패했습니다: %v", err)
	}
	defer db.Close()

	// DB 연결 상태 확인
	err = db.Ping()
	if err != nil {
		t.Fatalf("DB 연결 상태 확인 실패: %v", err)
	}

	t.Log("DB 연결 성공")
}

func TestSelectQuery(t *testing.T) {
	db, err := openDB()
	if err != nil {
		t.Fatalf("DB 연결에 실패했습니다: %v", err)
	}
	defer db.Close()

	// 쿼리 실행
	rows, err := db.Query("SELECT TOP 1 id, name FROM your_table")
	if err != nil {
		t.Fatalf("쿼리 실행에 실패했습니다: %v", err)
	}
	defer rows.Close()

	// 결과 처리
	var id int
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			t.Fatalf("결과 처리 중 오류 발생: %v", err)
		}
		t.Logf("ID: %d, Name: %s", id, name)
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("결과 처리 중 오류 발생: %v", err)
	}
}
