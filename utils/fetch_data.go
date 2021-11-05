package utils

import (
	"encoding/json"
	"fmt"
	"inshorts/models"
	"io/ioutil"
	"net/http"
)

func FetchData() []models.CovidData {

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://data.covid19india.org/v4/min/data.min.json", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
	}

	m := make(map[string]models.CovidData)

	data, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &m)

	var covid_data_list []models.CovidData
	for k, v := range m {
		curr_data := v
		curr_data.State = k
		covid_data_list = append(covid_data_list, curr_data)
	}

	return covid_data_list
}
