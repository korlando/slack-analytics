package slackanalytics

import (
  "fmt"
  "sort"
  "strconv"
  "strings"
)

var (
  commonWords = []string{"i", "you", "he", "she", "it", "we", "they", "me", "him", "her", "us", "them", "what", "who", "whom", "this", "that", "these", "those", "the", "to", "a", "is", "of", "and", "in", "on", "for", "not", "like", "have", "my"}
)

type WordStats struct {
  TotalWords   int
  AvgLength    float64
  WordCountMap map[string]int
}

type WordCount struct {
  Word  string
  Count int
}

type sortByCount []WordCount

func (s sortByCount) Len() int {
  return len(s)
}

func (s sortByCount) Less(i, j int) bool {
  return s[i].Count > s[j].Count
}

func (s sortByCount) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}

// GetWordStats takes in a slice of messages and calculates
// the total # of words, avg word length, and frequency counts.
func GetWordStats(messages []Message) (ws WordStats) {
  ws = WordStats{0, 0, make(map[string]int)}
  totalLength := 0
  for _, m := range messages {
    if m.Text == "" {
      continue
    }
    words := strings.Fields(m.Text)
    for _, w := range words {
      if w == "" {
        continue
      }
      ws.TotalWords += 1
      totalLength += len(w)
      count, ok := ws.WordCountMap[w]
      if !ok {
        ws.WordCountMap[w] = 1
      } else {
        ws.WordCountMap[w] = count + 1
      }
    }
  }
  ws.AvgLength = float64(totalLength) / float64(ws.TotalWords)
  return
}

// GetSortedWords takes in word stats and returns
// sorted word counts by frequency descending.
func GetSortedWords(ws WordStats) (wordCounts []WordCount) {
  wordCounts = make([]WordCount, len(ws.WordCountMap))
  i := 0
  for word, count := range ws.WordCountMap {
    wordCounts[i] = WordCount{word, count}
    i += 1
  }
  sort.Sort(sortByCount(wordCounts))
  return
}

// GetTopWords takes in a slice of words/frequencies (sorted by frequency) and
// returns the top amount of words of highest frequency, skipping common ones
// if includeCommon is false
func GetTopWords(wordCounts []WordCount, amount int, includeCommon bool) (topWordCounts []WordCount) {
  topWordCounts = make([]WordCount, amount)
  i := 0
  j := 0
  for j < amount && i < len(wordCounts) {
    wc := wordCounts[i]
    i += 1
    w := wc.Word
    isCommon := false
    if !includeCommon {
      for _, c := range commonWords {
        if c == strings.ToLower(w) {
          isCommon = true
          continue
        }
      }
    }
    if isCommon {
      continue
    }
    topWordCounts[j] = wc
    j += 1
  }
  return
}

func GetAndPrintStats(messages []Message) {
  ws := GetWordStats(messages)
  wordCounts := GetSortedWords(ws)
  top20 := GetTopWords(wordCounts, 20, false)
  fmt.Println("Total words: " + strconv.Itoa(ws.TotalWords))
  fmt.Println("Avg word length: " + strconv.FormatFloat(ws.AvgLength, 'f', 3, 64))
  for _, wc := range top20 {
    fmt.Println(wc.Word + " " + strconv.Itoa(wc.Count))
  }
}
