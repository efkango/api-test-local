package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	apiKey   = "8e172b1d-7696-4dab-9ab6-fc16b727893d:fx"
	apiURL   = "https://api-free.deepl.com/v2/translate"
	numTests = 500 // cevirilecek kelimeler
)

type TranslationRequest struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

type TranslationResponse struct {
	Translations []struct {
		Text string `json:"text"`
	} `json:"translations"`
}

var englishWords = []string{
	"apple", "banana", "cat", "dog", "elephant", "frog", "giraffe", "house", "ice cream", "jacket",
	"kite", "lemon", "monkey", "notebook", "orange", "pencil", "queen", "rabbit", "sun", "tree",
	"umbrella", "violin", "water", "xylophone", "yacht", "zebra", "book", "car", "door", "egg",
	"flower", "guitar", "hat", "island", "juice", "key", "lamp", "moon", "nest", "ocean",
	"piano", "quilt", "river", "star", "table", "unicorn", "vase", "window", "x-ray", "yarn",
	"airplane", "ball", "cloud", "dance", "eagle", "fan", "garden", "hill", "ice", "jungle",
	"kangaroo", "leaf", "mountain", "nose", "octopus", "pen", "quicksand", "rocket", "shark", "tiger",
	"umbrella", "vulture", "whale", "xenon", "yogurt", "zephyr", "alarm", "brush", "clock", "desk",
	"engine", "flag", "goal", "helicopter", "idea", "jewel", "kettle", "ladder", "mirror", "needle",
	"oven", "plate", "question", "rope", "ship", "telephone", "ufo", "van", "wheel", "xerox", "yawn",
	"zoo", "arrow", "button", "candle", "dolphin", "eraser", "fork", "glasses", "hammer", "ink",
	"jet", "knife", "lion", "magnet", "nail", "owl", "plant", "queen", "ring", "spoon", "train",
	"umbrella", "viper", "wallet", "xylophone", "yellow", "zipper", "ant", "bread", "cheese", "duck",
	"envelope", "fire", "globe", "honey", "island", "jacket", "kite", "leopard", "marble", "needle",
	"onion", "pan", "quill", "robot", "snow", "tractor", "unicorn", "volcano", "wind", "x-ray",
	"yarn", "zucchini", "air", "balloon", "cup", "drum", "elephant", "fence", "gold", "helmet", "iceberg",
	"jar", "keychain", "lamp", "microscope", "net", "ocean", "paint", "quiver", "radio", "sunflower",
	"tractor", "umbrella", "village", "watermelon", "xylophone", "yellow", "zipper", "avocado", "beach",
	"cloud", "daisy", "elephant", "fan", "giraffe", "hat", "igloo", "jug", "kite", "lemon", "map",
	"nut", "owl", "pizza", "queen", "rose", "sun", "tree", "umbrella", "vase", "whistle", "xylophone",
	"yo-yo", "zebra", "apple", "basket", "crayon", "dog", "egg", "fire", "guitar", "house", "ink", "jelly",
	"kitten", "lion", "mushroom", "nest", "octopus", "pencil", "quilt", "ring", "sand", "tent", "umbrella",
	"vase", "whale", "x-ray", "yacht", "zoo", "anchor", "ball", "chair", "door", "envelope", "feather",
	"guitar", "hat", "igloo", "jar", "kite", "lamp", "mirror", "needle", "orange", "pear", "queen",
	"robot", "stone", "tree", "unicorn", "vase", "wheel", "xylophone", "yo-yo", "zebra", "ant",
	"butterfly", "carrot", "duck", "eagle", "fish", "grape", "hat", "insect", "jacket", "kangaroo",
	"lemon", "moon", "nest", "ostrich", "penguin", "queen", "rabbit", "star", "turtle", "umbrella",
	"violin", "whale", "x-ray", "yarn", "zebra", "airplane", "banana", "cloud", "drum", "egg", "fan",
	"guitar", "hat", "ice", "jacket", "kite", "lemon", "mountain", "nest", "octopus", "pizza", "quilt",
	"ring", "star", "tree", "umbrella", "violin", "whale", "xylophone", "yarn", "zebra", "ant",
	"bird", "cat", "dog", "elephant", "frog", "giraffe", "hat", "insect", "jacket", "kangaroo",
	"lemon", "moon", "nest", "octopus", "penguin", "queen", "rabbit", "sun", "tiger", "umbrella",
	"violin", "whale", "x-ray", "yacht", "zebra", "apple", "ball", "cat", "dog", "elephant",
	"frog", "giraffe", "house", "ice", "jacket", "kite", "lemon", "monkey", "nest", "octopus",
	"penguin", "queen", "rabbit", "star", "tiger", "umbrella", "violin", "whale", "xylophone",
	"yacht", "zebra", "apple", "banana", "cat", "dog", "elephant", "frog", "giraffe", "house",
	"ice", "jacket", "kite", "lemon", "monkey", "notebook", "orange", "pencil", "queen", "rabbit",
	"sun", "tree", "umbrella", "violin", "water", "xylophone", "yacht", "zebra"}

func getRandomWord() string {
	return englishWords[rand.Intn(len(englishWords))]
}

func translateText(client *http.Client, text string, targetLang string) (string, error) {
	reqBody, err := json.Marshal(TranslationRequest{
		Text:       []string{text},
		TargetLang: targetLang,
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", body)
	}

	var translationResp TranslationResponse
	err = json.Unmarshal(body, &translationResp)
	if err != nil {
		return "", err
	}

	if len(translationResp.Translations) > 0 {
		return translationResp.Translations[0].Text, nil
	}

	return "", fmt.Errorf("no translation found")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	client := &http.Client{Timeout: 10 * time.Second}
	var wg sync.WaitGroup
	results := make(chan string, numTests)

	startTime := time.Now()

	for i := 0; i < numTests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			word := getRandomWord()
			translation, err := translateText(client, word, "TR")
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
	for i, translation := range translations[:250] {
		fmt.Printf("%d. %s\n", i+1, translation)
	}
}
