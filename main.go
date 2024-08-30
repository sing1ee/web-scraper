package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/scrape", func(c *fiber.Ctx) error {
		var input struct {
			URL string `json:"url"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}

		// 使用go-rod加载URL
		browser := rod.New().MustConnect()
		defer browser.MustClose()

		page := browser.MustPage(input.URL)
		page.MustWaitStable()

		// 获取页面HTML
		html := page.MustHTML()

		// 使用goquery处理HTML
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse HTML"})
		}

		// 移除JavaScript和CSS
		doc.Find("script").Remove()
		doc.Find("style").Remove()
		doc.Find("link[rel='stylesheet']").Remove()

		// 获取处理后的HTML
		cleanedHTML, err := doc.Html()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate cleaned HTML"})
		}
		// save to file with filename set to md5(url)
		filename := md5.Sum([]byte(input.URL))
		os.WriteFile(fmt.Sprintf("%x.html", filename), []byte(cleanedHTML), 0644)
		return c.JSON(fiber.Map{"html": cleanedHTML})
	})

	app.Listen(":3456")
}
