package cmd

import (
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func replaceStringTerms(packet string, terms []string) (string, bool) {
	if len(terms) == 0 {
		return packet, true
	}

	foundMatch := false
	for _, term := range terms {
		if strings.Contains(packet, term) {
			// Colour the output to highlight the found value
			packet = strings.ReplaceAll(packet, term, color.CyanString(term))
			foundMatch = true
		}
	}

	return packet, foundMatch
}

func replaceTerms(packet string, terms []string) (string, bool) {
	if len(terms) == 0 {
		return packet, true
	}

	if os.Getenv("WATCHTOWER_USE_REGEX") != "1" {
		return replaceStringTerms(packet, terms)
	}

	foundMatch := false
	for _, term := range terms {
		regex := regexp.MustCompile(term)
		if regex.MatchString(packet) {
			packet = regex.ReplaceAllStringFunc(packet, func(s string) string {
				return color.CyanString(s)
			})
			foundMatch = true
		}
	}

	return packet, foundMatch

}
