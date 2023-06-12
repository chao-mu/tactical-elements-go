package data;

import (
	"database/sql"
)

type Player struct {
	ID        int    `json:"id"`
	Health    int    `json:"health"`
	Exp       int    `json:"exp"`
	Gold      int    `json:"gold"`
	BrowserID string `json:"browserId"`
	Dead      bool   `json:"dead"`
}

func Damage(tx *sql.Tx, id int, n int) error {
	_, err := tx.Exec("UPDATE players SET health = max(health - ?, 0) WHERE id = ?", n, id)

	return err
}

func Gain(tx *sql.Tx, id int, gold int, exp int) error {
	_, err := tx.Exec("UPDATE players SET gold = gold + ?, exp = exp + ? WHERE id = ?", gold, exp, id)

	return err
}

func GetPlayer(tx *sql.Tx, id int) (Player, error) {
	p := Player{}
	row := tx.QueryRow("SELECT health, exp, gold, browser_id FROM players WHERE id = ?", id)
	err := row.Scan(&p.Health, &p.Exp, &p.Gold, &p.BrowserID)
	if err != nil {
		return p, err
	}

	p.Dead = p.Health <= 0

	return p, nil
}
