package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/cron"
	schedule_service "github.com/unendlicherping2/school-schedule-server/src"
)

func main() {
	app := pocketbase.New()

	entryCollection := &models.Collection{
		Name: "schedule",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "hint",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "date",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "schedule",
				Type:     schema.FieldTypeJson,
				Required: true,
				Options: &schema.JsonOptions{
					MaxSize: 2000000,
				},
			},
		),
	}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		if err := app.Dao().SaveCollection(entryCollection); err != nil {
			fmt.Printf("Could not create database: %s")
			return nil
		}

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		scheduler := cron.New()

		// runs every 15 minutes
		scheduler.MustAdd("checkSchedule", "*/1 * * * *", func() {
			date := time.Now()

			if date.Hour() > 17 {
				date = date.AddDate(0, 0, 1)
			}

			formatted := date.Format("2006-01-02")

			schedule, err := schedule_service.GetSchedule(formatted)
			if err != nil {
				return
			}

			collection, err := app.Dao().FindCollectionByNameOrId("schedule")
			if err != nil {
				return
			}

			var old_schedule_dto schedule_service.ScheduleDto
			app.Dao().DB().
				Select("hint", "date", "schedule").
				From("schedule").
				Where(dbx.NewExp("date = {:date}")).
				Limit(1).
				Bind(dbx.Params{
					"date": formatted,
				}).
				One(&old_schedule_dto)

			var old_schedule schedule_service.Schedule
			old_schedule.Hint = old_schedule_dto.Hint
			old_schedule.Date = old_schedule_dto.Date

			var old_schedule_entries schedule_service.ScheduleEntries

			bytes := []byte(old_schedule_dto.Schedule)

			json.Unmarshal(bytes, &old_schedule_entries)

			old_schedule.Schedule = old_schedule_entries

			if reflect.DeepEqual(old_schedule, schedule) {
				return
			}

			if old_schedule.Date != formatted {
				record := models.NewRecord(collection)

				record.Set("hint", schedule.Hint)
				record.Set("date", schedule.Date)
				record.Set("schedule", schedule.Schedule)

				if err := app.Dao().SaveRecord(record); err != nil {
					log.Printf("error saving record: %s\n", err)
					return
				}

				log.Println("Added new schedule entry!")

				return
			}

			record, err := app.Dao().FindFirstRecordByData("schedule", "date", formatted)
			if err != nil {
				log.Printf("an error occured: %s", err)
			}

			record.Set("hint", schedule.Hint)
			record.Set("schedule", schedule.Schedule)

			log.Printf("Updated record")
		})

		scheduler.Start()

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
