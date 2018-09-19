package slackanalytics

import (
  "fmt"
  "sort"
  "strconv"
  "strings"
)

type WordStats struct {
  TotalWords        int
  AvgLength         float64
  AvgWordsPerMsg    float64
  AvgCloutPerMsg    float64
  AvgTonePerMsg     float64
  AvgAnalyticPerMsg float64
  WordCountMap      map[string]int
  CategoryCounts    map[string]int
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
  ws = WordStats{
    TotalWords:        0,
    AvgLength:         0,
    AvgWordsPerMsg:    0,
    AvgCloutPerMsg:    0,
    AvgTonePerMsg:     0,
    AvgAnalyticPerMsg: 0,
    WordCountMap:      make(map[string]int),
    CategoryCounts:    make(map[string]int),
  }
  totalLength := 0
  totalClout := 0
  totalTone := 0
  totalAnalytic := 0
  for _, m := range messages {
    if m.Text == "" {
      continue
    }
    words := MessageToWords(m, true, true)
    totalClout += GetClout(words)
    totalTone += GetTone(words)
    totalAnalytic += GetAnalytic(words)
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
    counts := GetCategoryCounts(words)
    // merge counts
    for k, v := range counts {
      count, ok := ws.CategoryCounts[k]
      if ok {
        ws.CategoryCounts[k] = count + v
      } else {
        ws.CategoryCounts[k] = v
      }
    }
  }
  ws.AvgLength = float64(totalLength) / float64(ws.TotalWords)
  ws.AvgWordsPerMsg = float64(ws.TotalWords) / float64(len(messages))
  ws.AvgCloutPerMsg = float64(totalClout) / float64(ws.TotalWords)
  ws.AvgTonePerMsg = float64(totalTone) / float64(ws.TotalWords)
  ws.AvgAnalyticPerMsg = float64(totalAnalytic) / float64(ws.TotalWords)
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
      for _, c := range CommonWords {
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

// GetClout loosely calculates the clout of a
// slice of words (+1 for we/you and -1 for i)
func GetClout(words []string) (clout int) {
  for _, w := range words {
    wLower := strings.ToLower(w)
    if inList(wLower, IWords) {
      clout -= 1
      continue
    }
    if inList(wLower, YouWords) || inList(wLower, WeWords) {
      clout += 1
    }
  }
  return
}

// GetTone loosely calculates the tone of a slice
// of words (+1 for pos emo and -1 for neg emo)
func GetTone(words []string) (tone int) {
  for _, w := range words {
    wLower := strings.ToLower(w)
    if inList(wLower, PosEmo) {
      tone += 1
      continue
    }
    if inList(wLower, NegEmo) {
      tone -= 1
    }
  }
  return
}

func GetAnalytic(words []string) (analytic int) {
  analytic = 30
  for _, w := range words {
    w = strings.ToLower(w)
    if inList(w, Articles) || inList(w, Prepositions) {
      analytic += 1
      continue
    }
    if inList(w, PersonalPronouns) || inList(w, ImpersonalPronouns) || inList(w, AuxiliaryVerbs) || inList(w, Conjunctions) || inList(w, Adverbs) || inList(w, Negations) {
      analytic -= 1
      continue
    }
  }
  return
}

func GetCategoryCounts(words []string) (counts map[string]int) {
  counts = make(map[string]int)
  for _, w := range words {
    for k, v := range Categories {
      if inList(w, v) {
        count, ok := counts[k]
        if ok {
          counts[k] = count + 1
        } else {
          counts[k] = 1
        }
      }
    }
  }
  return
}

// GetAndPrintStats takes in a slice of messages
// and prints some stats about them
func GetAndPrintStats(messages []Message) {
  ws := GetWordStats(messages)
  wordCounts := GetSortedWords(ws)
  topWords := GetTopWords(wordCounts, 0, false)
  fmt.Println("Total words: " + strconv.Itoa(ws.TotalWords))
  fmt.Println("Avg word length: " + floatStr(ws.AvgLength, 4))
  fmt.Println("Avg words per message: " + floatStr(ws.AvgWordsPerMsg, 4))
  fmt.Println("Avg message clout: " + floatStr(ws.AvgCloutPerMsg, 4))
  fmt.Println("Avg message tone: " + floatStr(ws.AvgTonePerMsg, 4))
  fmt.Println("Avg message analytic: " + floatStr(ws.AvgAnalyticPerMsg, 4))
  fmt.Println("Category counts:")
  fmt.Println(ws.CategoryCounts)
  for _, wc := range topWords {
    fmt.Println(wc.Word + " " + strconv.Itoa(wc.Count))
  }
}

// inList determines whether a word
// is contained in a slice of words
func inList(word string, words []string) bool {
  for _, w := range words {
    if w == word {
      return true
    }
  }
  return false
}

// floatStr converts a float64 to a string
// with decimals worth of precision
func floatStr(f float64, decimals int) string {
  return strconv.FormatFloat(f, 'f', decimals, 64)
}
