package service

import (
	"context"
	"sort"
	"strings"
	"unicode/utf8"
)

type SearchService interface {
	SearchAllPhrasesContainingFirstChar(c context.Context, text string) ([][]string, error)
	SearchAllPhrasesInInput(c context.Context, text string) ([]string, error)
}

type searchService struct {
	characters, phrases                  map[string]bool
	startWithSortedHan, containSortedHan map[string][]string
	sortedArray                          []string
}

func NewSearchService(allChars []string, allPhrases []string) SearchService {
	// chars & phrases set
	characters := map[string]bool{}
	phrases := map[string]bool{}
	for _, char := range allChars {
		characters[char] = true
	}
	for _, phrase := range allPhrases {
		phrases[phrase] = true
	}

	// generate startwith & contain map
	startWithHan := map[string]map[string]bool{}
	containHan := map[string]map[string]bool{}
	for _, char := range allChars {
		startWithHan[char] = map[string]bool{char: true}
		containHan[char] = map[string]bool{char: true}
	}
	for _, phrase := range allPhrases {
		// start with
		first := string([]rune(phrase)[0])
		if _, ok := characters[first]; ok {
			startWithHan[first][phrase] = true
		}
		// contain
		for _, char := range phrase {
			if _, ok := characters[string(char)]; ok {
				containHan[string(char)][phrase] = true
			}
		}
	}

	// sort
	startWithSortedHan := map[string][]string{}
	containSortedHan := map[string][]string{}
	for key, val := range startWithHan {
		newVals := make([]string, len(val))
		i := 0
		for k, _ := range val {
			newVals[i] = k
			i++
		}
		sort.Strings(newVals)
		startWithSortedHan[key] = newVals
	}
	for key, val := range containHan {
		newVals := make([]string, len(val))
		i := 0
		for k, _ := range val {
			newVals[i] = k
			i++
		}
		sort.Strings(newVals)
		containSortedHan[key] = newVals
	}
	return &searchService{
		characters:         characters,
		phrases:            phrases,
		startWithSortedHan: startWithSortedHan,
		containSortedHan:   containSortedHan,
	}
}

func (ps *searchService) SearchAllPhrasesContainingFirstChar(c context.Context, text string) ([][]string, error) {
	first := string([]rune(text)[0])
	if _, ok := ps.characters[first]; !ok {
		return make([][]string, 0), nil
	}
	fullStartWith := make([]string, 0)
	halfStartWith := make([]string, 0)
	rest := make([]string, 0)
	if list, ok := ps.containSortedHan[first]; ok {
		for _, han := range list {
			if strings.HasPrefix(han, first) {
				// phrase starts with first & text starts with phrase
				if strings.HasPrefix(text, han) {
					fullStartWith = append(fullStartWith, han)
				} else { // phrase starts with first but text does not start with phrase
					halfStartWith = append(halfStartWith, han)
				}
			} else {
				rest = append(rest, han)
			}
		}
	}
	result := make([][]string, 3)
	result[0] = fullStartWith
	result[1] = halfStartWith
	result[2] = rest
	return result, nil
}

func (ps *searchService) SearchAllPhrasesInInput(c context.Context, text string) ([]string, error) {
	hans := make(map[string]bool, 0) // processed han
	sameLenPhrases := make(map[int]map[string]int)
	for _, char := range text {
		_, ok1 := ps.characters[string(char)]
		_, ok2 := hans[string(char)]
		if !ok1 || ok2 {
			continue
		}
		// mark as processed
		hans[string(char)] = true

		// get all phrase starting with each char
		list, _ := ps.startWithSortedHan[string(char)]
		for _, phrase := range list {
			// check phrase in text only
			if strings.Contains(text, phrase) {
				lenKey := utf8.RuneCountInString(phrase)
				if lenKeyPhrases, found := sameLenPhrases[lenKey]; found { // len already in sameLenPhrases
					if _, existed := lenKeyPhrases[phrase]; !existed { // phrase already in the map
						sameLenPhrases[lenKey][phrase] = len(sameLenPhrases[lenKey])
					}
				} else { // len not in sameLenPhrases
					sameLenPhrases[lenKey] = map[string]int{phrase: 0}
				}
			}
		}
	}

	// sort keys
	keys := make([]int, len(sameLenPhrases))
	i := 0
	for k, _ := range sameLenPhrases {
		keys[i] = k
		i++
	}
	sort.Ints(keys)

	// loop through keys
	result := make([]string, 0)
	for _, k := range keys {
		// sort values & get keys in order
		pVals := make([]int, len(sameLenPhrases[k]))
		inversedMap := make(map[int]string, len(sameLenPhrases[k]))
		it := 0
		for p, pVal := range sameLenPhrases[k] {
			pVals[it] = pVal
			inversedMap[pVal] = p
			it++
		}
		sort.Ints(pVals)

		// add to result
		for _, pVal := range pVals {
			result = append(result, inversedMap[pVal])
		}
	}

	return result, nil
}
