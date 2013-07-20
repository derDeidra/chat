package main

import (
	"database/sql"
	"time"
)

var (
	insertstatement *sql.Stmt
	banstatement    *sql.Stmt
	unbanstatement  *sql.Stmt
	// TODO privmsgstatement *sql.Stmt
)

func initEventlog() {

	var err error
	insertstatement, err = db.Prepare(`
		INSERT INTO chatlog
		SET
			userid       = ?,
			targetuserid = ?,
			event        = ?,
			data         = ?,
			timestamp    = ?
	`)

	if err != nil {
		B("Unable to create insert statement: ", err)
	}

	banstatement, err = db.Prepare(`
		INSERT INTO bans
		SET
			userid         = ?,
			targetuserid   = ?,
			ipaddress      = ?,
			reason         = ?,
			starttimestamp = ?,
			endtimestamp   = ?
	`)

	unbanstatement, err = db.Prepare(`
		UPDATE bans
		SET endtimestamp = NOW()
		WHERE
			targetuserid = ? AND
			(
				endtimestamp IS NULL OR
				endtimestamp > NOW()
			)
	`)

	if err != nil {
		B("Unable to create ban statement: ", err)
	}

}

type Event struct {
	Userid       Userid `json:"userid"`
	Targetuserid Userid `json:"targetuserid,omitempty"`
	Event        string `json:"event"`
	Data         string `json:"data,omitempty"`
	Timestamp    int64  `json:"timestamp"`
}

func logEvent(userid Userid, event string, data *EventDataOut) {

	targetuserid := &sql.NullInt64{}
	if data.Targetuserid != 0 {
		targetuserid.Int64 = int64(data.Targetuserid)
		targetuserid.Valid = true
	}

	d := &sql.NullString{}
	if len(data.Data) != 0 {
		d.String = data.Data
		d.Valid = true
	}

	// the timestamp is milisecond precision
	ts := time.Unix(data.Timestamp/1000, 0).UTC()
	_, err := insertstatement.Exec(userid, targetuserid, event, d, ts)
	if err != nil {
		D("Unable to insert event: ", err)
	}

}
