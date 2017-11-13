package main

import (
	"database/sql"
)

type device struct {
	Name       string `json:"name"`
	LocationID string `json:"location_id"`
	Location   string `json:"location"`
}

func (d *device) getDevice(db *sql.DB) error {
	return db.QueryRow("SELECT name, location FROM devices WHERE name=$1",
		d.Name).Scan(&d.Name)
}

func (d *device) deleteDevice(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM devices WHERE name=$1", d.Name)

	return err
}

func (d *device) createDevice(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO devices(name, location_id, location) VALUES($1, $2, $3) RETURNING name",
		d.Name, "box", "box").Scan(&d.Name)
}

func (d *device) updateDevice(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE devices SET name=$1, location_id=$2, location=$3 WHERE name=$1",
			d.Name, d.LocationID, d.Location)

	return err
}

func getDevices(db *sql.DB) ([]device, error) {
	rows, err := db.Query(
		"SELECT name, location_id, location FROM devices")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := []device{}

	for rows.Next() {
		var d device
		if err := rows.Scan(&d.Name, &d.LocationID, &d.Location); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}
