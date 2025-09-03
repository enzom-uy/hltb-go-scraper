package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/enzom-uy/hltb-go-scraper/internal/db"
	"github.com/enzom-uy/hltb-go-scraper/internal/model"
	"github.com/enzom-uy/hltb-go-scraper/internal/models"
)

type QueryGameResponse struct {
	GameTitle     string
	GameDurations models.GameDurations
	GameID        string
	GameURL       string
}

type ParsedGameResponse struct {
	duration float64
}

func parseHoursStringToInt(string string) (*ParsedGameResponse, error) {
	parsedString := strings.ReplaceAll(string, "½", ".5")
	parsedString = strings.ReplaceAll(parsedString, " Hours", "")
	toNumber, err := strconv.ParseFloat(parsedString, 64)

	if err != nil {
		fmt.Println("Error: ", err)
		return &ParsedGameResponse{}, err
	}

	fmt.Println("nuevo main text: ", toNumber)
	return &ParsedGameResponse{
		duration: toNumber,
	}, nil

}

func QueryGame(ctx context.Context, gameName string) (*QueryGameResponse, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	db, _, dbErr := db.Init()
	if dbErr != nil {
		fmt.Println("Error al conectar a la base de datos: ", dbErr)
		return nil, dbErr
	}

	exists := db.Where("title LIKE ?", "%"+gameName+"%").First(&model.Game{}).RowsAffected > 0
	fmt.Println("Existe el juego ya?: ", exists)

	if exists {
		fmt.Println("%v exists in database.", gameName)
		// TODO: don't return error, just return the game.
		return nil, errors.New("Game already exists in database.")
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
		chromedp.Sleep(2*time.Second),
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
	fmt.Println("First game: ", firstGame.Text())
	if firstGame.Length() == 0 {
		firstGame = doc.Find("li.GameCard_search_list__IuMbi").First()
		if firstGame.Length() == 0 {
			return &QueryGameResponse{}, errors.New("No game found.")
		}
	}

	gameTitle := strings.TrimSpace(firstGame.Find("h2 a").Text())
	mainStoryLength, parseError := parseHoursStringToInt(strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").First().Text()))
	mainExtraLength, parseError := parseHoursStringToInt(strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").Eq(1).Text()))
	completionistLength, parseError := parseHoursStringToInt(strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").Eq(2).Text()))

	if parseError != nil {
		return &QueryGameResponse{}, errors.New("(parse) Error when trying to parse the durations.")
	}

	gameUrl := strings.TrimSpace(firstGame.Find("h2 a").AttrOr("href", ""))
	splitUrl := strings.Split(gameUrl, "/")
	gameID := splitUrl[len(splitUrl)-1]
	gameIDInt, strConvErr := strconv.ParseInt(gameID, 10, 64)

	if strConvErr != nil {
		fmt.Println("Error converting game ID to int64: ", strConvErr)
		return &QueryGameResponse{}, strConvErr
	}

	fmt.Println("✅ Website scrapped successfully.")
	fmt.Println("Game title: ", gameTitle)
	fmt.Println("Main story duration: ", mainStoryLength.duration)
	fmt.Println("Main story + extras duration: ", mainExtraLength.duration)
	fmt.Println("Completionist duration: ", completionistLength.duration)
	fmt.Println("Game URL: ", "https://howlongtobeat.com"+gameUrl)

	fmt.Println("Trying to save data to database.")
	// TODO: handle errors
	newGame := model.HowlongtobeatDatum{
		HltbID:              gameIDInt,
		MainStoryHours:      mainStoryLength.duration,
		MainStorySidesHours: mainExtraLength.duration,
		CompletionistHours:  completionistLength.duration,
	}

	db.Create(&newGame)

	return &QueryGameResponse{
		GameTitle: gameTitle,
		GameDurations: models.GameDurations{
			MainStory:     mainStoryLength.duration,
			MainsSides:    mainExtraLength.duration,
			Completionist: completionistLength.duration,
		},
		GameID:  gameID,
		GameURL: "https://howlongtobeat.com" + gameUrl,
	}, nil
}
