package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/BarTar213/movies-service/models"
)

const (
	doNothingStatement = "DO NOTHING"
)

func (p *Postgres) AddMovie(movie *models.TmdbMovie) error {
	_, insertErr := p.db.Model(movie).
		Relation(genres).
		Relation("Countries").
		Relation("Companies").
		Relation("Languages").
		OnConflict("(id) DO UPDATE").
		Set("budget=?budget").
		Set("poster_path=?poster_path").
		Set("backdrop_path=?backdrop_path").
		Set("revenue=?revenue").
		Set("runtime=?runtime").
		Set("vote_average=?vote_average").
		Insert()
	if insertErr != nil {
		return insertErr
	}

	wg := sync.WaitGroup{}
	errs := sync.Map{}

	wg.Add(1)
	go p.addGenres(movie, &wg, &errs)

	wg.Add(1)
	go p.addCountries(movie, &wg, &errs)

	wg.Add(1)
	go p.addCompanies(movie, &wg, &errs)

	wg.Add(1)
	go p.addLanguages(movie, &wg, &errs)

	wg.Wait()
	errMsg := ""
	errs.Range(func(key, value interface{}) bool {
		errMsg += fmt.Sprintf("%v: %v\n", key, value)
		return false
	})
	if len(errMsg) == 0 {
		return nil
	}
	return errors.New(errMsg)
}

func (p *Postgres) addGenres(movie *models.TmdbMovie, wg *sync.WaitGroup, errs *sync.Map) {
	defer wg.Done()
	if len(movie.Genres) == 0 {
		return
	}

	genres := make([]*models.MovieGenre, 0, len(movie.Genres))
	for _, genre := range movie.Genres {
		genres = append(genres, &models.MovieGenre{
			MovieId: movie.Id,
			GenreId: genre.Id,
		})
	}
	_, err := p.db.Model(&genres).OnConflict(doNothingStatement).Insert()
	if err != nil {
		errs.Store("genres", err)

		err = nil
		_, err := p.db.Model(&(movie.Genres)).OnConflict(doNothingStatement).Insert()
		if err != nil{
			fmt.Printf("Additional genres: %s", err)
		}
	}
}

func (p *Postgres) addCountries(movie *models.TmdbMovie, wg *sync.WaitGroup, errs *sync.Map) {
	defer wg.Done()
	if len(movie.Countries) == 0 {
		return
	}

	countries := make([]*models.MovieCountry, 0, len(movie.Countries))
	for _, country := range movie.Countries {
		countries = append(countries, &models.MovieCountry{
			MovieId:     movie.Id,
			CountryCode: country.Code,
		})
	}
	_, err := p.db.Model(&countries).OnConflict(doNothingStatement).Insert()
	if err != nil {
		fmt.Printf("countries: %v\n", movie.Countries)
		errs.Store("countries", err)

		err = nil
		_, err := p.db.Model(&(movie.Countries)).OnConflict(doNothingStatement).Insert()
		if err != nil{
			fmt.Printf("Additional countires: %s", err)
		}
	}
}

func (p *Postgres) addCompanies(movie *models.TmdbMovie, wg *sync.WaitGroup, errs *sync.Map) {
	defer wg.Done()
	if len(movie.Companies) == 0 {
		return
	}

	companies := make([]*models.MovieCompany, 0, len(movie.Companies))
	for _, company := range movie.Companies {
		companies = append(companies, &models.MovieCompany{
			MovieId:   movie.Id,
			CompanyId: company.Id,
		})
	}
	_, err := p.db.Model(&companies).OnConflict(doNothingStatement).Insert()
	if err != nil {
		errs.Store("companies", err)

		err = nil
		_, err := p.db.Model(&(movie.Companies)).OnConflict(doNothingStatement).Insert()
		if err != nil{
			fmt.Printf("Additional companies: %s", err)
		}
	}
}

func (p *Postgres) addLanguages(movie *models.TmdbMovie, wg *sync.WaitGroup, errs *sync.Map) {
	defer wg.Done()
	if len(movie.Languages) == 0 {
		return
	}

	languages := make([]*models.MovieLanguage, 0, len(movie.Languages))
	for _, language := range movie.Languages {
		languages = append(languages, &models.MovieLanguage{
			MovieId: movie.Id,
			IsoCode: language.IsoCode,
		})
	}
	_, err := p.db.Model(&languages).OnConflict(doNothingStatement).Insert()
	if err != nil {
		errs.Store("languages", err)

		err = nil
		_, err := p.db.Model(&(movie.Languages)).OnConflict(doNothingStatement).Insert()
		if err != nil{
			fmt.Printf("Additional languages: %s", err)
		}
	}
}
