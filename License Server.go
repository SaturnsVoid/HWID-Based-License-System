package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	PORT int = 9347
)

func checkFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func createFile(pathFile string) error {
	file, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func randomString(n int) string {
	var letterRunes = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//LICENSE:EXPDATE:EMAIL:HWID = 0,1,2
func checkHandler(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	license := request.FormValue("license")
	hwid := request.FormValue("hwid")

	database, _ := readLines("db")
	for _, table := range database {

		row := strings.Split(table, ":")

		t, err := time.Parse("2006-01-02", row[1])
		if err != nil {
			fmt.Println("ERROR: Error reading database")
		}

		t2, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

		if license == row[0] && t.After(t2) {
			if hwid == row[3] {
				fmt.Fprintf(response, "0") //Registed, Good licnese
			} else if row[3] == "NOTSET" {
				b, err := ioutil.ReadFile("db")
				if err != nil {
					fmt.Println("READfromCHECK")
					os.Exit(0)
				}

				str := string(b)
				edit := row[0] + ":" + row[1] + ":" + row[2] + ":" + hwid
				res := strings.Replace(str, table, edit, -1)

				err = ioutil.WriteFile("db", []byte(res), 0644)
				if err != nil {
					fmt.Println("WRITEfromCHECK")
					os.Exit(0)
				}

				fmt.Fprintf(response, "2") //Registed, Good licnese
			}
		} else if license == row[0] && !t.After(t2) {
			fmt.Fprintf(response, "1") //registerd but license experied
		}
	}
}

func serverAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/", checkHandler).Methods("POST")
	http.Handle("/", router)

	http.ListenAndServe(":"+string(strconv.Itoa(PORT)), nil)
}

func main() {
	fmt.Println("License Server")
	fmt.Println("Github: github.com/SaturnsVoid")

	if !checkFileExist("db") {
		fmt.Println("Database does not exist, creating new database.")
		_ = createFile("db")
	}

	database, _ := readLines("db")

	fmt.Println("Total Licenses:", len(database))

	go serverAPI()
	for {
		fmt.Println(" ")
		fmt.Print("$> ")
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		switch scan.Text() {
		case "list":
			database, _ = readLines("db")
			for _, table := range database {
				fmt.Println(table)
			}
		case "add":
			var email string
			var experation string
			var license string

			fmt.Print("License Email: ")
			scan = bufio.NewScanner(os.Stdin)
			scan.Scan()
			email = scan.Text()

		exp:
			fmt.Print("License Experation (YYYY-MM-DD): ")
			scan = bufio.NewScanner(os.Stdin)
			scan.Scan()
			_, err := time.Parse("2006-01-02", scan.Text())
			if err != nil {
				fmt.Println("Experation must be in the YYYY-MM-DD Format.")
				goto exp
			}
			experation = scan.Text()

			license = randomString(4) + "-" + randomString(4) + "-" + randomString(4)

			b, err := ioutil.ReadFile("db")
			if err != nil {
				os.Exit(0)
			}

			str := string(b)
			str = str + "\r\n" + license + ":" + experation + ":" + email + ":NOTSET"

			re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
			str2 := strings.Trim(re.ReplaceAllString(str, ""), "\r\n")

			err = ioutil.WriteFile("db", []byte(str2), 0644)
			if err != nil {
				os.Exit(0)
			}

			fmt.Println("New License Generated:", license, "for", email)
		case "add bulk":
			var experation string

			fmt.Println("Bulk accounts will be added to database without emails. You can add emails at a later time.")
			fmt.Println(" ")

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("How many keys to generate? (#): ")
			bytes, _, err := reader.ReadLine()
			if err != nil {
				os.Exit(0)
			}

			amount := string(bytes)

			n, err := strconv.Atoi(amount)
			if err != nil {
				os.Exit(0)
			}

		expb:
			fmt.Print("License Experation (YYYY-MM-DD): ")
			scan = bufio.NewScanner(os.Stdin)
			scan.Scan()
			_, err = time.Parse("2006-01-02", scan.Text())
			if err != nil {
				fmt.Println("Experation must be in the YYYY-MM-DD Format.")
				goto expb
			}
			experation = scan.Text()

			for i := 0; i < n; i++ {
			restart:
				var old string
				license := randomString(4) + "-" + randomString(4) + "-" + randomString(4)
				if license != old {
					b, err := ioutil.ReadFile("db")
					if err != nil {
						os.Exit(0)
					}

					str := string(b)
					str = str + "\r\n" + license + ":" + experation + ":null" + ":NOTSET"

					re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
					str2 := strings.Trim(re.ReplaceAllString(str, ""), "\r\n")

					err = ioutil.WriteFile("db", []byte(str2), 0644)
					if err != nil {
						os.Exit(0)
					}
					fmt.Println("New License Generated:", license)
					old = license
				} else {
					goto restart
				}
			}

		case "remove":
			fmt.Print("What email would you like to remove?: ")
			scan := bufio.NewScanner(os.Stdin)
			scan.Scan()

			for _, table := range database {

				row := strings.Split(table, ":")

				if scan.Text() == row[2] { //Found in DB

					b, err := ioutil.ReadFile("db")
					if err != nil {
						os.Exit(0)
					}

					str := string(b)
					res := strings.Replace(str, table, "", -1)

					re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
					reres := strings.Trim(re.ReplaceAllString(res, ""), "\r\n")

					err = ioutil.WriteFile("db", []byte(reres), 0644)
					if err != nil {
						os.Exit(0)
					}
				}
			}

			fmt.Println("Done")
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Unknown Command")
		}
	}
}
