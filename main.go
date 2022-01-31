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

func main() {
	csvFile, err := os.Create(fmt.Sprintf("%v - Schedule Sempro.csv", time.Now().Unix()))
	if err != nil {
		log.Fatalln("Error create csv file", err)
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	defer w.Flush()

	titleRow := []string{"NIM", "No.", "Hari, Tanggal", "Waktu (WIB)", "Ruang Daring", "Nama Mahasiswa", "Judul", "Dosen Pembimbing"}
	csvRow := [][]string{}
	csvRow = append(csvRow, titleRow)

	c := colly.NewCollector(colly.AllowedDomains("tugasakhir.jti.polinema.ac.id"))

	c.OnHTML("table", func(h *colly.HTMLElement) {
		if h.Index == 1 {
			row := []string{}
			row = append(row, h.Response.Ctx.Get("nim"))

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
		c.OnRequest(func(r *colly.Request) {
			r.Ctx.Put("nim", strNim)
		})
		c.Post("http://tugasakhir.jti.polinema.ac.id/mhs.php", generateFormData(strNim))
	}

	c.Wait()
	w.WriteAll(csvRow)
}
