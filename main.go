package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
	"transaction/common"
)

type Users struct {
	UserID     string `json:"id"`
	Name       string `json:"name"`
	Balance    int    `json:"balance"`
	Created_at string `json:"created_at"`
}

type Transactions struct {
	User   string `json:"user"`
	Amount int    `json:"amount"`
}

type Journal struct {
	JournalID  string `json:"id"`
	UserID     string `json:"user_id"`
	Amount     int    `json:"amount"`
	CreatedAt string `json:"created_at"`
	SuccessTask string `json:"success_task"`
	SuccessOperation string `json:"success_operation"`
}

func task(ctx context.Context) {
	// запускаем бесконечный цикл
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return

		// выполняем код
		default:
			println("Запущена фоновая гоурутина...")
			journalList := getJournal(false)
			for _, journal := range journalList {
				description, err := transaction(journal)
				if err != nil {
					fmt.Println(err)
				}
				updateJournal(description, journal)
			}
		}
		// делаем паузу перед следующей итерацией
		time.Sleep(5 * time.Second)
	}
}

func handleRequest() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//получаем POST параметр
		r.ParseForm()
		x := r.Form.Get("parameter_name")
		fmt.Println(x)

		http.ServeFile(w, r, "static/index.html")
	})

	// Получение списка пользователей
	http.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		users := getUsers()
		jsonData, _ := json.Marshal(users)
		fmt.Fprint(w, string(jsonData))
	})

	// Получение журнала
	http.HandleFunc("/api/v1/journal", func(w http.ResponseWriter, r *http.Request) {
		journal := getJournal(true)
		jsonData, _ := json.Marshal(journal)
		fmt.Fprint(w, string(jsonData))
	})

	// Создание пользователя
	http.HandleFunc("/api/v1/ctreate_user", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		responce, err := createUser(name)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
		}
		fmt.Fprint(w, responce)
	})

	// Подтверждение транзакции
	http.HandleFunc("/api/v1/transaction", func(w http.ResponseWriter, r *http.Request) {
		detailTransaction := Transactions{}
		err := json.NewDecoder(r.Body).Decode(&detailTransaction)
		if err != nil {
			fmt.Println(err)
		}
		writeJournal(detailTransaction)
		fmt.Fprint(w, "Данные записаны в журнал")
	})

	http.ListenAndServe(":8080", nil)
}

func main() {
	// создаём контекст с функцией завершения
	ctx, cancel := context.WithCancel(context.Background())
	// запускаем нашу горутину
	go task(ctx)

	handleRequest()

	// завершаем контекст, чтобы завершить горутину
	cancel()
}

func transaction(detail Journal) (string, error) {
	settings := get_settings.GetSettings("")

	// Подключаемся к существующей базе данных
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword, settings.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Начинаем транзакцию
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Откладываем откат на случай, если что-то пойдет не так.
	defer tx.Rollback()

	if detail.Amount < 0 {
		// Подтверждаем, что денег хватает
		enough, err := enoughMoney(tx, ctx, detail.Amount, detail.UserID)
		if !enough {
			fmt.Println("\n", (err), "\n Не хватает денег!")
			return "Не хватает денег!", errors.New("Не хватает денег!")
		}
	}

	// Обновляем баланс
	err = updateMoney(tx, ctx, detail.Amount, detail.UserID)
	if err != nil {
		fmt.Println("\n", (err), "\n ....Откат транзакции!")
		return "Откат транзакции!", err
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return "Успешно!", nil
}

func enoughMoney(tx *sql.Tx, ctx context.Context, money int, userID string) (bool, error) {
	enough := false
	query := fmt.Sprintf("SELECT (balance >= %d) from users where id=%s", -money, userID)
	err := tx.QueryRowContext(ctx, query).Scan(&enough)

	return enough, err
}

func updateMoney(tx *sql.Tx, ctx context.Context, money int, userID string) error {
	query := fmt.Sprintf("UPDATE users SET balance = balance + %d WHERE id=%s", money, userID)
	_, err := tx.ExecContext(ctx, query)

	return err
}

func updateJournal(description string, journal Journal) error {
	settings := get_settings.GetSettings("")

	// Подключаемся к существующей базе данных
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword, settings.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
    
	success_operation := false
	if description == "Успешно!" {
		success_operation = true 
	}
	query := fmt.Sprintf("UPDATE journal SET description='%s', success_task=true, success_operation=%t WHERE id=%s", description, success_operation, journal.JournalID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func getUsers() []Users {
	settings := get_settings.GetSettings("")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword, settings.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, name, balance, created_at FROM users ORDER BY id")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := []Users{}
	for rows.Next() {
		u := Users{}
		err := rows.Scan(&u.UserID, &u.Name, &u.Balance, &u.Created_at)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u)
	}

	return users
}

func createUser(name string) (string, error) {
	settings := get_settings.GetSettings("")

	// Подключаемся к существующей базе данных
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword, settings.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO users (name, balance) VALUES ('%s', 0)", name)
	_, err = db.Exec(query)
	if err != nil {
		return "Ошибка", err
	}

	return "Ok", nil
}

func writeJournal(detailTransaction Transactions) (string, error) {
	settings := get_settings.GetSettings("")

	// Подключаемся к существующей базе данных
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword, settings.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO journal (user_id, amount) VALUES ('%s', %d)", detailTransaction.User, detailTransaction.Amount)
	_, err = db.Exec(query)
	if err != nil {
		return "Ошибка", err
	}

	return "Ok", nil
}

func getJournal(all bool) []Journal {
	settings := get_settings.GetSettings("")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUsername, settings.DBPassword, settings.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	query := `SELECT id, user_id, amount, created_at, success_task, success_operation 
	          FROM journal 
			  ORDER BY created_at`
	if !all {
		query = `SELECT id, user_id, amount, created_at, success_task, success_operation 
				FROM journal 
				WHERE success_task = FALSE
				ORDER BY user_id, created_at`
	}
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	journal := []Journal{}
	for rows.Next() {
		j := Journal{}
		err := rows.Scan(&j.JournalID, &j.UserID, &j.Amount, &j.CreatedAt, &j.SuccessTask, &j.SuccessOperation)
		if err != nil {
			fmt.Println(err)
			continue
		}
		journal = append(journal, j)
	}

	return journal
}
