package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"transaction/common"
)

func main() {
	settings := get_settings.GetSettings("/..")
	createBD(settings)
	createUsersTable(settings)
	createJournalTable(settings)
}

func createBD(settings get_settings.Settings) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword)
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Проверяем на существование
	var datname string
	err = db.QueryRow("SELECT datname FROM pg_database WHERE datname = $1", settings.DBName).Scan(&datname)
	if err == sql.ErrNoRows {
		// Создаем базу
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", settings.DBName))
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	fmt.Printf("База данных %s успешно создана", settings.DBName)
}

func createUsersTable(settings get_settings.Settings) error {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword, settings.DBName)
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := `CREATE TABLE IF NOT EXISTS users(
		id serial4  NOT NULL,
		name varchar(100) NOT NULL, 
		balance integer NOT NULL,
		created_at timestamp default CURRENT_TIMESTAMP,
		updated_at timestamp default CURRENT_TIMESTAMP,
		CONSTRAINT users_pkey PRIMARY KEY (id))`
	res, err := db.Exec(query)
	if err != nil {
		fmt.Printf("\nError %s when creating users table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("\nError %s when getting rows affected", err)
		return err
	}
	fmt.Printf("\nRows affected when creating table: %d", rows)
	return nil
}

func createJournalTable(settings get_settings.Settings) error {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword, settings.DBName)
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := `CREATE TABLE IF NOT EXISTS journal(
		id serial4 NOT NULL ,
		user_id int4 NOT NULL ,
		description varchar(100) NULL, 
		amount integer NOT NULL,
		created_at timestamp default CURRENT_TIMESTAMP,
		updated_at timestamp default CURRENT_TIMESTAMP,
		successful BOOLEAN NOT NULL default FALSE,
		CONSTRAINT journal_pkey PRIMARY KEY (id),
	    CONSTRAINT journal_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id))`
	res, err := db.Exec(query)
	if err != nil {
		fmt.Printf("\nError %s when creating journal table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("\nError %s when getting rows affected", err)
		return err
	}
	fmt.Printf("\nRows affected when creating table: %d", rows)
	return nil
}
