package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/models"
)

type Client interface {
	GetCredits(movieId int, credit *models.Credit) (int, error)
	GetTrendingMovies() ([]models.TmdbMovie, int, error)
}

type Tmdb struct {
	HttpClient *http.Client
	ApiKey     string
	BaseUrl    string
}

func NewClient(timeout time.Duration, config *config.Config) Client {
	client := http.DefaultClient
	client.Timeout = timeout

	return &Tmdb{
		HttpClient: client,
		ApiKey:     config.Tmdb.Key,
		BaseUrl:    config.Tmdb.Url,
	}
}

func (c *Tmdb) GetCredits(movieId int, credit *models.Credit) (int, error) {
	url := fmt.Sprintf("%s/movie/%d/credits?api_key=%s", c.BaseUrl, movieId, c.ApiKey)

	request, err := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(credit)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}

func (c *Tmdb) GetTrendingMovies() ([]models.TmdbMovie, int, error) {
	url := fmt.Sprintf("%s/trending/movie/day?api_key=%s", c.BaseUrl, c.ApiKey)

	request, err := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	response := &LatestResponse{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	movies := make([]models.TmdbMovie, 0, len(response.Results))
	for _, id := range response.Results {
		movie := &models.TmdbMovie{}
		status, err := c.GetMovieDetails(id.Id, movie)
		if err != nil || status != http.StatusOK {
			continue
		}
		movies = append(movies, *movie)
	}

	return movies, resp.StatusCode, nil
}

func (c *Tmdb) GetMovieDetails(id int, movie *models.TmdbMovie) (int, error) {
	url := fmt.Sprintf("%s/movie/%d?api_key=%s", c.BaseUrl, id, c.ApiKey)

	request, err := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(movie)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}
