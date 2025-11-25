// cmd/server/main.go
package main

import (
	"flag"
	"fmt"
	"log"

	"auth-service/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	mod := mustMode()
	application, err := app.New(mod)
	if err != nil {
		log.Fatal()
		fmt.Println("Ошибка сборки приложение")
	}

	application.Run()
}

func mustMode() string {

	mode := flag.String(
		"app-mode",
		"",
		"application launch mode for selecting settings in env",
	)

	flag.Parse()

	if *mode == "" { // ← ПРАВИЛЬНАЯ ПРОВЕРКА!
		fmt.Println("No mode specified, defaulting to 'dev'")
		return "dev"
	}

	return *mode

}
