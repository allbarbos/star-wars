package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"star-wars/entity"
	"star-wars/importer"
	"star-wars/planet"
	"star-wars/swapi"
	"time"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func openCsv() *os.File {
	csvfile, err := os.Open(os.Getenv("PATH_CSV"))
	if err != nil {
		log.Fatalln("couldn't open the csv file", err)
	}

	return csvfile
}

func readCsv(file *os.File) []entity.Planet {
	r := csv.NewReader(file)
	r.Comma = ';'

	if _, err := r.Read(); err != nil {
		panic(err)
	}

	var planets []entity.Planet

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		newPlanet := entity.Planet{
			Name:    rec[0],
			Climate: rec[1],
			Terrain: rec[2],
		}

		empty := newPlanet.IsEmpty([]string{"Name", "Climate", "Terrain"})

		if empty {
			log.Printf("row error: %s; %s; %s", rec[0], rec[1], rec[2])
			continue
		}

		planets = append(planets, newPlanet)
	}

	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return planets
}

func main() {
	loadEnv()

	csvfile := openCsv()
	planets := readCsv(csvfile)

	s := swapi.New()
	p := planet.NewService(planet.NewRepository()	, s)
	srv := importer.NewImporter(p, s)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	errors := srv.Import(ctx, planets)

	for _, err := range errors {
		log.Print(err)
	}

	fmt.Println("> completed - errors:", len(errors))
}
