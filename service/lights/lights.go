package lights

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func ToggleLights() {
	godotenv.Load()
	token := os.Getenv("TOKEN")

	haURL := "http://localhost:8123/api/services/light/toggle"

	body := map[string]string{
		"entity_id": "light.office_lights",
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", haURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
