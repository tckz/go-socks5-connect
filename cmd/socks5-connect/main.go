package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/proxy"
	"golang.org/x/sync/errgroup"
)

var version string

var (
	optSocks5  = flag.String("socks5", "", "socks5 host:port")
	optDest    = flag.String("dest", "", "destination host:port")
	optVersion = flag.Bool("version", false, "show version")
)

func main() {
	flag.Parse()

	if *optVersion {
		fmt.Println(version)
		return
	}

	if *optSocks5 == "" {
		log.Fatal("--socks5 must be specified")
	}
	if *optDest == "" {
		log.Fatal("--dest must be specified")
	}

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	dialer, err := proxy.SOCKS5("tcp", *optSocks5, nil, proxy.Direct)
	if err != nil {
		return fmt.Errorf("proxy.SOCKS5: %w", err)
	}

	conn, err := dialer.Dial("tcp", *optDest)
	if err != nil {
		return fmt.Errorf("dialer.Dial: %w", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	relay := func(ctx context.Context, name string, bufSize int, r io.Reader, w io.Writer) func() error {
		return func() (retErr error) {
			defer func() {
				if retErr != nil {
					cancel()
					retErr = fmt.Errorf("%s: %s", name, retErr)
				}
			}()
			buf := make([]byte, bufSize)
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					n, err := r.Read(buf)
					if n > 0 {
						_, err := w.Write(buf[0:n])
						if err != nil {
							return err
						}
					}
					if err != nil {
						return err
					}
					if n == 0 {
						return nil
					}
				}
			}
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(relay(ctx, "upstream", 4096, os.Stdin, conn))
	eg.Go(relay(ctx, "downstream", 4096, conn, os.Stdout))

	if err := eg.Wait(); err == io.EOF {
		return nil
	}
	return err
}
