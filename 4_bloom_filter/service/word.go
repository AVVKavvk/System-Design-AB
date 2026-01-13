package service

import (
	"github.com/AVVKavvk/bloom_filter/bloomFilter"
	"github.com/AVVKavvk/bloom_filter/models"
)

func AddWordService(word *models.Word) (*models.ResponseAddWord, error) {
	blf := bloomFilter.GetBloomFilter()
	rowInd, colInd := blf.Add([]byte(word.Word))
	return &models.ResponseAddWord{
		RowIdx: rowInd,
		ColIdx: colInd,
	}, nil
}

func CheckWeatherWordIsExistService(word *models.Word) (*models.ResponseWordProbability, error) {
	blf := bloomFilter.GetBloomFilter()
	isFound, rawInd, colInd := blf.Contains([]byte(word.Word))

	return &models.ResponseWordProbability{
		IsFound: isFound,
		RowIdx:  rawInd,
		ColIdx:  colInd,
	}, nil
}
