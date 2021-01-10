package slackanalytics

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SlackStats struct {
	AllStats     *WordStats
	UserStats    map[string]*WordStats
	ChannelStats map[string]*WordStats
}

type SlackMessageStats struct {
	Time         int
	AllStats     *MessageStats
	UserStats    map[string]*MessageStats
	DailyStats   map[string]*MessageStats
	MonthlyStats map[string]*MessageStats
}

type MessageStats struct {
	NumMessages       int
	NumWords          int
	NumEmojis         int
	TotalTextLength   int
	AvgWordLength     float64
	AvgWordsPerMsg    float64
	AvgEmojisPerMsg   float64
	AvgCloutPerMsg    float64
	AvgTonePerMsg     float64
	AvgAnalyticPerMsg float64
	WordCountMap      map[string]int
	EmojiCountMap     map[string]int
	CategoryCounts    map[string]int
}

type WordStats struct {
	TotalTextLength   int
	TotalWords        int
	TotalMessages     int
	AvgWordLength     float64
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

// GetSlackStats takes in a slice of users and channels and calculates the
// total # of words, avg word length, frequency counts, and sentiment analysis.
func GetSlackStats(users []*User, channels []*Channel) (ss SlackStats) {
	SortCategories()
	ss = SlackStats{
		AllStats:     newWordStats(),
		UserStats:    make(map[string]*WordStats),
		ChannelStats: make(map[string]*WordStats),
	}
	// init user stats
	for _, u := range users {
		ss.UserStats[u.Id] = newWordStats()
	}
	// init channel stats
	for _, c := range channels {
		ss.ChannelStats[c.Id] = newWordStats()
	}
	var words []string
	var clout float64
	var tone float64
	var analytic float64
	for _, c := range channels {
		for _, m := range c.Messages {
			if m.Text == "" {
				continue
			}
			words = MessageToWords(m, true, true)
			clout = float64(GetClout(words))
			tone = float64(GetTone(words))
			analytic = float64(GetAnalytic(words))
			ss.AllStats.TotalTextLength += len(m.Text)
			ss.AllStats.TotalMessages += 1
			ss.AllStats.AvgCloutPerMsg += clout
			ss.AllStats.AvgTonePerMsg += tone
			ss.AllStats.AvgAnalyticPerMsg += analytic
			userStats, userOk := ss.UserStats[m.User]
			if userOk {
				userStats.TotalTextLength += len(m.Text)
				userStats.TotalMessages += 1
				userStats.AvgCloutPerMsg += clout
				userStats.AvgTonePerMsg += tone
				userStats.AvgAnalyticPerMsg += analytic
			}
			channelStats, channelOk := ss.ChannelStats[c.Name]
			if channelOk {
				channelStats.TotalTextLength += len(m.Text)
				channelStats.TotalMessages += 1
				channelStats.AvgCloutPerMsg += clout
				channelStats.AvgTonePerMsg += tone
				channelStats.AvgAnalyticPerMsg += analytic
			}
			for _, w := range words {
				if w == "" {
					continue
				}
				l := float64(len(w))
				ss.AllStats.TotalWords += 1
				ss.AllStats.AvgWordLength += l
				updateWordCountMap(w, &ss.AllStats.WordCountMap)
				if userOk {
					userStats.TotalWords += 1
					userStats.AvgWordLength += l
					updateWordCountMap(w, &userStats.WordCountMap)
				}
				if channelOk {
					channelStats.TotalWords += 1
					channelStats.AvgWordLength += l
					updateWordCountMap(w, &channelStats.WordCountMap)
				}
			}
		}
	}
	wordCategoriesCache := make(map[string][]string)
	populateCategoryCounts(ss.AllStats, &wordCategoriesCache)
	setAverages(ss.AllStats)
	for _, ws := range ss.UserStats {
		if ws.TotalWords == 0 {
			continue
		}
		populateCategoryCounts(ws, &wordCategoriesCache)
		setAverages(ws)
	}
	for _, ws := range ss.ChannelStats {
		if ws.TotalWords == 0 {
			continue
		}
		populateCategoryCounts(ws, &wordCategoriesCache)
		setAverages(ws)
	}
	return
}

// GetSortedWords takes in word stats and returns
// sorted word counts by frequency descending.
func GetSortedWords(ws *WordStats) (wordCounts []WordCount) {
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

// GetAnalytic loosely calculates the analytical thinking
// of a slice of words (+1 for article, prep and -1 for
// ppron, ipron, auxverb, conj, adverb, negation)
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

func GetCategories(word string) (categories []string) {
	categories = []string{}
	for cat, catWords := range Categories {
		if inListBinary(word, catWords) {
			categories = append(categories, cat)
		}
	}
	return
}

// GetAndPrintStats takes in a slice of users and
// slice of channels, and prints some stats about them
func GetAndPrintStats(users []*User, channels []*Channel) {
	ss := GetSlackStats(users, channels)
	wordCounts := GetSortedWords(ss.AllStats)
	topWords := GetTopWords(wordCounts, 0, false)
	printStats(ss.AllStats)
	fmt.Println("Category counts:")
	fmt.Println(ss.AllStats.CategoryCounts)
	fmt.Println()
	for _, wc := range topWords {
		fmt.Println(wc.Word + " " + strconv.Itoa(wc.Count))
	}
	for _, u := range users {
		if u.Deleted {
			continue
		}
		ws, ok := ss.UserStats[u.Id]
		if !ok || ws.TotalWords == 0 {
			continue
		}
		name := u.Profile.RealName
		if u.Profile.DisplayName != "" {
			name = u.Profile.DisplayName
		}
		fmt.Println(name + "\n")
		printStats(ws)
		fmt.Println()
	}
}

// AnalyzeMessages uses messages to get statistics on a per-user basis,
// per-day basis, and per-month basis. Overall stats are also included.
func AnalyzeMessages(messages []Message) (s SlackMessageStats) {
	SortCategories()
	s = SlackMessageStats{
		Time:         int(time.Now().Unix()),
		AllStats:     newMessageStats(),
		UserStats:    make(map[string]*MessageStats),
		DailyStats:   make(map[string]*MessageStats),
		MonthlyStats: make(map[string]*MessageStats),
	}

	var words []string
	var emojis []string
	var clout float64
	var tone float64
	var analytic float64

	for _, m := range messages {
		if m.Text == "" {
			continue
		}
		var t int64
		t64, err := strconv.ParseFloat(m.TimeStamp, 64)
		if err != nil {
			t = time.Now().Unix()
		} else {
			t = int64(t64)
		}
		tm := time.Unix(t, 0)
		d := tm.Format("2006-01-02")
		mo := tm.Format("2006-01")
		u := m.User

		words, emojis = ParseWords(m, false)
		clout = float64(GetClout(words))
		tone = float64(GetTone(words))
		analytic = float64(GetAnalytic(words))

		s.AllStats.TotalTextLength += len(m.Text)
		s.AllStats.NumMessages += 1
		s.AllStats.NumWords += len(words)
		s.AllStats.NumEmojis += len(emojis)
		s.AllStats.AvgCloutPerMsg += clout
		s.AllStats.AvgTonePerMsg += tone
		s.AllStats.AvgAnalyticPerMsg += analytic

		userStats, userOk := s.UserStats[u]
		if !userOk {
			userStats = newMessageStats()
			s.UserStats[u] = userStats
		}
		dailyStats, dailyOk := s.DailyStats[d]
		if !dailyOk {
			dailyStats = newMessageStats()
			s.DailyStats[d] = dailyStats
		}
		monthlyStats, monthlyOk := s.MonthlyStats[mo]
		if !monthlyOk {
			monthlyStats = newMessageStats()
			s.MonthlyStats[mo] = monthlyStats
		}

		userStats.TotalTextLength += len(m.Text)
		userStats.NumMessages += 1
		userStats.NumWords += len(words)
		userStats.NumEmojis += len(emojis)
		userStats.AvgCloutPerMsg += clout
		userStats.AvgTonePerMsg += tone
		userStats.AvgAnalyticPerMsg += analytic

		dailyStats.TotalTextLength += len(m.Text)
		dailyStats.NumMessages += 1
		dailyStats.NumWords += len(words)
		dailyStats.NumEmojis += len(emojis)
		dailyStats.AvgCloutPerMsg += clout
		dailyStats.AvgTonePerMsg += tone
		dailyStats.AvgAnalyticPerMsg += analytic

		monthlyStats.TotalTextLength += len(m.Text)
		monthlyStats.NumMessages += 1
		monthlyStats.NumWords += len(words)
		monthlyStats.NumEmojis += len(emojis)
		monthlyStats.AvgCloutPerMsg += clout
		monthlyStats.AvgTonePerMsg += tone
		monthlyStats.AvgAnalyticPerMsg += analytic

		for _, w := range words {
			if w == "" {
				continue
			}
			l := float64(len(w))
			s.AllStats.AvgWordLength += l
			updateWordCountMap(w, &s.AllStats.WordCountMap)

			userStats.AvgWordLength += l
			updateWordCountMap(w, &userStats.WordCountMap)

			dailyStats.AvgWordLength += l
			updateWordCountMap(w, &dailyStats.WordCountMap)

			monthlyStats.AvgWordLength += l
			updateWordCountMap(w, &monthlyStats.WordCountMap)
		}

		for _, e := range emojis {
			if e == "" {
				continue
			}
			updateWordCountMap(e, &s.AllStats.EmojiCountMap)
			updateWordCountMap(e, &userStats.EmojiCountMap)
			updateWordCountMap(e, &dailyStats.EmojiCountMap)
			updateWordCountMap(e, &monthlyStats.EmojiCountMap)
		}
	}

	wordCategoriesCache := make(map[string][]string)
	populateMsgCategoryCounts(s.AllStats, &wordCategoriesCache)
	setStatAverages(s.AllStats)
	for _, ms := range s.UserStats {
		if ms.NumWords == 0 {
			continue
		}
		populateMsgCategoryCounts(ms, &wordCategoriesCache)
		setStatAverages(ms)
	}
	for _, ms := range s.DailyStats {
		if ms.NumWords == 0 {
			continue
		}
		populateMsgCategoryCounts(ms, &wordCategoriesCache)
		setStatAverages(ms)
	}
	for _, ms := range s.MonthlyStats {
		if ms.NumWords == 0 {
			continue
		}
		populateMsgCategoryCounts(ms, &wordCategoriesCache)
		setStatAverages(ms)
	}
	return
}

func ExportMessageAnalysis(messages []Message) {
	s := AnalyzeMessages(messages)
	ExportSlackMessageStats(s)
}

func printStats(ws *WordStats) {
	fmt.Println("Total text length: " + strconv.Itoa(ws.TotalTextLength))
	fmt.Println("Total words: " + strconv.Itoa(ws.TotalWords))
	fmt.Println("Total messages: " + strconv.Itoa(ws.TotalMessages))
	fmt.Println("Avg word length: " + floatStr(ws.AvgWordLength, 4))
	fmt.Println("Avg words per message: " + floatStr(ws.AvgWordsPerMsg, 4))
	fmt.Println("Avg message clout: " + floatStr(ws.AvgCloutPerMsg, 4))
	fmt.Println("Avg message tone: " + floatStr(ws.AvgTonePerMsg, 4))
	fmt.Println("Avg message analytic: " + floatStr(ws.AvgAnalyticPerMsg, 4))
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

// inListBinary determines whether a word is contained in a slice
// of words using binary search (assumes words is sorted)
func inListBinary(word string, words []string) bool {
	i := 0
	j := len(words) - 1
	for i <= j {
		k := (i + j) / 2
		w := words[k]
		if word == w {
			return true
		}
		if word < w {
			j = k - 1
		} else if word > w {
			i = k + 1
		}
	}
	return false
}

// floatStr converts a float64 to a string
// with decimals worth of precision
func floatStr(f float64, decimals int) string {
	return strconv.FormatFloat(f, 'f', decimals, 64)
}

func newWordStats() *WordStats {
	return &WordStats{
		TotalTextLength:   0,
		TotalWords:        0,
		TotalMessages:     0,
		AvgWordLength:     0,
		AvgWordsPerMsg:    0,
		AvgCloutPerMsg:    0,
		AvgTonePerMsg:     0,
		AvgAnalyticPerMsg: 0,
		WordCountMap:      make(map[string]int),
		CategoryCounts:    make(map[string]int),
	}
}

func newMessageStats() *MessageStats {
	return &MessageStats{
		NumMessages:       0,
		NumWords:          0,
		NumEmojis:         0,
		AvgWordLength:     0,
		AvgWordsPerMsg:    0,
		AvgEmojisPerMsg:   0,
		AvgCloutPerMsg:    0,
		AvgTonePerMsg:     0,
		AvgAnalyticPerMsg: 0,
		WordCountMap:      make(map[string]int),
		EmojiCountMap:     make(map[string]int),
		CategoryCounts:    make(map[string]int),
	}
}

func updateWordCountMap(w string, wcm *map[string]int) {
	count, ok := (*wcm)[w]
	if !ok {
		(*wcm)[w] = 1
	} else {
		(*wcm)[w] = count + 1
	}
}

func populateCategoryCounts(ws *WordStats, wordCategoriesCache *map[string][]string) {
	for word, count := range (*ws).WordCountMap {
		categories, hit := (*wordCategoriesCache)[word]
		if !hit {
			categories = GetCategories(word)
			(*wordCategoriesCache)[word] = categories
		}
		for _, cat := range categories {
			c, ok := (*ws).CategoryCounts[cat]
			if ok {
				(*ws).CategoryCounts[cat] = c + count
			} else {
				(*ws).CategoryCounts[cat] = count
			}
		}
	}
}

func populateMsgCategoryCounts(ms *MessageStats, wordCategoriesCache *map[string][]string) {
	for word, count := range (*ms).WordCountMap {
		categories, hit := (*wordCategoriesCache)[word]
		if !hit {
			categories = GetCategories(word)
			(*wordCategoriesCache)[word] = categories
		}
		for _, cat := range categories {
			c, ok := (*ms).CategoryCounts[cat]
			if ok {
				(*ms).CategoryCounts[cat] = c + count
			} else {
				(*ms).CategoryCounts[cat] = count
			}
		}
	}
}

func setAverages(ws *WordStats) {
	totalWords := float64(ws.TotalWords)
	totalMessages := float64(ws.TotalMessages)
	ws.AvgWordLength /= totalWords
	ws.AvgWordsPerMsg = totalWords / totalMessages
	ws.AvgCloutPerMsg /= totalMessages
	ws.AvgTonePerMsg /= totalMessages
	ws.AvgAnalyticPerMsg /= totalMessages
}

func setStatAverages(ms *MessageStats) {
	numWords := float64(ms.NumWords)
	numEmojis := float64(ms.NumEmojis)
	numMsg := float64(ms.NumMessages)
	ms.AvgWordLength /= numWords
	ms.AvgWordsPerMsg = numWords / numMsg
	ms.AvgEmojisPerMsg = numEmojis / numMsg
	ms.AvgCloutPerMsg /= numMsg
	ms.AvgTonePerMsg /= numMsg
	ms.AvgAnalyticPerMsg /= numMsg
}
