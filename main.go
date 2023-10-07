package main

import (
  "fmt"
  "os"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strings"
)

type DictionaryEntry []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text      string `json:"text"`
		Audio     string `json:"audio"`
		SourceURL string `json:"sourceUrl,omitempty"`
		License   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license,omitempty"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Synonyms   []string  `json:"synonyms"`
			Antonyms   []string  `json:"antonyms"`
      Example    string  `json:"example"`
		} `json:"definitions"`
		Synonyms []string `json:"synonyms"`
		Antonyms []string `json:"antonyms"`
	} `json:"meanings"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	SourceUrls []string `json:"sourceUrls"`
}

type Color int

const (
  RED = 31
  GREEN = 32
  YELLOW = 33
  BLUE = 34
  MAGENTA = 35
)

func colorize (s string, c Color) string {
  return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, s)
}

func italicize (s string) string {
  return fmt.Sprintf("\x1b[3m%s\x1b[0m", s)
}

func main() {
  word := strings.Join(os.Args[1:], " ")
  resp, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + word)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer resp.Body.Close()
  if resp.StatusCode == 404 {
    fmt.Println("Word not found.")
    os.Exit(1)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  var result DictionaryEntry
  if err := json.Unmarshal(body, &result); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  foundWord := colorize(strings.ToUpper(result[0].Word), RED)

  var meaningStrings string
  for index, meaning := range(result[0].Meanings) {
    partOfSpeech := colorize(meaning.PartOfSpeech, MAGENTA)
    meaningStr := fmt.Sprintf(`Meaning #%v:

  Part of Speech: %s
  `, index + 1, partOfSpeech)
    for index, definition := range(meaning.Definitions) {
      definitionStyled := colorize(italicize(definition.Definition), GREEN)
      meaningStr += fmt.Sprintf(`
  Definition #%v:
    %s
  `, index + 1, definitionStyled)
      if (definition.Example != "") {
        meaningStr += fmt.Sprintf("  Example: %s\n", italicize(colorize(definition.Example, YELLOW)))
      }
    }
    
    meaningStrings += meaningStr + "\n"
  }
  fmt.Printf(`
Definition for %s:

%s`, foundWord, meaningStrings) 
} 
