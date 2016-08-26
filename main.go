package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	joinFlag bool
	lineFlag bool
	tabWidth int
	splitNum int
	delim    string
)

func init() {
	flag.BoolVar(&joinFlag, "join", false, "join files?")
	flag.BoolVar(&lineFlag, "line", false, "line flag")
	flag.IntVar(&tabWidth, "tab", 4, "tab width")
	flag.IntVar(&splitNum, "n", -1, "split num if -1 all")
	flag.StringVar(&delim, "delim", " \t,", "delim")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if flag.NArg() == 0 {
		args = []string{"/dev/stdin"}
	}
	if joinFlag {
		join(args)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	_delim := "\u001f"
	delimReg := regexp.MustCompile(_delim + "+")
	n := 0
	for _, v := range args {
		func() {
			f, err := os.Open(v)
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				n++
				text := scanner.Text()

				var record []string
				if lineFlag {
					record = []string{fmt.Sprintf("%d", n)}
				}

				if strings.HasPrefix(text, "#") {
					record = append(record, text)
					table.Append(record)
					continue
				}

				for _, r := range delim {
					text = strings.Replace(text, string(r), _delim, -1)
				}
				record = append(record, delimReg.Split(text, splitNum)...)
				for i, v := range record {
					//	NOTE if you use " ", tablewriter work wrong
					v = strings.Replace(v, _delim, "_", -1)
					v = strings.Replace(v, " ", "_", -1)
					record[i] = v
				}
				table.Append(record)
			}
			if err := scanner.Err(); err != nil {
				log.Fatalln(err)
			}
		}()
	}
	table.Render()
}
