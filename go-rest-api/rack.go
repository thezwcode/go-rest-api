// rack.go
package main

import (
	"database/sql"
	"errors"
)

type rack struct {
	ID       int    `json:"id"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	Location string `json:"location"`
}

func (p *rack) createRack(db *sql.DB) error {
	sqlStatement := "INSERT INTO racks(height, width, location) VALUES($1, $2, $3) RETURNING id"
	row := db.QueryRow(sqlStatement, p.Height, p.Width, p.Location)
	err := row.Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func (p *rack) getRack(db *sql.DB) error {
	sqlStatement := `SELECT * FROM racks WHERE id=$1`
	row := db.QueryRow(sqlStatement, p.ID)
	err := row.Scan(&p.ID, &p.Height, &p.Width, &p.Location)
	return err
}

func (p *rack) updateRack(db *sql.DB) error {
	if p.Height == 0 || p.Width == 0 || p.Location == "" {
		return errors.New("some needed variables (height/width/location) are empty")
	}
	sqlStatement := "UPDATE racks SET height=$1, width=$2, location=$3 WHERE id=$4"
	_, err := db.Exec(sqlStatement, p.Height, p.Width, p.Location, p.ID)
	return err
}
