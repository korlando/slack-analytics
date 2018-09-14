package slackanalytics

import (
  "sort"
  "strings"
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
