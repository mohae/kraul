package app

import "github.com/mohae/kraul/indexer"

// Index indexes the site at the starturl. URLs not in the startUrl domain
// are not crawled
func Index(startUrl string) (string, error) {
	return indexer.Index(startUrl)
}
