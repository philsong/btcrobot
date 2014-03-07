// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package huobiapi

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"logger"
	"os"
	"strconv"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string, skipline int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for i, line := range lines {
		if i < skipline {
			continue
		}

		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func data2csv(filename string, skipline int) {
	lines, err := readLines(filename)
	if err != nil {
		logger.Fatalf("readLines: %s", err)
	}
	/*
		for i, line := range lines {
			fmt.Println(i, line)
		}
	*/

	if err := writeLines(lines, filename+".csv", skipline); err != nil {
		logger.Fatalf("writeLines: %s", err)
	}
}

func ParseMinuteCSV(filename string) (MinuteRecords []MinuteRecord) {
	data2csv(filename, 3)
	file, err := os.Open(filename + ".csv")
	if err != nil {
		fmt.Println("ParseMinuteCSV Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	/*
		record, err := reader.ReadAll()
		fmt.Println(record)
		return
	*/

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(record) == 0 {
			fmt.Println("Error:", "record is zero, maybe it is not a cvs format!!!")
			return
		}

		var minRecord MinuteRecord
		minRecord.Time = record[0]
		minRecord.Price, err = strconv.ParseFloat(record[1], 64)
		if err != nil {
			logger.Fatalln(record)
			logger.Fatal("ParseMinuteCSV item price is not number")
		}
		minRecord.Volumn, err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			logger.Fatal("ParseMinuteCSV item Volumn is not number")
		}
		minRecord.Amount, err = strconv.ParseFloat(record[3], 64)
		if err != nil {
			logger.Fatal("ParseMinuteCSV item Amount is not number")
		}

		MinuteRecords = append(MinuteRecords, minRecord)
	}

	return
}

func ParsePeroidCSV(filename string) (PeroidRecords []PeroidRecord) {
	data2csv(filename, 2)
	file, err := os.Open(filename + ".csv")
	if err != nil {
		fmt.Println("ParsePeroidCSV Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	/*
		record, err := reader.ReadAll()
		fmt.Println(record)
		return
	*/

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(record) < 8 {
			fmt.Println("Error:", "record is zero, maybe it is not a cvs format!!!", len(record))
			return
		}

		var peroidRecord PeroidRecord
		peroidRecord.Date = record[0]
		peroidRecord.Time = record[1]
		peroidRecord.Open, err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		peroidRecord.High, err = strconv.ParseFloat(record[3], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		peroidRecord.Low, err = strconv.ParseFloat(record[4], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}
		peroidRecord.Close, err = strconv.ParseFloat(record[5], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item price is not number")
		}

		peroidRecord.Volumn, err = strconv.ParseFloat(record[6], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item Volumn is not number")
		}
		peroidRecord.Amount, err = strconv.ParseFloat(record[7], 64)
		if err != nil {
			logger.Fatal("ParsePeroidCSV item Amount is not number")
		}

		PeroidRecords = append(PeroidRecords, peroidRecord)
	}
	return
}
