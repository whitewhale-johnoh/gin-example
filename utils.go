package main

import (
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
	"gopkg.in/redis.v5"
)

func generateuuid() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	} else {
		return u.String(), nil
	}
}

func formatTime(hour int, min int) string {
	formattedhour := strconv.Itoa(hour)
	formattedmin := strconv.Itoa(min)
	if formattedhour == "0" {
		formattedhour = "00"
	}
	if formattedmin == "0" {
		formattedmin = "00"
	}
	return formattedhour + ":" + formattedmin

}

func ping(client *redis.Client) error {
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return err

	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}
