package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int64
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int64, error) {
	sqlStr := "INSERT INTO t_snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"

	result, err := m.DB.Exec(sqlStr, title, content, expires)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	snippet := &Snippet{}
	sqlStr := "SELECT id, title, content, created, expires FROM t_snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"

	err := m.DB.QueryRow(sqlStr, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return snippet, nil
}

// 返回最新的10条记录
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	sqlStr := "SELECT id, title, content, created, expires FROM t_snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10"

	rows, err := m.DB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*Snippet{}
	for rows.Next() {
		snippet := &Snippet{}
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		result = append(result, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
