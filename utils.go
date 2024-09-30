package main

import (
	"log"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
)

// arrayFlags is a custom flag type that allows us to pass multiple values
type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	newStr := strings.TrimSpace(value)
	if slices.Contains(*i, newStr) {
		return nil
	}

	*i = append(*i, newStr)

	return nil
}

func getFromEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getFromEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalln("Invalid value for " + key + ": " + value)
	}

	return i
}

func getFromEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalln("Invalid value for " + key + ": " + value)
	}

	return b
}

func parseUrl(rawUrl string) *url.URL {
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Fatalln("Invalid URL: " + rawUrl)
	}

	return u
}

func getSliceFromString(value string) []string {
	slice := make([]string, 0)
	for _, v := range strings.Split(value, ",") {
		str := strings.TrimSpace(v)
		if str == "" {
			continue
		}

		slice = append(slice, str)
	}

	return slice
}

func getFromEnvStringSlice(key string, defaults []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaults
	}

	return getSliceFromString(value)
}
