package models

import (
    "database/sql"
    "time"
)

type {{.ModelName}} struct {
    ID        uint      `json:"id" db:"id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new {{.ModelName}} into the database
func (m *{{.ModelName}}) Create(db *sql.DB) error {
    query := `INSERT INTO {{.TableName}} (created_at, updated_at) VALUES (?, ?) RETURNING id`
    now := time.Now()
    m.CreatedAt = now
    m.UpdatedAt = now
    
    err := db.QueryRow(query, m.CreatedAt, m.UpdatedAt).Scan(&m.ID)
    return err
}

// Update modifies an existing {{.ModelName}} in the database
func (m *{{.ModelName}}) Update(db *sql.DB) error {
    query := `UPDATE {{.TableName}} SET updated_at = ? WHERE id = ?`
    m.UpdatedAt = time.Now()
    
    _, err := db.Exec(query, m.UpdatedAt, m.ID)
    return err
}

// Delete removes a {{.ModelName}} from the database
func (m *{{.ModelName}}) Delete(db *sql.DB) error {
    query := `DELETE FROM {{.TableName}} WHERE id = ?`
    _, err := db.Exec(query, m.ID)
    return err
}

// FindByID retrieves a {{.ModelName}} by its ID
func Find{{.ModelName}}ByID(db *sql.DB, id uint) (*{{.ModelName}}, error) {
    query := `SELECT id, created_at, updated_at FROM {{.TableName}} WHERE id = ?`
    
    m := &{{.ModelName}}{}
    err := db.QueryRow(query, id).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
    if err != nil {
        return nil, err
    }
    
    return m, nil
}

// FindAll retrieves all {{.ModelName}} records
func FindAll{{.ModelName}}s(db *sql.DB) ([]*{{.ModelName}}, error) {
    query := `SELECT id, created_at, updated_at FROM {{.TableName}}`
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var models []*{{.ModelName}}
    for rows.Next() {
        m := &{{.ModelName}}{}
        err := rows.Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
        if err != nil {
            return nil, err
        }
        models = append(models, m)
    }
    
    return models, nil
}
