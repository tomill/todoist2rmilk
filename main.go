package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/kobtea/go-todoist/todoist"
	"github.com/parroty/rtmgo/rtm"
)

func main() {
	todoi, _ := todoist.NewClient("", os.Getenv("TODOIST_TOKEN"), "", "", nil)
	rmilk := rtm.NewClient(os.Getenv("RTM_TOKEN"), os.Getenv("RTM_API_KEY"))

	if err := todoi.FullSync(context.Background(), []todoist.Command{}); err != nil {
		log.Fatal(err)
	}

	timeline, err := rmilk.Timelines.Create()
	if err != nil {
		log.Fatal(err)
	}

	proj := todoi.Project.FindOneByName("Alexa To-do List")
	items := todoi.Item.FindByProjectIDs([]todoist.ID{proj.ID})
	for _, v := range items {
		if v.IsChecked() {
			continue
		}

		if _, _, err = rmilk.Tasks.Add(timeline, v.Content); err != nil {
			log.Fatal(err)
		}
		if err := todoi.Item.Complete(v.ID, todoist.Time{Time: time.Now().UTC()}, true); err != nil {
			log.Fatal(err)
		}
	}

	if err := todoi.Commit(context.Background()); err != nil {
		log.Fatal(err)
	}
}
