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
	"The cat jumped over the fence to chase butterflies in garden.",
	"My favorite fruit is a banana, but I love apples too.",
	"The dog barked loudly at the moon shining brightly in sky.",
	"I bought a new pencil and notebook for my art class.",
	"The giraffe reached up to eat leaves from the tallest tree.",
	"We played with a kite near the beach and flew it high.",
	"The ice cream melted quickly in the hot summer sun today.",
	"I saw a frog jump into the pond near my grandmother's house.",
	"She carried an umbrella to protect herself from the heavy rain.",
	"He played a beautiful song on the violin for the audience.",

	"The apple fell from the tree and rolled down the hill.",
	"My new jacket kept me warm in the cold winter weather.",
	"The elephant drank water from the river using its long trunk.",
	"I wrote a letter with my pencil on a clean notebook.",
	"The sun shone brightly over the garden full of colorful flowers.",
	"We saw a monkey swinging between trees at the big zoo.",
	"She painted a beautiful picture of a house by the lake.",
	"He used a key to unlock the door of the old house.",
	"They built a sandcastle near the ocean and enjoyed the waves.",
	"I saw a rainbow after the rain near my big backyard.",

	"The frog jumped into the water and swam back to shore.",
	"He played his guitar while sitting under the tall oak tree.",
	"She wore a hat to protect her head from the hot sun.",
	"The bird flew from the tree to the top of the roof.",
	"He wrote his homework in a notebook with his favorite pencil.",
	"They flew a kite in the park on a very windy day.",
	"The dog chased the cat around the house all afternoon long.",
	"The elephant used its trunk to pick up the fallen apple.",
	"She sat under the tree and read a book very quietly.",
	"The sun set behind the mountains, turning the sky bright orange.",

	"He played the piano for hours practicing his favorite classical songs.",
	"The monkey grabbed a banana and climbed up the big tree.",
	"She brought a water bottle and an umbrella to the park.",
	"The rabbit hopped through the grass, looking for some fresh food.",
	"The star shone brightly in the clear night sky last night.",
	"He opened the door with a key and walked inside slowly.",
	"The boat sailed down the river on a bright sunny day.",
	"I saw a giraffe eating leaves from the tallest tree nearby.",
	"The flower bloomed in the garden next to the little house.",
	"He sat on a bench, reading a book under the big tree.",

	"We saw a tiger pacing back and forth in the big cage.",
	"The violinist played a beautiful melody at the concert last night.",
	"The airplane flew over the mountains and disappeared into the clouds.",
	"I found my jacket hanging on a hook by the door.",
	"The cat curled up in the sun and fell asleep peacefully.",
	"The river flowed through the valley, surrounded by tall green trees.",
	"He packed a sandwich, juice, and book for the picnic today.",
	"The frog leaped from one rock to another near the flowing stream.",
	"She wore her new jacket while walking through the city park.",
	"The sun was setting as we walked along the sandy beach.",

	"The dog barked loudly when he saw the delivery man arrive.",
	"She bought a new notebook and pencil for her English class.",
	"The tree was full of apples, ready to be picked soon.",
	"He played the guitar while his friends sang along happily together.",
	"The house on the hill had a beautiful view of the city.",
	"The ice cream truck drove by, playing its familiar catchy tune.",
	"The kite soared high in the sky on the windy afternoon.",
	"He painted a picture of the ocean with his new watercolors.",
	"The giraffe bent down to drink water from the small pond.",
	"She took a picture of the rainbow after the heavy storm.",

	"He wrote a letter to his friend in a lined notebook.",
	"The sun set over the horizon, casting a warm orange glow.",
	"She sat by the river, watching the boats sail peacefully by.",
	"The bird built a nest in the tree near the house.",
	"The rabbit hopped through the garden, nibbling on the fresh carrots.",
	"The monkey swung from branch to branch in the dense jungle.",
	"He unlocked the door with the key and stepped inside quietly.",
	"The violinist performed a beautiful piece at the evening concert hall.",
	"The airplane took off from the runway and flew into the sky.",
	"She wore a hat to protect her face from the bright sun.",

	"The frog jumped into the pond and swam towards the lily pad.",
	"The tiger roared loudly in its cage at the busy zoo.",
	"The sun rose early, bathing the mountains in a golden light.",
	"He packed his bag with a notebook and water bottle for school.",
	"The giraffe reached up to eat the leaves from the tall branches.",
	"The house on the corner had a red door and big windows.",
	"She sat under the umbrella and read a book very quietly.",
	"He played the violin while everyone listened in complete silence.",
	"The river flowed peacefully through the valley surrounded by tall trees.",
	"She bought a new jacket for the upcoming cold winter season.",

	"The bird sang a beautiful song from its perch in the tree.",
	"The rabbit dug a hole in the garden near the white fence.",
	"He played soccer with his friends in the park after school.",
	"The cat chased a butterfly through the colorful flowers in the garden.",
	"The sun shone brightly over the ocean on a warm summer day.",
	"The monkey stole a banana and ran away into the jungle quickly.",
	"She found a book on the table and began reading it slowly.",
	"He drew a picture of the mountains with his favorite pencil.",
	"The ice cream melted quickly in the hot afternoon sun today.",
	"We flew a kite by the beach and watched it soar high.",

	"The frog croaked loudly near the pond as the sun was setting.",
	"She carried an umbrella because it was raining heavily in the city.",
	"The elephant splashed water with its trunk to cool off in heat.",
	"He used a key to open the door of the old house.",
	"The cat curled up in the sun and took a peaceful nap.",
	"The rabbit hopped through the grass, looking for something to eat quickly.",
	"The airplane flew over the ocean and landed safely at the airport.",
	"He played the piano for hours, practicing his favorite classical music.",
	"The dog chased the ball down the street and caught it happily.",
	"The giraffe bent down to drink water from the river in park.",

	"The monkey grabbed a banana from the tree and ate it quickly.",
	"The bird built its nest in the tree near the park entrance.",
	"She wore a hat and sunglasses to protect herself from the sun.",
	"The frog jumped into the pond and swam towards the lily pad.",
	"The tiger roared loudly in the zoo, scaring the nearby visitors.",
	"The sun rose early, painting the sky with beautiful shades of orange.",
	"He packed a sandwich and a juice box for his picnic today.",
	"The giraffe reached up and grabbed some leaves from the tall tree.",
	"The house on the hill had a red door and white windows.",
	"She sat under the umbrella and read a book while sipping tea.",

	"He played the violin in the park while people gathered to listen.",
	"The river flowed gently through the valley, surrounded by tall green trees.",
	"She bought a new jacket and scarf for the upcoming winter season.",
	"The bird sang a cheerful song from the tree branch near the house.",
	"The rabbit hopped through the field, nibbling on grass and flowers nearby.",
	"The airplane flew high in the sky, leaving a trail of white clouds.",
	"The monkey swung from tree to tree, looking for bananas to eat.",
	"He used a pencil to sketch a picture of the mountains in notebook.",
	"The sun shone brightly over the ocean, creating a beautiful golden reflection.",
	"The cat stretched out in the sun and fell asleep on porch.",

	"The dog chased a butterfly across the yard, barking excitedly as it ran.",
	"She carried an umbrella to stay dry during the sudden afternoon rainstorm.",
	"The ice cream truck played a cheerful tune as it drove by quickly.",
	"He played his guitar while sitting under the tree near the river.",
	"The giraffe bent down to drink water from the small pond in park.",
	"The frog croaked loudly as it jumped into the cool water of pond.",
	"The sun set behind the mountains, casting a warm orange glow over the sky.",
	"The rabbit dug a hole in the garden, looking for fresh vegetables to eat.",
	"The elephant splashed water on its back, cooling off from the hot sun.",
	"She wore a sun hat to protect her face from the bright summer sun."}

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
	for i, translation := range translations[:100] {
		fmt.Printf("%d. %s\n", i+1, translation)
	}
}
