package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	release "szymm.org/rkt/dev/dev_release/src"
)

func main() {
	resPath := flag.String("res", "res", "resource directory path")
	pkgName := flag.String("pkg", "base", "package name inside of the resource directory")
	outName := flag.String("out", "rkt", "output executable name")
	isBuildAll := flag.Bool("all", false, "rebuild all")
	isDevel := flag.Bool("dev", false, "run development build")
	flag.Parse()

	outPath := filepath.Join(*resPath, *pkgName+".zip")

	prevSize, err := release.PathSize(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "abort: %v", err)
		os.Exit(2)
	}

	if err := release.CreateZip(*resPath, *pkgName, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "abort: %v", err)
		os.Exit(3)
	}

	nextSize, err := release.PathSize(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "abort: %v", err)
		os.Exit(4)
	}

	diffSize := int(nextSize) - int(prevSize)
	fmt.Fprintf(os.Stderr, "pkg> %s:\n", outPath)
	fmt.Fprintf(os.Stderr, "pkg> %d -> %d (%+d)\n", prevSize, nextSize, diffSize)
	fmt.Fprintln(os.Stderr, "-- pkg OK --")

	if err := release.Build(".", *outName, *isDevel, *isBuildAll); err != nil {
		fmt.Fprintf(os.Stderr, "abort: %v", err)
		os.Exit(5)
	}

	if *isDevel {
		if err := release.Run(*outName, true); err != nil {
			fmt.Fprintf(os.Stderr, "abort: %v", err)
			os.Exit(6)
		}
	}
}
