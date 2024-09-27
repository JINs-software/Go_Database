package mssql

import (
	"database/sql"
	"testing"
)

func TestSigleRowQuery(t *testing.T) {
	db, err := OpenDB(defaultServer, defaultDatabase)
	if err != nil {
		t.Fatalf("DB 연결에 실패했습니다: %v", err)
	}

	var query string

	var nameFirst, nameLast, birthCity string
	query = `SELECT nameFirst, nameLast, birthCity
			FROM dbo.players
			WHERE birthYear = 1866`
	row, err := SELECT_SingleRow(db, query)
	if err != nil {
		t.Fatalf("[FAILED] SELECT_SingleRow(db, query): %v", err)
	}
	err = row.Scan(&nameFirst, &nameLast, &birthCity)
	if err != nil {
		t.Fatalf("[FAILED] Scan(..): %v", err)
	}
	t.Logf("Name: %s %s, BirthCity: %s", nameFirst, nameLast, birthCity)
}

func TestMultipleRowQuery(t *testing.T) {
	db, err := OpenDB(defaultServer, defaultDatabase)
	if err != nil {
		t.Fatalf("DB 연결에 실패했습니다: %v", err)
	}

	query := `SELECT nameFirst, nameLast, birthCity
			FROM dbo.players
			WHERE birthYear = 1866`
	rows, err := SELECT_MultipleRow(db, query)
	if err != nil {
		t.Fatalf("[FAILED] SELECT_MultipleRow(db, query): %v", err)
	}
	defer rows.Close()

	//var nameFirst, nameLast, birthCity string
	// -> 널 값 허용 필드를 위해
	var (
		nameFirst sql.NullString
		nameLast  sql.NullString
		birthCity sql.NullString
	)

	for rows.Next() {
		err := rows.Scan(&nameFirst, &nameLast, &birthCity)
		if err != nil {
			t.Fatalf("[ERROR] rows.Scan : %v", err)
		}
		t.Logf(nameFirst.String, nameLast.String, birthCity.String)
	}
}
