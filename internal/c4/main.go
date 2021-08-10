package main

import (
	"context"
	"os"

	"github.com/krzysztofreczek/go-structurizr/pkg/model"
	"github.com/krzysztofreczek/go-structurizr/pkg/scraper"
	"github.com/krzysztofreczek/go-structurizr/pkg/view"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/service"
)

func main() {
	config := scraper.NewConfiguration(
		"github.com/",
		"cloud.google.com",
	)

	s := scraper.NewScraper(config)

	r, err := scraper.NewRule().
		WithApplyFunc(
			func(name string, groups ...string) model.Info {
				return model.ComponentInfo(name)
			},
		).Build()
	if err != nil {
		panic(err)
	}

	err = s.RegisterRule(r)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	a := service.NewApplication(ctx)

	structure := s.Scrape(a)

	f, err := os.Create("out/output.plantuml")
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = f.Close()
	}()

	v := view.NewView().
		WithTitle("Components").
		Build()
	err = v.RenderStructureTo(structure, f)
	if err != nil {
		panic(err)
	}
}
