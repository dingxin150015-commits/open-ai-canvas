package database

import (
	"bytes"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type loggerFixture struct {
	ID string `gorm:"primaryKey"`
}

func TestGORMLoggerSuppressesNotFoundAndParameterizesErrors(t *testing.T) {
	var output bytes.Buffer
	db, err := gorm.Open(sqlite.Open("file:logger-test?mode=memory&cache=shared"), &gorm.Config{Logger: newGORMLogger(&output)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&loggerFixture{}); err != nil {
		t.Fatal(err)
	}
	var missing loggerFixture
	if err := db.First(&missing, "id = ?", "private-record-id").Error; err == nil {
		t.Fatal("missing row should still return an error")
	}
	if output.Len() != 0 {
		t.Fatalf("record-not-found should not be logged: %q", output.String())
	}
	secret := "stage15-sensitive-query-value"
	if err := db.Raw("SELECT * FROM definitely_missing_table WHERE secret = ?", secret).Scan(&[]map[string]any{}).Error; err == nil {
		t.Fatal("invalid table query should fail")
	}
	logged := output.String()
	if strings.Contains(logged, secret) {
		t.Fatalf("query parameters must be hidden: %q", logged)
	}
	if !strings.Contains(strings.ToLower(logged), "no such table") {
		t.Fatalf("real database errors must remain diagnosable: %q", logged)
	}
	if !strings.Contains(logged, "[SQL statement redacted]") {
		t.Fatalf("SQL text must be replaced with a stable marker: %q", logged)
	}
}
