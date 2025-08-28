package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func searchGame(gameName string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.Flag("disable-gpu", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var finalURL string

	fmt.Printf("Intentando obtener los datos de %v en HowLongToBeat.", gameName)
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://howlongtobeat.com/"),
		chromedp.WaitVisible(`input[type="search"]`),
		// chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(`input[type="search"]`, gameName),
		chromedp.KeyEvent("\r"),
		chromedp.Sleep(3*time.Second),
		chromedp.WaitVisible(`#search-results-header`, chromedp.ByQuery),
		chromedp.Location(&finalURL),
	)

	fmt.Println("Error?: ", err)

	if err != nil {
		return fmt.Errorf("chromedp error: %v", err)
	}

	fmt.Println("Final URL: ", finalURL)

	var htmlContent string
	err = chromedp.Run(ctx,
		// chromedp.Sleep(2*time.Second),
		chromedp.InnerHTML("body", &htmlContent),
	)

	if err != nil {
		return fmt.Errorf("error obteniendo HTML: %v", err)
	}

	fmt.Printf("HTML obtenido, tamaño: %d caracteres\n", len(htmlContent))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Errorf("error parseando HTML: %v", err)
	}

	firstGame := doc.Find("#search-results-header ul li").First()
	if firstGame.Length() == 0 {
		fmt.Println("❌ No se encontró ningún juego con el selector #search-results-header ul li")

		firstGame = doc.Find("li.GameCard_search_list__IuMbi").First()
		if firstGame.Length() == 0 {
			return fmt.Errorf("no se encontraron juegos en los resultados")
		}
	}

	gameTitle := strings.TrimSpace(firstGame.Find("h2 a").Text())
	mainStoryLength := strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").First().Text())
	mainExtraLength := strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").Eq(1).Text())
	completionistLength := strings.TrimSpace(firstGame.Find(".GameCard_search_list_details_block__XEXkr .GameCard_search_list_tidbit__0r_OP.center.time_100").Eq(2).Text())

	fmt.Println("✅ Página procesada correctamente")
	fmt.Println("Título del juego: ", gameTitle)
	fmt.Println("Duración de la historia principal: ", mainStoryLength)
	fmt.Println("Duración de la historia principal + extras: ", mainExtraLength)
	fmt.Println("Duración para completarlo al 100%: ", completionistLength)

	return nil

}
