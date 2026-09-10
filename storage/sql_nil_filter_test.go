package storage_test

import (
	"testing"

	"github.com/tink3rlabs/magic/storage"
)

type sqlNullItem struct {
	Id        string  `json:"id" gorm:"column:id;primaryKey"`
	Name      string  `json:"name"`
	DeletedAt *string `json:"deleted_at"`
}

func (sqlNullItem) TableName() string { return "sql_null_items" }

func TestSQLAdapterNilOnlyFilterGetListCount(t *testing.T) {
	m := storage.GetMemoryAdapterInstance()
	if err := m.Execute(`CREATE TABLE IF NOT EXISTS sql_null_items (id TEXT PRIMARY KEY, name TEXT, deleted_at TEXT)`); err != nil {
		t.Fatal(err)
	}
	if err := m.Execute(`DELETE FROM sql_null_items`); err != nil {
		t.Fatal(err)
	}
	if err := m.Execute(`INSERT INTO sql_null_items (id, name, deleted_at) VALUES ('1', 'alive', NULL)`); err != nil {
		t.Fatal(err)
	}
	if err := m.Execute(`INSERT INTO sql_null_items (id, name, deleted_at) VALUES ('2', 'gone', '2020-01-01')`); err != nil {
		t.Fatal(err)
	}

	sqlAd := m.DB
	var got sqlNullItem
	if err := sqlAd.Get(&got, map[string]any{"deleted_at": nil}); err != nil {
		t.Fatalf("Get nil-only filter: %v", err)
	}
	if got.Id != "1" {
		t.Fatalf("Get id = %q", got.Id)
	}

	var page []sqlNullItem
	if _, err := sqlAd.List(&page, "id", map[string]any{"deleted_at": nil}, 10, ""); err != nil {
		t.Fatalf("List nil-only filter: %v", err)
	}
	if len(page) != 1 || page[0].Id != "1" {
		t.Fatalf("List page = %#v", page)
	}

	n, err := sqlAd.Count(&sqlNullItem{}, map[string]any{"deleted_at": nil})
	if err != nil {
		t.Fatalf("Count nil-only filter: %v", err)
	}
	if n != 1 {
		t.Fatalf("Count = %d", n)
	}
}
