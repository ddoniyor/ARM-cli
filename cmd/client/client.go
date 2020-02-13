package main

import (
	"database/sql"
	"fmt"

	"io"
	"log"
	"os"
	"strings"
)



// TODO: для тех, кто хочет попробовать, можете использовать структуры и методы:
type manager struct {
	db  *sql.DB
	out io.Writer
	in  io.Reader
}

func newManagerCLI(db *sql.DB, out io.Writer, in io.Reader) *manager {
	return &manager{db: db, out: out, in: in}
}

// Writer, Reader


func main() {
	// os.Stdin, os.Stout, os.Stderr, File
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
			return false
		}
		if !ok {
			fmt.Println("Неправильно введён логин или пароль. Попробуйте ещё раз.")

			return false
		}
		operationsLoop(db, authorizedOperations, authorizedOperationsLoop)
	case "2":
		atms, err := core.GetAllAtms(db)
		if err != nil {
			log.Printf("can't get all atms: %v", err)
			return true
		}
		printAtm(atms)

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
		accounts, err := core.GetAllAccounts(db)
		if err != nil {
			log.Printf("can't get all accounts: %v", err)
			return true
		}
		printAccounts(accounts)
	case "2":
		err := handlePhoneNumberTransaction(db)
		if err != nil {
			log.Printf("can't pay with phone number: %v", err)
			return true
		}
	case "3":
		err := handleAccountsIdTransaction(db)
		if err != nil {
			log.Printf("can't pay with account: %v", err)
			return true
		}

	case "4":
		services, err := core.GetAllServices(db)
		if err != nil {
			log.Printf("can't get all services: %v", err)
			return true
		}
		printServices(services)
	case "5":
		err := handleServiceTransaction(db)
		if err != nil {
			log.Printf("can't pay service: %v", err)
			return true
		}
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}


//PRINTS STUFFS
func printAtm(atms []core.Atm) {
	for _, atm := range atms {
		fmt.Printf(
			"id: %d, name: %s, adress: %s\n",
			atm.Id,
			atm.Name,
			atm.Adress,
		)
	}
}

func printAccounts(accounts []core.Account) {
	for _, account := range accounts {
		fmt.Printf(
			"id: %d,name: %s, balance: %d, client_id: %d\n",
			account.Id,
			account.Name,
			account.Balance,
			account.Client_id,
		)
	}
}

func printServices(services []core.Service) {
	for _, service := range services {
		fmt.Printf(
			"id: %d, name: %s, balance: %d\n",
			service.Id,
			service.Name,
			service.Balance,
		)
	}
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


	ok, err = core.LoginClient(login, password,db)
	if err != nil {
		return false, err
	}

	return ok, err
}

func handleServiceTransaction(db *sql.DB)( err error) {
	var sum int
	fmt.Println("Введите сумму перевода: ")
	_, err = fmt.Scan(&sum)
	if err != nil {
		return  err
	}

	var idService int
	fmt.Println("Введите id сервиса: ")
	_, err = fmt.Scan(&idService)
	if err != nil {
		return  err
	}

	err =core.TransferMoneyToService(sum,idService, db)
	if err != nil {
		return  err
	}

	fmt.Println("Транзакция успешно выполнена!")
	return nil

}

func handlePhoneNumberTransaction(db *sql.DB)( err error) {
	var sum int
	fmt.Println("Введите сумму перевода: ")
	_, err = fmt.Scan(&sum)
	if err != nil {
		return  err
	}

	var phoneNum int
	fmt.Println("Введите номер клиента: ")
	_, err = fmt.Scan(&phoneNum)
	if err != nil {
		return  err
	}

	err =core.TransferMoneyWithPhoneNumber(sum,phoneNum, db)
	if err != nil {
		return  err
	}

	fmt.Println("Транзакция успешно выполнена!")
	return nil

}

func handleAccountsIdTransaction(db *sql.DB)( err error) {
	var sum int
	fmt.Println("Введите сумму перевода: ")
	_, err = fmt.Scan(&sum)
	if err != nil {
		return  err
	}

	var idAccount int
	fmt.Println("Введите id счета: ")
	_, err = fmt.Scan(&idAccount)
	if err != nil {
		return  err
	}

	err =core.TransferMoneyWithAccountId(sum,idAccount, db)
	if err != nil {
		return  err
	}

	fmt.Println("Транзакция успешно выполнена!")
	return nil

}