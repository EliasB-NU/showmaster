package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"showmaster/src/util"
	"time"
)

type Project struct {
	Name    string
	Table   string
	Creator string
}

func GetProjects(db *sql.DB) ([]Project, error) {
	var (
		execSQL = fmt.Sprintf(`SELECT * FROM showmaster.projects`)

		projects []Project
		err      error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error querying rows: %d\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var project Project
		err = rows.Scan(&project.Name, &project.Table, &project.Creator)
		if err != nil {
			log.SetFlags(log.LstdFlags & log.Lshortfile)
			log.Printf("Error scanning rows: %d\n", err)
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return projects, nil
}

func NewProject(name string, creator string, db *sql.DB) error {
	var (
		execSqlProjectRow = fmt.Sprintf(`
			INSERT INTO showmaster.projects (name, projecttable, creator) VALUES ('%s', '%s','%s');`, name, name+"table", creator)

		execSqlProjectTable = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS showmaster.%s (
			    id SERIAL PRIMARY KEY,
			    pos DOUBLE PRECISION,
			    name TEXT NOT NULL,
			    audio TEXT,
			    light TEXT,
			    pptx TEXT,
			    notes TEXT,
			    timer interval
			);`, name+"table")

		err error
	)

	for _, project := range util.PROJECTS {
		if project == name {
			err = errors.New("project already exists")
			return err
		} else {
			continue
		}
	}

	_, err = db.Exec(execSqlProjectRow)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Error inserting new project: %d\n", err)
		return err
	}

	_, err = db.Exec(execSqlProjectTable)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Error creating new project table: %d\n", err)
		return err
	}

	util.CacheProjects(db)
	return nil
}

func UpdateProject(oldName string, newName string, db *sql.DB) error {
	var (
		execSql = fmt.Sprintf(`UPDATE showmaster.projects SET name = '%s' WHERE name = '%s';`, newName, oldName)

		err error
	)

	_, err = db.Exec(execSql)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Error updating project table: %d\n", err)
		return err
	}

	return nil
}

func DeleteProject(name string, db *sql.DB) error {
	var (
		execSQLRow   = fmt.Sprintf(`DELETE FROM showmaster.projects WHERE name = '%s';`, name)
		execSQLTable = fmt.Sprintf(`DROP TABLE showmaster.%s;`, name+"table")
		err          error
	)

	_, err = db.Exec(execSQLRow)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Error deleting project: %d\n", err)
		return err
	}

	_, err = db.Exec(execSQLTable)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Error deleting project table: %d\n", err)
		return err
	}

	return nil
}

// Row the row used by the individual projects
type Row struct {
	id    int
	pos   float32
	name  string
	audio string
	light string
	pptx  string
	notes string
	timer time.Duration
}

// NewRow to insert into a project
func NewRow(name string, r Row, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`INSERT INTO showmaster.%s (pos, name, audio, light, pptx, notes, timer) VALUES (%f, '%s', '%s', '%s', '%s', '%s', %d);`,
			name+"table",
			r.pos,
			r.name,
			r.audio,
			r.light,
			r.pptx,
			r.notes,
			r.timer)

		err error
	)

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Error inserting new row: %d\n", err)
		return err
	}

	return nil
}

// UpdateRow to update an existing row
func UpdateRow(name string, r Row, id int, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`UPDATE showmaster.%s SET pos = %f, name = '%s', audio = '%s', light = '%s', pptx = '%s', notes = '%s', timer = %d WHERE id = %d;`,
			name+"table",
			r.pos,
			r.name,
			r.audio,
			r.light,
			r.pptx,
			r.notes,
			r.timer,
			id)

		err error
	)

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Error updating row: %d\n", err)
		return err
	}

	return nil
}

// DeleteRow to delete a specific row
func DeleteRow(name string, id int, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`DELETE FROM showmaster.%s WHERE id = %d;`, name+"table", id)

		err error
	)

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Fatalf("Error deleting row: %d\n", err)
		return err
	}

	return nil
}
