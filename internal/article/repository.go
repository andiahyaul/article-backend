package article

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("article not found")

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const cols = "id, title, content, category, status, created_date, updated_date"

func scan(s interface{ Scan(...any) error }) (Article, error) {
	var a Article
	var c, u sql.NullTime
	if err := s.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.Status, &c, &u); err != nil {
		return a, err
	}
	a.CreatedDate = c.Time.Format("2006-01-02T15:04:05Z07:00")
	a.UpdatedDate = u.Time.Format("2006-01-02T15:04:05Z07:00")
	return a, nil
}

func (r *Repository) Create(a Article) (int64, error) {
	res, err := r.db.Exec(`INSERT INTO posts (title, content, category, status) VALUES (?, ?, ?, ?)`,
		a.Title, a.Content, a.Category, a.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// List returns one page, newest first, plus the total count. Empty status = no filter.
func (r *Repository) List(limit, offset int, status string) ([]Article, int, error) {
	where, args := "", []any{}
	if status != "" {
		where, args = " WHERE status = ?", append(args, status)
	}
	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM posts"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query("SELECT "+cols+" FROM posts"+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []Article{}
	for rows.Next() {
		a, err := scan(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, a)
	}
	return out, total, rows.Err()
}

func (r *Repository) Get(id int) (Article, error) {
	a, err := scan(r.db.QueryRow("SELECT "+cols+" FROM posts WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return a, ErrNotFound
	}
	return a, err
}

func (r *Repository) Update(id int, a Article) error {
	if _, err := r.Get(id); err != nil {
		return err
	}
	_, err := r.db.Exec(`UPDATE posts SET title=?, content=?, category=?, status=? WHERE id=?`,
		a.Title, a.Content, a.Category, a.Status, id)
	return err
}

func (r *Repository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
