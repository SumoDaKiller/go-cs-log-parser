package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go-cs-log-parser/database"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var k = koanf.New(".")
var queries *database.Queries
var ctx context.Context

func main() {
	confFile := flag.String("c", "cs-log-parser.yaml", "configuration file")
	parseOnly := flag.Bool("p", false, "Only parse logfile and skip generating output HTML files")
	outputOnly := flag.Bool("o", false, "Only output HTML files and skip parsing log files")
	startServer := flag.Bool("s", false, "Start API server. -o and -p will be ignored when -s is used")

	flag.Parse()

	if err := k.Load(file.Provider(*confFile), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	fmt.Println("log path is ", k.String("log_path"))

	// If we want to start the server, we ignore output and parse parameters
	if *startServer {
		if err := runServer(); err != nil {
			log.Fatal(err)
		}
	} else {
		if *parseOnly && *outputOnly {
			log.Fatal("-o and -p cannot be used at the same time")
		}

		if *parseOnly || !*outputOnly {
			if err := runParser(k.String("log_path"), os.Stdout); err != nil {
				log.Fatal(err)
			}
		}

		if *outputOnly || !*parseOnly {
			if err := runOutput(); err != nil {
				log.Fatal(err)
			}
		}
	}

}

func runOutput() error {
	err := generateAllPages()
	if err != nil {
		return err
	}
	return nil
}

func runParser(path string, out io.Writer) error {

	errCh := make(chan error)
	doneCh := make(chan struct{})
	filesCh := make(chan string)

	wg := sync.WaitGroup{}

	go func() {
		defer close(filesCh)
		err := filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {

					return err
				}
				if filepath.Ext(path) == ".log" {
					fmt.Printf("Now parsing log file: %s\n", path)
					filesCh <- path
				}
				return nil
			})
		if err != nil {
			fmt.Printf("filepath.Walk() returned %v\n", err)
			errCh <- err
		}
	}()

	re := initializeRegexpPatterns()
	// TODO: Only use one thread for now as there are some table that might be written to with the same data in each thread, and then hits unique errors
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fname := range filesCh {
				f, err := os.Open(fname)
				if err != nil {
					errCh <- fmt.Errorf("cannot open file: %w", err)
					return
				}
				err = parseFile(f, re)
				if err != nil {
					errCh <- fmt.Errorf("cannot parse file: %w", err)
				}
				if err := f.Close(); err != nil {
					errCh <- err
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case <-doneCh:
			return nil
		}
	}
}
