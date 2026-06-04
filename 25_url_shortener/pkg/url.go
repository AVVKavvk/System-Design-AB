package pkg

import (
	"context"
	"fmt"

	"github.com/AVVKavvk/url-shortener/constants"
	"github.com/AVVKavvk/url-shortener/db"
	"github.com/AVVKavvk/url-shortener/models"
	redisClient "github.com/AVVKavvk/url-shortener/redis"
	"github.com/AVVKavvk/url-shortener/shortener"
	"github.com/jackc/pgx/v5"
)

type UrlPkg struct{}

func (up *UrlPkg) CreateLongUrl(req *models.UrlShortenerRequest) (*models.UrlShortenerResponse, error) {

	urlDb, err := db.GetDb()
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	// Keep generating until we find a unique short_url
	var shortUrl string
	for {
		shortUrl, err = shortener.ShortenURL(req.LongUrl, constants.SHORT_URL_LENGTH)
		if err != nil {
			fmt.Println("error:", err)
			return nil, err
		}

		var existing string
		query := "SELECT short_url FROM urls WHERE short_url=$1"
		err = urlDb.QueryRow(context.Background(), query, shortUrl).Scan(&existing)

		if err == pgx.ErrNoRows {
			// Unique — no collision, break out
			break
		}
		if err != nil {
			fmt.Println("error:", err)
			// Real DB error
			return nil, err
		}
		// err == nil means a row was found → collision, loop again
	}

	query := "INSERT INTO urls (short_url, long_url, user_id) VALUES ($1, $2, $3)"
	_, err = urlDb.Exec(context.Background(), query, shortUrl, req.LongUrl, req.UserId)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	ur := redisClient.UrlRedis{}

	err = ur.AddUrlMapping(req.LongUrl, shortUrl, constants.REDIS_TTL_SECONDS)
	if err != nil {
		fmt.Println("error:", err)
		fmt.Println(err)
	}

	return &models.UrlShortenerResponse{
		ShortUrl: shortUrl,
		LongUrl:  req.LongUrl,
	}, nil

}

func (up *UrlPkg) GetLongUrl(shortUrl string) (string, error) {
	ur := redisClient.UrlRedis{}
	longUrl, err := ur.GetUrlMapping(shortUrl)

	if err != nil {
		return "", err
	}

	if longUrl != "" {
		// Cache hit
		return longUrl, nil
	}

	// Cache miss — fetch from DB
	longUrl, err = getLongUrl(shortUrl)
	if err != nil {
		return "", err
	}
	if longUrl == "" {
		return "", nil
	}

	// Populate cache
	if err = ur.AddUrlMapping(longUrl, shortUrl, constants.REDIS_TTL_SECONDS); err != nil {
		fmt.Println("redis set error:", err)
	}

	return longUrl, nil
}
func getLongUrl(shortUrl string) (string, error) {
	urlDb, err := db.GetDb()
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	var existing string
	query := "SELECT long_url FROM urls WHERE short_url=$1"
	err = urlDb.QueryRow(context.Background(), query, shortUrl).Scan(&existing)

	if err == pgx.ErrNoRows {
		fmt.Println("error:", err)
		return "", nil
	}
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	return existing, nil
}
