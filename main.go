package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	apiKey   = "8e172b1d-7696-4dab-9ab6-fc16b727893d:fx"
	apiURL   = "https://api-free.deepl.com/v2/translate"
	numTests = 1 // çevrilecek kelimeler
)

type TranslationRequest struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
	IgnoreTags []string `json:"ignore_tags"`
}

type TranslationResponse struct {
	Translations []struct {
		Text string `json:"text"`
	} `json:"translations"`
}

var englishWords = []string{
	"Merhaba su icer misin?<@4c1ef347-25e3-42d6-9bcf-45a250dca6c9:merhaba>",
}

func getRandomWord() string {
	return englishWords[rand.Intn(len(englishWords))]
}

func translateText(ctx context.Context, text string, targetLang string) (string, error) {
	text = wrapUserMentionsWithXTags(text)

	reqBody, err := json.Marshal(map[string]interface{}{
		"text":        []string{text},
		"target_lang": targetLang,
		"ignore_tags": []string{},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			err := "DeepL API request unauthorized, check API key"
			return "", fmt.Errorf(err)
		} else {
			errMsg := fmt.Sprintf("DeepL API request failed with status code: %d", resp.StatusCode)
			return "", fmt.Errorf(errMsg)
		}
	}

	var translationResp TranslationResponse
	if err := json.NewDecoder(resp.Body).Decode(&translationResp); err != nil {
		return "", fmt.Errorf("failed to decode translation response: %v", err)
	}

	if len(translationResp.Translations) > 0 {
		// <x> kaldir
		return removeXTags(translationResp.Translations[0].Text), nil
	}

	return "", fmt.Errorf("no translation found")
}

func wrapUserMentionsWithXTags(text string) string {
	text = regexp.MustCompile(`(@\S+)`).ReplaceAllString(text, "<x>$1</x>")
	return regexp.MustCompile(`(<@[a-f0-9-]+:[^>]+>)`).ReplaceAllString(text, "<x>$1</x>")
}

func removeXTags(text string) string {
	text = strings.ReplaceAll(text, "<x>", "")
	return strings.ReplaceAll(text, "</x>", "")
}

func main() {
	var wg sync.WaitGroup
	results := make(chan string, numTests)
	ctx := context.Background()
	startTime := time.Now()

	for i := 0; i < numTests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			word := getRandomWord()
			translation, err := translateText(ctx, word, "EN")
			if err != nil {
				results <- fmt.Sprintf("Error translating '%s': %v", word, err)
			} else {
				results <- fmt.Sprintf("%s: %s", word, translation)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var translations []string
	for result := range results {
		translations = append(translations, result)
	}

	duration := time.Since(startTime)

	fmt.Printf("Toplam %d kelime %.2f saniyede çevrildi.\n", numTests, duration.Seconds())
	fmt.Printf("Saniyede ortalama %.2f çeviri yapıldı.\n", float64(numTests)/duration.Seconds())
	fmt.Println("İlk 250 çeviri sonucu:")
	for i, translation := range translations[:1] {
		fmt.Printf("%d. %s\n", i+1, translation)
	}
}
