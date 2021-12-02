package main

import (
	"bufio"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var words []string = make([]string, 0)

var specialCharacters = []string{
	" ",
	"!",
	"\"",
	"#",
	"$",
	"%",
	"&",
	"'",
	"(",
	")",
	"*",
	"+",
	",",
	"-",
	".",
	"/",
	":",
	";",
	"<",
	"=",
	">",
	"?",
	"@",
	"[",
	"\\",
	"]",
	"^",
	"_",
	"`",
	"{",
	"|",
	"}",
	"~",
}

var defaultWordCount = 3

func readWords() {
	log.Println("Reading wordlist")
	file, err := os.Open("words.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, strings.Title(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func getAvailableWord(usedWords map[string]bool, iteration int) string {
	wordCandidate := words[rand.Intn(len(words))]
	if usedWords[wordCandidate] && iteration < 10 {
		return getAvailableWord(usedWords, iteration+1)
	}
	return wordCandidate
}

func buildPassword(wordCount int, specialChars bool) string {
	var password string
	usedWords := map[string]bool{}

	for i := 0; i < wordCount; i++ {
		word := getAvailableWord(usedWords, 0)
		password = password + word
		usedWords[word] = true
		if specialChars {
			password = password + specialCharacters[rand.Intn(len(specialCharacters))]
		}
	}

	return password
}

func generatePassword(wordCount int, specialChars bool) string {
	switch {
	case wordCount <= 0:
		return generatePassword(3, specialChars)
	case wordCount > 10:
		return generatePassword(3, specialChars)
	default:
		return buildPassword(wordCount, specialChars)
	}
}

func main() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	readWords()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/", func(c *gin.Context) {
		wordCount, wordCountErr := strconv.Atoi(c.Query("words"))
		specialChars, specialCharErr := strconv.ParseBool(c.Query("special"))
		if wordCountErr != nil {
			wordCount = defaultWordCount
		}
		if specialCharErr != nil {
			specialChars = true
		}
		c.JSON(http.StatusOK, gin.H{"password": generatePassword(wordCount, specialChars)})
	})

	r.Run()
}
