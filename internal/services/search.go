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
	targets := make([]string, len(parts))
	nameToPart := make(map[string]models.Part)

	for i, p := range parts {
		name := strings.ToLower(p.Name)
		targets[i] = name
		nameToPart[name] = p
	}

	ranks := fuzzy.RankFind(query, targets)

	results := make([]models.SearchResult, 0, len(ranks))
	for _, r := range ranks {
		if part, ok := nameToPart[r.Target]; ok {
			score := calculateScore(query, r.Target, r.Distance)
			if score >= s.minScore {
				results = append(results, models.SearchResult{
					Part:  part,
					Score: score,
				})
			}
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
		if filters.Type != "" && p.Type != filters.Type {
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
