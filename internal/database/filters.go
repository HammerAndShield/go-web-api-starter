package database

import (
	"regexp"
	"strings"
)

type OrderDirection int

const (
	_ = iota
	Ascending
	Descending
)

const (
	ValidOrderDirectionAsc  = "asc"
	ValidOrderDirectionDesc = "desc"
)

var (
	ValidOrderDirections = []string{
		ValidOrderDirectionAsc,
		ValidOrderDirectionDesc,
	}
)

func DirectionFromString(direction string) OrderDirection {
	switch direction {
	case ValidOrderDirectionAsc:
		return Ascending
	case ValidOrderDirectionDesc:
		return Descending
	default:
		return 0
	}
}

type TsQuery string

// FilterTsQueryOperators removes special operators and numbers used in distance operators
// from the TsQuery, returning a cleaned version of the query.
func (q TsQuery) FilterTsQueryOperators() TsQuery {
	// Remove special operators for ts query functions
	re := regexp.MustCompile(`[&|!<>:*()]`)
	filtered := re.ReplaceAllString(string(q), "")

	// Remove numbers that might be used in distance operators
	re = regexp.MustCompile(`\s+\d+\s+`)
	filtered = re.ReplaceAllString(filtered, " ")

	return TsQuery(strings.TrimSpace(filtered))
}

// AddPrefixMatching modifies the TsQuery by appending a prefix matching operator
// to each word in the query, enabling prefix-based search.
//
// Must be used after FilterTsQueryOperators or they will be removed.
func (q TsQuery) AddPrefixMatching() TsQuery {
	words := strings.Fields(string(q))
	for i, word := range words {
		words[i] = word + ":*"
	}
	search := strings.Join(words, " & ")

	return TsQuery(search)
}
