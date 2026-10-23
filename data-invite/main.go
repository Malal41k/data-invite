package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib" // Драйвер для PostgreSQL
)

// Строка подключения к базе данных
// ВАЖНО: замени "твой_пароль" на свой пароль от postgres!
const connStr = "postgres://postgres:486250@localhost:5432/invite_db?sslmode=disable"

func main() {
	// 1. Подключаемся к базе данных
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения к БД:", err)
		return
	}
	defer db.Close()

	// Проверяем, что связь реально работает
	err = db.Ping()
	if err != nil {
		fmt.Println("База данных не отвечает:", err)
		return
	}
	fmt.Println("Успешно подключено к PostgreSQL!")

	// --- МАРШРУТЫ ---

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="utf-8">
				<title>Приглашение</title>
				<link rel="stylesheet" href="/static/css/style.css">
			</head>
			<body>
				<h1>Привет! У меня к тебе вопрос...</h1>
				<h2>Пойдёшь со мной на свидание?</h2>
				<a href="/yes"><button>Да</button></a>
				<button id="no-btn">Нет</button>
				<script src="/static/js/script.js"></script>
			</body>
			</html>
		`)
	})

	http.HandleFunc("/yes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="utf-8">
				<title>Ура!</title>
				<link rel="stylesheet" href="/static/css/style.css">
			</head>
			<body>
				<h1>Ура! Я знал, что ты согласишься ❤️</h1>
				<p>Теперь выбери дату и ресторан:</p>
				<form action="/confirm" method="POST">
					<p>
						<label>Дата:</label>
						<input type="date" name="date" required>
					</p>
					<p>
						<label>Ресторан:</label>
						<select name="restaurant">
							<option>Пушкин</option>
							<option>Тюбетей</option>
							<option>Приют холостяка</option>
							<option>Марусовка</option>
							<option>Дом чая</option>
						</select>
					</p>
					<p><button type="submit">Отправить</button></p>
				</form>
			</body>
			</html>
		`)
	})

	// СОХРАНЕНИЕ В БАЗУ ДАННЫХ
	http.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		date := r.FormValue("date")
		restaurant := r.FormValue("restaurant")

		// SQL запрос на вставку данных ($1 и $2 - это placeholders)
		_, err := db.Exec("INSERT INTO answers (date, restaurant) VALUES ($1, $2)", date, restaurant)
		if err != nil {
			fmt.Println("Ошибка сохранения в БД:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="utf-8">
				<title>Договорились!</title>
				<link rel="stylesheet" href="/static/css/style.css">
			</head>
			<body>
				<h1>Договорились! ❤️</h1>
				<p>Жду тебя <b>%s</b> в ресторане <b>%s</b>.</p>
				<p><a href="/"><button>На главную</button></a></p>
			</body>
			</html>
		`, date, restaurant)
	})

	// ЧТЕНИЕ ИЗ БАЗЫ ДАННЫХ
	http.HandleFunc("/answer", func(w http.ResponseWriter, r *http.Request) {
		var date, restaurant string
		// Берем последний добавленный ответ (сортируем по id по убыванию)
		err := db.QueryRow("SELECT date, restaurant FROM answers ORDER BY id DESC LIMIT 1").Scan(&date, &restaurant)

		var content string
		if err != nil {
			content = "Пока никто не ответил или произошла ошибка."
			fmt.Println("Ошибка чтения из БД:", err)
		} else {
			content = fmt.Sprintf("Дата: %s\nРесторан: %s\n", date, restaurant)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="utf-8">
				<title>Ответ</title>
				<link rel="stylesheet" href="/static/css/style.css">
			</head>
			<body>
				<h1>Секретный ответ 💌</h1>
				<pre style="font-size: 18px;">%s</pre>
				<p><a href="/"><button>На главную</button></a></p>
			</body>
			</html>
		`, content)
	})

	// Статика (CSS и JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Сервер запущен: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
