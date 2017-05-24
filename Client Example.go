package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
)

var licenseServer string = "http://127.0.0.1:9347/" //Your license server address

func checkFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func licenseCheck() {
	if !checkFileExist("license.dat") {

		fmt.Println("License file not found.")

		fmt.Print("Would you like to register?: ")
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()

		if scan.Text() == "yes" || scan.Text() == "Yes" || scan.Text() == "YES" {
			fmt.Println("If you do not have a key contact your supplier.")
			fmt.Print("Please entor your key: ")
			scan := bufio.NewScanner(os.Stdin)
			scan.Scan()

			key := scan.Text()

			name, _ := os.Hostname()
			usr, _ := user.Current()

			fmt.Println("Key:", key)

			ioutil.WriteFile("license.dat", []byte(key), 0600)

			fmt.Println("HWID:", md5Hash(name+usr.Username))

			fmt.Println("Connecting to license server...")

			client := &http.Client{}
			data := url.Values{}
			data.Set("license", key)
			data.Set("hwid", md5Hash(name+usr.Username))
			u, _ := url.ParseRequestURI(licenseServer)
			urlStr := fmt.Sprintf("%v", u)
			r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			resp, err := client.Do(r)
			if err != nil {
				fmt.Println("Unable to connect to license server.")
				os.Exit(0)
			} else {
				defer resp.Body.Close()
				resp_body, _ := ioutil.ReadAll(resp.Body)
				if resp.StatusCode == 200 {
					if string(resp_body) != "0" {
						if string(resp_body) == "1" {
							fmt.Println("License is Expired.")
							os.Exit(0)
						} else if string(resp_body) == "2" {
							fmt.Println("Registered!")
							fmt.Println("DO NOT DELETE THE 'license.dat' FILE!")
							fmt.Println(" ")
						} else {
							fmt.Println("Unable to verify to license server.")
							os.Exit(0)
						}
					}
				} else {
					fmt.Println(resp.StatusCode)
				}
			}

		} else {
			os.Exit(0)
		}
	} else {
		name, _ := os.Hostname()
		usr, _ := user.Current()

		dat, _ := ioutil.ReadFile("license.dat")

		key := string(dat)

		client := &http.Client{}
		data := url.Values{}
		data.Set("license", key)
		data.Set("hwid", md5Hash(name+usr.Username))
		u, _ := url.ParseRequestURI(licenseServer)
		urlStr := fmt.Sprintf("%v", u)
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(r)
		if err != nil {
			fmt.Println("Unable to connect to license server.")
			os.Exit(0)
		} else {
			defer resp.Body.Close()
			resp_body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				if string(resp_body) != "0" {
					if string(resp_body) == "1" {
						fmt.Println("License is Expired.")
						os.Exit(0)
					} else {
						fmt.Println("Unable to verify to license server.")
						os.Exit(0)
					}
				}
			}
		}

	}
}

func main() {
	licenseCheck()
	for {

	}
}
