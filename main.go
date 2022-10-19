package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Cache struct {
	Name      string `gorm:"type:varchar(100);primayKey"`
	Value     string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	db := database()

	db.AutoMigrate(Cache{})

	delay, _ := strconv.Atoi(env("DELAY", "10"))
	count, _ := strconv.Atoi(env("COUNT", "5"))

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	fmt.Printf("delay: %d\n", delay)
	fmt.Printf("count: %d\n", count)

	for delay > 0 {
		time.Sleep(time.Second)

		delay -= 1

		if delay > 1 {
			continue
		}

		for count > 0 {
			name := uuid.New()

			x := make([]rune, 128)
			y := make([]rune, 128)

			for i := range x {
				x[i] = letters[rand.Intn(len(letters))]
				y[i] = letters[rand.Intn(len(letters))]
			}

			res := db.Exec(
				"INSERT INTO caches (name, value, created_at, updated_at) VALUES (?, ?, ?, ?)",
				name.String(),
				string([]byte(fmt.Sprintf(`{"x": "%s","y": "%s"}`, string(x), string(y)))),
				time.Now(),
				time.Now(),
			)

			if res.Error != nil {
				fmt.Printf("Failed insert to database: %v\n", res.Error)
			} else {
				fmt.Println("Success insert to database")
			}

			count -= 1
		}

		delay, _ = strconv.Atoi(env("DELAY", "10"))
		count, _ = strconv.Atoi(env("COUNT", "5"))
	}
}

func env(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

func database() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env("DB_USERNAME", "root"),
		env("DB_PASSWORD", "root"),
		env("DB_HOST", "127.0.0.1"),
		env("DB_PORT", "3306"),
		env("DB_DATABASE", "test__database"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database: %v", err))
	}

	return db
}
