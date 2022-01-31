package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
)

func generateFormData(nim string) map[string]string {
	return map[string]string{
		"nim": nim,
	}
}

func contain(elements []string, value string) bool {
	for _, element := range elements {
		if element == value {
			return true
		}
	}
	return false
}

func main() {
	csvFile, err := os.Create(fmt.Sprintf("%v - Schedule Sempro.csv", time.Now().Unix()))
	if err != nil {
		log.Fatalln("Error create csv file", err)
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	defer w.Flush()

	titleRow := []string{"NIM"}
	csvRow := [][]string{}

	c := colly.NewCollector(colly.AllowedDomains("tugasakhir.jti.polinema.ac.id"))

	c.OnHTML("table", func(h *colly.HTMLElement) {
		if h.Index == 1 {
			row := []string{}
			row = append(row, h.Response.Ctx.Get("nim"))

			h.ForEach("tr th", func(_ int, h *colly.HTMLElement) {
				if !contain(titleRow, h.Text) {
					titleRow = append(titleRow, h.Text)
				}
			})

			h.ForEach("tr td", func(_ int, el *colly.HTMLElement) {
				row = append(row, el.Text)
			})

			csvRow = append(csvRow, row)
		}
	})

	// c.Limit(&colly.LimitRule{
	// 	Parallelism: 2,
	// 	RandomDelay: 2 * time.Second,
	// })

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting:", r.URL.String())
	// })

	for nim := 1841720001; nim < 1841720250; nim++ {
		strNim := strconv.Itoa(nim)
		fmt.Printf("get data nim: %v\n", strNim)
		c.OnRequest(func(r *colly.Request) {
			r.Ctx.Put("nim", strNim)
		})
		c.Post("http://tugasakhir.jti.polinema.ac.id/mhs.php", generateFormData(strNim))
	}

	c.Wait()

	allRow := [][]string{titleRow}

	allRow = append(allRow, csvRow...)
	w.WriteAll(allRow)
}
