package main

import (
	"log"
	"os"
	"strings"

	"github.com/zMrKrabz/hidemyemail"
)

func main() {
	// f, err := os.OpenFile("cookies.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	// b, err := io.ReadAll(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))

	f, err := os.ReadFile("cookies.txt")
	if err != nil {
		// log.Fatal(err)
		log.Fatalf("unable to read file: %v", err)
	}

	cookies := string(f)
	hme := hidemyemail.HideMyEmail{
		Cookies: cookies,
		Label:   "Test",
	}

	emails, err := hme.List()
	if err != nil {
		log.Fatal(err)
	}
	emailsString := strings.Join(emails, "\n")
	if err := os.WriteFile("emails.txt", []byte(emailsString), 0644); err != nil {
		log.Fatal(err)
	}

	// for {
	// 	gen, err := hme.Generate()
	// 	if err != nil {
	// 		log.Println(err)
	// 		continue
	// 	}

	// 	reserve, err := hme.Reserve(gen.Result.Hme)
	// 	if err != nil {
	// 		log.Printf("error when reserving: %v\n", err)
	// 		time.Sleep(1 * time.Minute)
	// 	} else {
	// 		log.Printf("Successfully genned email: %v\n", reserve.Result.Hme.Hme)
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }
}
