package mssql

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	defaultServer   = "JIN_SURFACE\\SQLEXPRESS"
	defaultDatabase = "BaseballData"
)

// [데이터베이스 연결 함수]
func OpenDB(server string, database string) (*sql.DB, error) {
	// MSSQL에 연결하기 위한 DSN(Data Source Name)
	dsn := "server=" + server + ";database=" + database + ";integrated security=true"
	return sql.Open("mssql", dsn)
}

// [단일 행 쿼리] : 'QueryRow'
// 결과가 1개인 경우 사용되고, QueryRow 메서드를 사용함
// 단일 행 쿼리는 기대한 대로 결과가 1개만 반환될 것이라고 가정하고 쿼리를 실행.
// QueryRow 메서드는 내부적으로 Query 메서드를 호출하여 쿼리를 실행함.
// 결과로 반환된 행이 1개인 경우 그 행의 데이터를 성공적으로 가져오고, 없는 경우 sql.ErrNoRows 에러 반환
// 만약 여러 행이 반환되면, QueryRow는 첫 번째 행의 데이터만을 가져오고 나머지 행은 무시
// [복수 행 쿼리] : 'Query'
// 쿼리 조건에 맞는 여러 개의 결과가 반환될 것으로 예상되는 경우에 사용함.
// 일반적으로 WHERE절을 사용하여 특정 조건을 만족하는 모든 레코드를 검색

// [단일 행 쿼리 실행 함수]
// 주어진 SQL 쿼리를 실행하여 단일 행의 결과를 반환하는 함수
// 쿼리와 쿼리의 파라미털르 입력받아, 결과를 문자열 슬라이스로 반환
func SELECT_SingleRow(db *sql.DB, query string, args ...interface{}) (*sql.Row, error) {
	if db == nil {
		db, err := OpenDB(defaultServer, defaultDatabase)
		if err != nil {
			return nil, err
		}
		defer db.Close()
	}

	// 주어진 쿼리를 실행하고, 결과를 row로 가져옴
	row := db.QueryRow(query, args...) // args...를 통해 추가적인 파라미터를 쿼리에 전달
	return row, nil
}

// [복수 행 쿼리 실행 함수]
func SELECT_MultipleRow(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		db, err := OpenDB(defaultServer, defaultDatabase)
		if err != nil {
			return nil, err
		}
		defer db.Close()
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
