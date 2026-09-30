package helpers

import (
	"regexp"
	"strings"
)

// GenerateSlug generates a clean URL slug from string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(strings.TrimSpace(text))

	// Replace spaces with hyphens
	reSpaces := regexp.MustCompile(`\s+`)
	slug = reSpaces.ReplaceAllString(slug, "-")

	// Remove all non-alphanumeric and non-hyphen characters
	reInvalid := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reInvalid.ReplaceAllString(slug, "")

	// Remove duplicate hyphens
	reDuplicates := regexp.MustCompile(`\-+`)
	slug = reDuplicates.ReplaceAllString(slug, "-")

	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
