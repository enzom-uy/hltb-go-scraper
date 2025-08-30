package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/enzom-uy/hltb-go-scraper/internal/models"
)

type QueryGameResponse struct {
	GameTitle     string
	GameDurations models.GameDurations
}

func QueryGame(ctx context.Context, gameName string) (*QueryGameResponse, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	gameName = strings.TrimSpace(gameName)
	gameName = strings.Trim(gameName, `"'`)

	if gameName == "" {
		fmt.Println("Game name is empty.")
		return nil, errors.New("Game name is empty.")
	}
	if len(gameName) > 50 {
		fmt.Println("Game name is too long (max 50 characters).")
		return nil, errors.New("Game name is too long (max 50 characters).")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromedpCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	chromedpCtx, cancel = context.WithTimeout(chromedpCtx, 60*time.Second)
	defer cancel()

	var finalURL string

	fmt.Printf("Trying to scrap %v data from HowLongToBeat.", gameName)

	select {
	case <-ctx.Done():
		log.Println("Request cancelled before navigation")
		return nil, ctx.Err()
	default:
	}

	err := chromedp.Run(chromedpCtx,
		chromedp.Navigate("https://howlongtobeat.com/"),
		chromedp.WaitVisible(`input[type="search"]`),
		chromedp.SendKeys(`input[type="search"]`, gameName),
		chromedp.KeyEvent("\r"),
		chromedp.Sleep(3*time.Second),
		chromedp.WaitVisible(`#search-results-header`, chromedp.ByQuery),
		chromedp.Location(&finalURL),
	)

	fmt.Println("Error?: ", err)

	if err != nil {
		if ctx.Err() == context.Canceled {
			log.Println("Navigation cancelled by client")
			return nil, ctx.Err()
		}
		return &QueryGameResponse{}, errors.New("(Chromedp) There was an error when trying to navigate and scrape the data.")
	}

	fmt.Println("Final URL: ", finalURL)

	select {
	case <-ctx.Done():
		log.Println("Request cancelled before getting HTML")
		return nil, ctx.Err()
	default:
	}

	var htmlContent string
	err = chromedp.Run(chromedpCtx,
		chromedp.InnerHTML("body", &htmlContent),
	)

	if err != nil {
		if ctx.Err() == context.Canceled {
			log.Println("HTML retrieval cancelled by client")
			return nil, ctx.Err()
		}
		return &QueryGameResponse{}, errors.New("(Chromedp) There was an error when trying to get scrapped website HTML.")
	}

	fmt.Printf("HTML retrieved: %d characters.\n", len(htmlContent))

	select {
	case <-ctx.Done():
		log.Println("Request cancelled before parsing HTML")
		return nil, ctx.Err()
	default:
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return &QueryGameResponse{}, errors.New("(goquery) Error when trying to parse the HTML document.")
	}

	firstGame := doc.Find("#search-results-header ul li").First()
	if firstGame.Length() == 0 {
		firstGame = doc.Find("li.GameCard_search_list__IuMbi").First()
		if firstGame.Length() == 0 {
			return &QueryGameResponse{}, errors.New("No game found.")
		}
	}

	gameTitle := strings.TrimSpace(firstGame.Find("h2 a").Text())
	mainStoryLength := strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").First().Text())
	mainExtraLength := strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").Eq(1).Text())
	completionistLength := strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").Eq(2).Text())

	fmt.Println("✅ Website scrapped successfully.")
	fmt.Println("Game title: ", gameTitle)
	fmt.Println("Main story duration: ", mainStoryLength)
	fmt.Println("Main story + extras duration: ", mainExtraLength)
	fmt.Println("Completionist duration: ", completionistLength)

	return &QueryGameResponse{
		GameTitle: gameTitle,
		GameDurations: models.GameDurations{
			MainStory:     mainStoryLength,
			MainsSides:    mainExtraLength,
			Completionist: completionistLength,
		},
	}, nil
}
