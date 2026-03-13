package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewConnection(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func RunMigrations(db *sqlx.DB, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}

		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			if strings.TrimSpace(q) == "" {
				continue
			}
			if _, err := db.Exec(q); err != nil {
				return fmt.Errorf("migration %s failed: %w", file.Name(), err)
			}
		}
		log.Printf("Migration applied: %s", file.Name())
	}
	return nil
}
