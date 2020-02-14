package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ddoniyor/ARM-core/pkg/core"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Clients struct {
	Id       int
	Name     string
	Login    string
	Password string
	PhoneNum int
}

type Atm struct {
	Id     int
	Name   string
	Adress string
}

type Accounts struct {
	Id        int
	Name      string
	Balance   int
	Client_id int
}

func main() {

	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Print("start application")
	log.Print("open db")
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}
	defer func() {
		log.Print("close db")
		if err := db.Close(); err != nil {
			log.Fatalf("can't close db: %v", err)
		}
	}()
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init db: %v", err)
	}
	fmt.Fprintln(os.Stdout, "Добро пожаловать в наше приложение")
	log.Print("start operations loop")
	operationsLoop(db, unauthorizedOperations, unauthorizedOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")

	jsonFormatClients(db)
	jsonFormatAtms(db)
	jsonFormatAccounts(db)

}

func operationsLoop(db *sql.DB, commands string, loop func(db *sql.DB, cmd string) bool) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err) // %v - natural ...
		}
		if exit := loop(db, strings.TrimSpace(cmd)); exit {
			return
		}
	}
}

func unauthorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		ok, err := handleLogin(db)
		if err != nil {
			log.Printf("can't handle login: %v", err)
			return true
		}
		if !ok {
			fmt.Println("Неправильно введён логин или пароль. Попробуйте ещё раз.")

			return false
		}
		operationsLoop(db, authorizedOperations, authorizedOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}

func authorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		err := handleAtm(db)
		if err != nil {
			log.Printf("can't add atm: %v", err)
			return true
		}
	case "2":
		err := handleClients(db)
		if err != nil {
			log.Printf("can't add client: %v", err)
			return true
		}
	case "3":
		err := handleAccount(db)
		if err != nil {
			log.Printf("can't add account: %v", err)
			return true
		}
	case "4":
		err := handleService(db)
		if err != nil {
			log.Printf("can't add service: %v", err)
			return true
		}

	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

//HANDLE STUFFS

func handleLogin(db *sql.DB) (ok bool, err error) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return false, err
	}
	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return false, err
	}

	ok, err = core.Login(login, password, db)
	if err != nil {
		return false, err
	}

	return ok, err
}

func handleAtm(db *sql.DB) (err error) {
	fmt.Println("Введите данные банкомата")

	var name string
	fmt.Print("Имя: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}

	fmt.Print("Адресс: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Can't read command: %v", err)
	}
	fmt.Println(text)

	err = core.AddAtm(name, text, db)
	if err != nil {
		return err
	}

	return nil
}

func handleClients(db *sql.DB) (err error) {
	fmt.Println("Введите данные клиента")

	var name string
	fmt.Print("Имя: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}

	var log string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&log)
	if err != nil {
		return err
	}

	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return err
	}
	var phone int64
	fmt.Print("Номер телефона: ")
	_, err = fmt.Scan(&phone)
	if err != nil {
		return err
	}

	err = core.AddClients(name, log, password, phone, db)
	if err != nil {
		return err
	}
	fmt.Println("Пользователь успешно добавлен!")

	return nil
}

func handleService(db *sql.DB) (err error) {
	fmt.Println("Введите данные услуги:")
	var name string
	fmt.Print("Имя услуги: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}
	var balance int
	balance = 0
	err = core.AddService(name, balance, db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Услуга успешно добавлена!")

	return nil
}

func handleAccount(db *sql.DB) (err error) {
	fmt.Println("Введите данные счёта:")
	var name string
	fmt.Print("Название счёта: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}

	var balance int64
	fmt.Print("Пополните счёт клиента: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	var client_id int64
	fmt.Print("Введите id владелец счёта: ")
	_, err = fmt.Scan(&client_id)
	if err != nil {
		fmt.Println("Нет такого id")
		log.Fatal(err)
	}

	err = core.AddAccount(name, balance, client_id, db)
	if err != nil {
		return err
	}

	fmt.Println("Счёт успешно добавлен!")

	return nil
}

//JSON FORMAT

func jsonFormatClients(db *sql.DB) {
	client := []Clients{{4, "Amir", "amir", "amir", 902222201},
		{5, "Lola", "lola", "lola", 9022002020},
	}
	data, err := json.Marshal(client)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("Clients.json", data, 0666, )
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadFile("Clients.json")
	if err != nil {
		log.Fatal(bytes)
	}
	err = json.Unmarshal(data, &client)
	if err != nil {
		log.Fatal(err)
	}

	for _, clients := range client {
		_, err = db.Exec(`INSERT INTO clients VALUES (:id, :name, :login, :password, :phoneNum) ON CONFLICT DO NOTHING `,
			sql.Named("id", clients.Id),
			sql.Named("name", clients.Name),
			sql.Named("login", clients.Login),
			sql.Named("password", clients.Password),
			sql.Named("phoneNum", clients.PhoneNum),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func jsonFormatAtms(db *sql.DB) {
	atm := []Atm{
		{5, "T5", "Ayni"},
		{6, "T6", "Shohmansur"},
	}

	data, err := json.Marshal(atm)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("Atms.json", data, 0666, )
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadFile("Atms.json")
	if err != nil {
		log.Fatal(bytes)
	}
	err = json.Unmarshal(data, &atm)
	if err != nil {
		log.Fatal(err)
	}

	for _, atms := range atm {
		_, err = db.Exec(`INSERT INTO atm VALUES (:id, :name, :adress) ON CONFLICT DO NOTHING `,
			sql.Named("id", atms.Id),
			sql.Named("name", atms.Name),
			sql.Named("adress", atms.Adress),

		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func jsonFormatAccounts(db *sql.DB) {
	account := []Accounts{
		{4, "Alif", 2000, 4},
		{5, "Milli", 3000, 5},
	}

	/*exec, err := db.Exec(`SELECT name FROM clients `)
	if err != nil {
		log.Fatal(err)
	}*/

	data, err := json.Marshal(account)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("Accounts.json", data, 0666, )
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadFile("Accounts.json")
	if err != nil {
		log.Fatal(bytes)
	}
	err = json.Unmarshal(data, &account)
	if err != nil {
		log.Fatal(err)
	}

	for _, accounts := range account {
		_, err = db.Exec(`INSERT INTO accounts VALUES (:id, :name, :balance, :client_id) ON CONFLICT DO NOTHING `,
			sql.Named("id", accounts.Id),
			sql.Named("name", accounts.Name),
			sql.Named("balance", accounts.Balance),
			sql.Named("client_id", accounts.Client_id),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

}

