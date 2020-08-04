package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"star-wars/db"
	"star-wars/entity"
	"star-wars/importer"
	"star-wars/planet"
	"star-wars/swapi"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func openCsv() *os.File {
	csvfile, err := os.Open("seed.csv")
	if err != nil {
		log.Fatalln("couldn't open the csv file", err)
	}

	return csvfile
}

func main() {
	loadEnv()

	csvfile := openCsv()
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	r.Comma = ';'

	if _, err := r.Read(); err != nil {
    panic(err)
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		newPlanet := entity.Planet{
			Name: record[0],
			Climate: record[1],
			Terrain: record[2],
		}

		empty := newPlanet.IsEmpty([]string{"Name", "Climate", "Terrain"})

		if empty {
			log.Printf("error: %s; %s; %s", record[0], record[1], record[2])
			continue
		}

		d := db.New()
		r := planet.NewRepository(d)
		p := planet.NewService(r)
		s := swapi.New()
		srv := importer.NewImporter(p, s)

		errchan := make(chan entity.Planet)

		go srv.Process(newPlanet, errchan)

		for p := range errchan {
			log.Println(p)
		}

	}
}
