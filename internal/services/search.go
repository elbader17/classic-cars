package services

import (
	"sort"
	"strings"

	"github.com/eduardo/classicCarSearch/internal/models"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type SearchService struct {
	minScore int
}

func NewSearchService() *SearchService {
	return &SearchService{
		minScore: 50,
	}
}

func (s *SearchService) FuzzySearch(query string, parts []models.Part) []models.SearchResult {
	if query == "" {
		return allAsResults(parts)
	}

	query = strings.ToLower(strings.TrimSpace(query))
	results := make([]models.SearchResult, 0, len(parts))

	for _, p := range parts {
		searchFields := []string{
			strings.ToLower(p.Name),
			strings.ToLower(p.Brand),
			strings.ToLower(p.Category),
			strings.ToLower(p.Model),
			strings.ToLower(p.Year),
			strings.ToLower(p.Description),
		}

		bestScore := 0
		for _, field := range searchFields {
			if field == "" {
				continue
			}

			ranks := fuzzy.RankFind(query, []string{field})
			if len(ranks) > 0 {
				r := ranks[0]
				score := calculateScore(query, r.Target, r.Distance)
				if score > bestScore {
					bestScore = score
				}
			}

			if strings.Contains(field, query) {
				bestScore = max(bestScore, 60)
			}
		}

		if bestScore >= s.minScore {
			results = append(results, models.SearchResult{
				Part:  p,
				Score: bestScore,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func (s *SearchService) FuzzySearchWithFilters(query string, parts []models.Part, filters models.FilterOptions) []models.SearchResult {
	filtered := applyFilters(parts, filters)
	return s.FuzzySearch(query, filtered)
}

func applyFilters(parts []models.Part, filters models.FilterOptions) []models.Part {
	result := make([]models.Part, 0, len(parts))

	for _, p := range parts {
		if filters.Brand != "" && p.Brand != filters.Brand {
			continue
		}
		if filters.Category != "" && p.Category != filters.Category {
			continue
		}
		if filters.Subcategoria != "" && p.Subcategoria != filters.Subcategoria {
			continue
		}
		result = append(result, p)
	}

	return result
}

func calculateScore(query, target string, distance int) int {
	maxLen := len(target)
	if maxLen == 0 {
		return 0
	}

	score := 100 - (distance * 100 / maxLen)
	if score < 0 {
		score = 0
	}

	if strings.HasPrefix(target, query) {
		score = min(100, score+20)
	}

	if strings.Contains(target, query) {
		score = min(100, score+10)
	}

	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func allAsResults(parts []models.Part) []models.SearchResult {
	results := make([]models.SearchResult, len(parts))
	for i, p := range parts {
		results[i] = models.SearchResult{
			Part:  p,
			Score: 100,
		}
	}
	return results
}
