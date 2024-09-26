package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Project struct {
	Name    string `json:"name"`
	Table   string `json:"table"`
	Creator string `json:"creator"`
}

func GetProjects(db *sql.DB) ([]Project, error) {
	var (
		execSQL = fmt.Sprintf(`SELECT * FROM showmaster.projects`)

		projects []Project
		err      error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying rows: %d\n", err)
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error iterating over rows: %d\n", err)
		return nil, err
	}

	for rows.Next() {
		var project Project
		err = rows.Scan(&project.Name, &project.Table, &project.Creator)
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning rows: %d\n", err)
			return nil, err
		}
		projects = append(projects, project)
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
	cacheProjects(db)

	for _, project := range projects {
		if project == name {
			err = errors.New("project already exists")
			return err
		}
	}

	_, err = db.Exec(execSqlProjectRow)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error inserting new project: %d\n", err)
		return err
	}

	_, err = db.Exec(execSqlProjectTable)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error creating new project table: %d\n", err)
		return err
	}

	return nil
}

func UpdateProject(oldName string, newName string, db *sql.DB) error {
	var (
		execSql = fmt.Sprintf(`UPDATE showmaster.projects SET name = '%s' WHERE name = '%s';`, newName, oldName)

		err error
	)

	_, err = db.Exec(execSql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
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
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error deleting project: %d\n", err)
		return err
	}

	_, err = db.Exec(execSQLTable)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error deleting project table: %d\n", err)
		return err
	}

	return nil
}

// UpdateTimer updates the time stored in the db
func UpdateTimer(t time.Duration, name string, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`UPDATE showmaster.projects SET timer= '%d' WHERE name = '%s'`, t, name)
		err     error
	)
	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error updating timer table: %d\n", err)
		return err
	}

	return nil
}

// Row the row used by the individual projects
type Row struct {
	Id    int            `json:"id"`
	Pos   float32        `json:"pos"`
	Name  string         `json:"name"`
	Audio *string        `json:"audio"`
	Light *string        `json:"light"`
	PPTX  *string        `json:"pptx"`
	Notes *string        `json:"notes"`
	Timer *time.Duration `json:"timer"`
}

// GetRows Get all the rows from a project
func GetRows(project string, db *sql.DB) ([]Row, error) {
	var (
		data    []Row
		execSQL = fmt.Sprintf(`SELECT * FROM showmaster.%s`, project)

		err error
	)
	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying rows: %d\n", err)
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error iterating over rows: %d\n", err)
		return nil, err
	}

	for rows.Next() {
		var row Row
		err = rows.Scan(&row.Id, &row.Pos, &row.Name, &row.Audio, &row.Light, &row.PPTX, &row.Notes, &row.Timer)
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning rows: %d\n", err)
			return nil, err
		}
		data = append(data, row)
	}

	return data, nil
}

// NewRow to insert into a project
func NewRow(project string, r Row, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`INSERT INTO showmaster.%s (pos, name, audio, light, pptx, notes, timer) VALUES (%f, '%s', '%s', '%s', '%s', '%s', %d);`,
			project,
			r.Pos,
			r.Name,
			r.Audio,
			r.Light,
			r.PPTX,
			r.Notes,
			r.Timer)

		err error
	)

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error inserting new row: %d\n", err)
		return err
	}

	return nil
}

// UpdateRow to update an existing row
func UpdateRow(project string, r Row, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`UPDATE showmaster.%s SET pos = %f, name = '%s', audio = '%s', light = '%s', pptx = '%s', notes = '%s', timer = %d WHERE id = %d;`,
			project,
			r.Pos,
			r.Name,
			r.Audio,
			r.Light,
			r.PPTX,
			r.Notes,
			r.Timer,
			r.Id)

		err error
	)

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error updating row: %d\n", err)
		return err
	}

	return nil
}

// DeleteRow to delete a specific row
func DeleteRow(project string, id int, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`DELETE FROM showmaster.%s WHERE id = %d;`, project, id)

		err error
	)

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error deleting row: %d\n", err)
		return err
	}

	return nil
}

// PROJECTS list of all projects currently registered
var projects []string

// CacheProjects caches the projects into a var stored in this file, will be called after every user action
func cacheProjects(db *sql.DB) {
	var (
		execSQL = fmt.Sprintf(`SELECT name FROM showmaster.projects`)

		err error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying rows from showmaster.users: %d\n", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Processing the results
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning row: %d\n", err)
		}
		projects = append(projects, name)
	}

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error with rows: %d\n", err)
	}
}
