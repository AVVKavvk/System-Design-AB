package api

import (
	"github.com/AVVKavvk/bloom_filter/models"
	"github.com/AVVKavvk/bloom_filter/service"
	"github.com/labstack/echo/v4"
)

// AddWord godoc
// @Summary      Add a new word to the Bloom Filter
// @Description  Stores a word in the bloom filter for future membership checks
// @Tags         BloomFilter
// @Accept       json
// @Produce      json
// @Param        word  body      models.Word  true  "Word to add"
// @Success      201   {object}  models.ResponseAddWord
// @Failure      400   {object}  map[string]string "Invalid request body"
// @Router       /words [post]
func AddWord(ctx echo.Context) error {
	var word models.Word

	if err := ctx.Bind(&word); err != nil {
		return err
	}
	result, err := service.AddWordService(&word)
	if err != nil {
		return err
	}
	return ctx.JSON(201, result)
}

// CheckWeatherWordIsExist godoc
// @Summary      Check if a word exists
// @Description  Checks the Bloom Filter for word membership. Note: may return false positives.
// @Tags         BloomFilter
// @Accept       json
// @Produce      json
// @Param        word  body      models.Word  true  "Word to check"
// @Success      200   {object}  map[string]interface{} "Returns a boolean existence check"
// @Failure      400   {object}  map[string]string "Invalid request body"
// @Router       /words/check [post]
func CheckWeatherWordIsExist(ctx echo.Context) error {
	var word models.Word

	if err := ctx.Bind(&word); err != nil {
		return err
	}

	response, err := service.CheckWeatherWordIsExistService(&word)
	if err != nil {
		return err
	}
	return ctx.JSON(200, response)
}
