package es

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/esapi"
	"log"
	"rider-device/initenal/basic/config"
	"rider-device/initenal/model"
	"strconv"
	"strings"
)

// net/http post demo

func SetEs(riders *model.Riders) {
	marshal, err := json.Marshal(riders)
	if err != nil {
		return
	}
	req := esapi.IndexRequest{
		Index:      "riders",
		DocumentID: strconv.Itoa(int(riders.Id)),
		Body:       strings.NewReader(string(marshal)),
		Refresh:    "true",
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), config.Es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

}
