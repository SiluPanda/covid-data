package utils

import (
	"encoding/json"
	"errors"
	"inshorts/models"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func FetchReverseGeo(lat string, long string) (models.ReverseGeo, error) {
	godotenv.Load()
	client := http.Client{}

	req, _ := http.NewRequest("GET", "http://api.positionstack.com/v1/reverse?access_key="+os.Getenv("POSITION_STACK_API_KEY")+"&query="+lat+","+long, nil)

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return models.ReverseGeo{}, err
	}

	if resp.StatusCode != 200 {
		return models.ReverseGeo{}, errors.New("can not fetch location")
	}

	m := make(map[string][]models.ReverseGeo)
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &m)

	return m["data"][0], nil

}
