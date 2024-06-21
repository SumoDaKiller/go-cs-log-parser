package main

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var k = koanf.New(".")

func main() {
	if err := k.Load(file.Provider("cs-log-parser.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	fmt.Println("log path is ", k.String("log_path"))
	if err := run(k.String("log_path"), os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(path string, out io.Writer) error {

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

	for i := 0; i < 4; i++ {
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
