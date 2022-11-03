package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kangata/db-test-go/database"
	"github.com/kangata/db-test-go/helpers"
)

type Cache struct {
	Name      string `gorm:"type:varchar(100);primayKey"`
	Value     string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	duration, _ := strconv.Atoi(helpers.Env("DURATION", "30"))
	delay, _ := strconv.Atoi(helpers.Env("DELAY", "10"))
	count, _ := strconv.Atoi(helpers.Env("COUNT", "5"))
	lCount := 0

	db := database.New()

	db.AutoMigrate(Cache{})

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	lastRow := Cache{}

	r := db.Order("created_at DESC").First(&lastRow)

	stopTime := time.Now().Add(time.Second * time.Duration(duration))

	fmt.Printf("process start: %v\n", time.Now())
	fmt.Printf("duration: %d second(s)\n", duration)
	fmt.Printf("delay: %d second(s)\n", delay)
	fmt.Printf("count: %d per second\n", count)

	for delay > 0 {
		if time.Since(stopTime) >= 0 {
			fmt.Printf("process end: %v\n", time.Now())

			rowCount := 0

			if r.RowsAffected < 1 {
				db.Raw("SELECT COUNT(*) FROM caches").Scan(&rowCount)
			} else {
				db.Raw("SELECT COUNT(*) FROM caches WHERE created_at > ?", lastRow.CreatedAt).Scan(&rowCount)
			}

			fmt.Printf("rows inserted: %d\n", rowCount)
			fmt.Printf("looped count: %d\n", lCount)

			os.Exit(0)
		}

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
			}

			count -= 1
		}

		delay, _ = strconv.Atoi(helpers.Env("DELAY", "10"))
		count, _ = strconv.Atoi(helpers.Env("COUNT", "5"))
		lCount += 1
	}
}
