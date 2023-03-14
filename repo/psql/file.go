package psql

import (
	"database/sql"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
)

type file struct {
	db *sql.DB
}

func NewPsqlFileRepo(db *sql.DB) repo.FileRepo {
	return &file{db}
}

func (self *file) SelectAllByUserId(id int) ([]model.File, error) {
	rows, err := self.db.Query(`
		SELECT	 id
			,path
			,ext
			,user_id
			,created_at
			,updated_at
		FROM m_file
		WHERE user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []model.File{}
	for rows.Next() {
		file := model.File{}
		err = rows.Scan(&file.ID, &file.Path, &file.Ext,
			&file.User.ID, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return files, nil
}

// BROKEN
func (self *file) Create(file *model.File) error {
	stmt, err := self.db.Prepare(`
		INSERT INTO m_file(
			 path
			,ext
			,user_id
			,created_at
			,updated_at
		)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	file.CreatedAt = time.Now()
	file.UpdatedAt = file.CreatedAt

	_, err = stmt.Exec(file.Path, file.Ext, file.User.ID, file.CreatedAt,
		file.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// BROKEN
func (self *file) Update(file *model.File) error {
	stmt, err := self.db.Prepare(`
		UPDATE m_file SET
			 path=?
			,ext=?
			,user_id=?
			,updated_at=?
		WHERE id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	file.UpdatedAt = time.Now()

	_, err = stmt.Exec(file.Path, file.Ext, file.User.ID, file.UpdatedAt,
		file.ID)
	if err != nil {
		return err
	}

	return nil
}

// BROKEN
func (self *file) Delete(id int) error {
	stmt, err := self.db.Prepare("UPDATE m_file SET deleted = ?, deleted_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(true, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (self *file) Search(userId int, query string) ([]model.File, error) {
	query = "%" + query + "%"
	files := []model.File{}
	rows, err := self.db.Query(`
        SELECT id, path, ext, user_id, created_at, updated_at
        FROM m_file
        WHERE user_id = $1 AND (path ILIKE $2 OR ext ILIKE $2)
    `, userId, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := model.File{}
		err = rows.Scan(&file.ID, &file.Path, &file.Ext,
			&file.User.ID, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, repo.ErrRepoNoData
	}

	return files, nil
}
