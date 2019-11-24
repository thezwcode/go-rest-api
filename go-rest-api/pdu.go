// pdu.go
package main

import (
	"database/sql"
	"errors"
)

type pdu struct {
	ID            int `json:"id"`
	OutletCount   int `json:"outlet_count"`
	OutletUsed    int `json:"outlet_used"`
	PowerCapacity int `json:"power_capacity"`
}

func (p *pdu) createPDU(db *sql.DB) error {
	sqlStatement := "INSERT INTO pdus(outletcount, outletused, powercapacity) VALUES($1, $2, $3) RETURNING id"
	row := db.QueryRow(sqlStatement, p.OutletCount, p.OutletUsed, p.PowerCapacity)
	err := row.Scan(&p.ID, &p.OutletCount, &p.OutletUsed, &p.PowerCapacity)

	if err != nil {
		return err
	}

	return nil
}

func (p *pdu) getPDU(db *sql.DB) error {
	sqlStatement := `SELECT * FROM pdus WHERE id=$1`
	row := db.QueryRow(sqlStatement, p.ID)
	err := row.Scan(&p.ID, &p.OutletCount, &p.OutletUsed, &p.PowerCapacity)
	return err
}

func (p *pdu) updatePDU(db *sql.DB) error {
	if p.OutletCount == 0 || p.OutletUsed == 0 || p.PowerCapacity == 0 {
		return errors.New("some needed variables (height/width/location) are empty")
	}
	sqlStatement := "UPDATE pdus SET outletcount=$1, outletused=$2, powercapacity=$3 WHERE id=$4"
	_, err := db.Exec(sqlStatement, p.OutletCount, p.OutletUsed, p.PowerCapacity, p.ID)
	return err
}
