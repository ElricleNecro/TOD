package db

import (
	"database/sql"
	"strings"
)

type Computer struct {
	IP, Name string
	Perf     float64
	NbProc   int
}

//Lit la base de donnée "name" pour récupérer les champs "fields" de la table "tables".
//La fonction remplit la structure Computer
func GetInfoFromDB(name, tables string, fields ...string) ([]Computer, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT" + strings.Join(fields, ",") + " FROM " + tables)
	if err != nil {
		return nil, err
	}

	/* for rows.Next() {*/
	//var tmp Computer
	//rows.Scan(
	/* }*/

	return nil, nil
}

//vim: spelllang=en
