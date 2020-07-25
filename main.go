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
	from, _ := todoist.NewClient("", os.Getenv("TODOIST_TOKEN"), "", "", nil)
	dest := rtm.NewClient(os.Getenv("RTM_TOKEN"), os.Getenv("RTM_API_KEY"))

	if err := from.FullSync(context.Background(), []todoist.Command{}); err != nil {
		log.Fatal(err)
	}

	proj := from.Project.FindOneByName("Alexa To-do List")
	items := from.Item.FindByProjectIDs([]todoist.ID{proj.ID})

	for _, v := range items {
		if v.IsChecked() {
			continue
		}

		if err := add(dest, v.Content); err != nil {
			log.Fatal(err)
		}
		if err := from.Item.Complete(v.ID, todoist.Time{Time: time.Now()}, true); err != nil {
			log.Fatal(err)
		}
		if err := from.Commit(context.Background()); err != nil {
			log.Fatal(err)
		}
	}
}

func add(rtm *rtm.Client, task string) error {
	timeline, err := rtm.Timelines.Create()
	if err != nil {
		return err
	}

	_, _, err = rtm.Tasks.Add(timeline, task)
	if err != nil {
		return err
	}

	return nil
}
