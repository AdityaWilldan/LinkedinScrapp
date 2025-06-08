package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mbndr/figlet4go"
)

func main() {

	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorRed,
	}
	renderStr, _ := ascii.RenderOpts("LinkedIn Scrapper", options)
	fmt.Print(renderStr)
	fmt.Print("tools by: Wolfgang\n")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input URL: ")
	inputURL, _ := reader.ReadString('\n')
	inputURL = strings.TrimSpace(inputURL)

	fmt.Print("Selector / Tag: ")
	selector, _ := reader.ReadString('\n')
	selector = strings.TrimSpace(selector)

	resp, err := http.Get(inputURL)
	if err != nil {
		log.Fatalf("error fetch the page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("error parsing html: %v", err)
	}

	var links []string

	count := 1
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		link, exsists := s.Attr("href")
		if exsists {
			if isLinkedInJobLink(link) {
				fmt.Printf("url found: %d | %s\n", count, link)
				links = append(links, link)
				count++
			}
		}
	})
	if len(links) == 0 {
		fmt.Println("no URLs found")
		return
	}

	fmt.Print("export to CSV?(y/n): ")
	exportChoice, _ := reader.ReadString('\n')
	exportChoice = strings.TrimSpace(strings.ToLower(exportChoice))

	if exportChoice == "y" {
		err := writeCSV(links)
		if err != nil {
			log.Fatalf("failed convert CSV %v", err)
		}
	}
}

func isLinkedInJobLink(link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	if strings.Contains(parsedURL.Host, "linkedin.com") && strings.HasPrefix(parsedURL.Path, "/jobs/view/") {
		return true
	}
	if strings.HasPrefix(link, "/jobs/views/") {
		return true
	}
	return false
}

func writeCSV(data []string) error {
	file, err := os.Create("outpus.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"no", "linkedin URL"})
	for i, url := range data {
		writer.Write([]string{fmt.Sprintf("%d", i+1), url})
	}
	return nil
}
