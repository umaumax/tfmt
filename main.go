package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	//	nFlag bool
	tabWidth int
	splitNum int
	delim    string
)

func init() {
	//	flag.BoolVar(&nFlag, "n", false, "line no flag")
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, v := range args {
		f, err := os.Open(v)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			text := scanner.Text()
			if strings.HasPrefix(text, "#") {
				table.Append([]string{text})
				continue
			}

			_delim := "\u001f"
			for _, r := range delim {
				text = strings.Replace(text, string(r), _delim, -1)
			}
			record := regexp.MustCompile(_delim+"+").Split(text, splitNum)
			for i, v := range record {
				//	NOTE if you use " ", tablewriter work wrong
				v = strings.Replace(v, _delim, "_", -1)
				v = strings.Replace(v, " ", "_", -1)
				record[i] = v
			}
			table.Append(record)
		}
		table.Render()
		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}
	}
}
