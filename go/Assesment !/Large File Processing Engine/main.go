package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

var mu sync.Mutex
var fmu sync.Mutex

type Failed_record struct {
	name string
	rec  []string
	err  string
}

type safe_strct_failed_record struct {
	fmu           sync.Mutex
	failed_record []Failed_record
}

type Department_stats struct {
	emp_cnt      int64
	avg_salaries float64
}

func main() {
	ctx := context.Background()
	add := "emp_recod.csv"
	FileProcessor(ctx, add)
}

func FileProcessor(ctx context.Context, address string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan []string, 100)
	stats := make(map[string]Department_stats)
	var failed_rec []Failed_record
	var wg sync.WaitGroup

	wg.Add(8)
	for i := 0; i < 8; i++ {
		go processRecord(ctx, jobs, stats, &wg, &failed_rec)
	}

	readfile(ctx, cancel, address, jobs)
	close(jobs)
	wg.Wait()

	printstats(stats)
	printfailed(failed_rec)
}

func processRecord(ctx context.Context, jobs <-chan []string, stats map[string]Department_stats, wg *sync.WaitGroup, failed_rec *[]Failed_record) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case rec, ok := <-jobs:
			if !ok {
				return
			}

			if !validateRecord(rec, failed_rec) {
				continue
			}

			updatestats(stats, rec)
		}
	}
}

func updatestats(stats map[string]Department_stats, rec []string) {
	// find department in rec
	dep := rec[7]
	sal := rec[4]
	salFloat, err := strconv.ParseFloat(sal, 64)
	if err != nil {
		return
	}

	mu.Lock()
	stat := stats[dep]
	stat.emp_cnt++
	stat.avg_salaries += salFloat
	stats[dep] = stat
	mu.Unlock()
}

func readfile(ctx context.Context, cancel context.CancelFunc, add string, jobs chan<- []string) {
	// Implement your file processing logic here
	file, err := os.Open(add)
	if err != nil {
		log.Printf("failed to open file %s: %v\n", add, err)
		cancel()
		return
	}
	defer file.Close()

	csv_read := csv.NewReader(file)

	// Read and skip header row
	header, err := csv_read.Read()

	if err != nil {
		if err == io.EOF {
			return
		}
		log.Printf("failed to read header: %v\n", err)
		cancel()
		return
	}
	_ = header

	for {

		record, err := csv_read.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("error reading record: %v\n", err)
			cancel()
			break
		}

		select {
		case <-ctx.Done():
			return
		case jobs <- record:
		}
	}
}

func validateRecord(rec []string, failed_rec *[]Failed_record) bool {
	if len(rec) < 8 {
		fmu.Lock()
		*failed_rec = append(*failed_rec, Failed_record{
			name: "unknown",
			rec:  rec,
			err:  "insufficient columns",
		})
		fmu.Unlock()
		return false
	}

	name := rec[0]
	dep := rec[7]
	sal := rec[4]

	if dep == "" {
		fmu.Lock()
		*failed_rec = append(*failed_rec, Failed_record{
			name: name,
			rec:  rec,
			err:  "missing department",
		})
		fmu.Unlock()
		return false
	}

	if _, err := strconv.ParseFloat(sal, 64); err != nil {
		fmu.Lock()
		*failed_rec = append(*failed_rec, Failed_record{
			name: name,
			rec:  rec,
			err:  "invalid salary: " + err.Error(),
		})
		mu.Unlock()
		return false
	}

	return true
}

func printstats(stats map[string]Department_stats) {
	fmt.Printf("\n%-24s | %-10s | %-15s | %-15s\n", "Department", "Employees", "Total Salary", "Average Salary")
	fmt.Println("-------------------------------------------------------------------------")
	for dep, stat := range stats {
		avg := 0.0
		if stat.emp_cnt > 0 {
			avg = stat.avg_salaries / float64(stat.emp_cnt)
		}
		fmt.Printf("%-24s | %-10d | $%-14.2f | $%-14.2f\n", dep, stat.emp_cnt, stat.avg_salaries, avg)
	}
}

func printfailed(failed_rec []Failed_record) {
	fmt.Printf("\n--- Failed Records (Total: %d) ---\n", len(failed_rec))
	for i, f := range failed_rec {
		if i >= 10 {
			fmt.Printf("... and %d more failed records\n", len(failed_rec)-10)
			break
		}
		fmt.Printf("Name: %-15s | Error: %-22s | Record: %v\n", f.name, f.err, f.rec)
	}
}
