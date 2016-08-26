package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

func join(args []string) {
	var err error
	fs := make([]*os.File, len(args))
	scanners := make([]*bufio.Scanner, len(args))
	for i, v := range args {
		fs[i], err = os.Open(v)
		if err != nil {
			log.Fatalln(err)
		}
		defer fs[i].Close()
		scanner := bufio.NewScanner(fs[i])
		scanners[i] = scanner
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	func() {
		record := make([]string, len(scanners))
		offset := 0
		if lineFlag {
			record = append(record, "")
			offset++
		}
		for n := 1; ; n++ {
			for i, scanner := range scanners {
				if lineFlag {
					record[0] = fmt.Sprintf("%d", n)
				}
				if scanner.Scan() {
					text := scanner.Text()
					//					if i == 0 {
					//						if lineFlag {
					//							fmt.Printf("%d%s", i, delim)
					//						}
					//					}
					//					if i > 0 {
					//						fmt.Printf("%s", delim)
					//					}
					//					fmt.Printf("%s", text)
					record[offset+i] = text
				} else {
					return
				}
			}
			table.Append(record)
			//			fmt.Println()
		}
	}()
	table.Render()
}
