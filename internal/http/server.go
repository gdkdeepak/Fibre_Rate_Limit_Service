package http

import "log"

func Start(app interface{ Listen(string) error }) {
	if err := app.Listen(":8080"); err != nil {
		log.Println("server error:", err)
	}
}
