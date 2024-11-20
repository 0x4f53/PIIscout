package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"),
		chromedp.Flag("user-data-dir", "./chromeData/"), // Persistent user data directory
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var pageSource string

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://perplexity.ai"),
		chromedp.WaitVisible(`textarea[placeholder*="aanything"]`),
		chromedp.OuterHTML("html", &pageSource),
	); err != nil {
		log.Fatalf("Failed to get page source: %v", err)
	}

	fmt.Println(pageSource)

}
