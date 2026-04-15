package main

import (
	"flag"
	"fmt"
	"os"

	bake "szymm.org/rkt/dev/src"
)

func abort(err error, exit int) {
	fmt.Fprintf(os.Stderr, "abort: %v\n", err)
	os.Exit(exit)
}

func info(f string, a ...any) {
	fmt.Fprintf(os.Stderr, f, a...)
}

func main() {
	resDirPath := flag.String("res", "res", "resource directory path")
	outName := flag.String("out", "rkt", "output executable and release prefix (<out>.exe, <out>v9_win64.zip)")
	pkgName := flag.String("pkg", "base", "package name inside of the resource directory")
	verName := flag.String("ver", "", "release version name (ex. 'v9')")
	platform := flag.String("plat", "", "release platform (ex. 'win64')")
	isBuildAll := flag.Bool("all", false, "rebuild all")
	isRelease := flag.Bool("rel", false, "build release zip")
	isDevel := flag.Bool("dev", false, "run development build")
	flag.Parse()

	if *isDevel && *isRelease {
		abort(fmt.Errorf("cannot run -dev and -rel together"), 1)
	}
	if *isRelease {
		if *verName == "" {
			abort(fmt.Errorf("no release version name"), 1)
		}
		if *platform == "" {
			abort(fmt.Errorf("no release platform"), 1)
		}
	}

	setupInfo := bake.NewSetupInfo("rkt", *resDirPath, *pkgName, *verName, *platform)
	setup := setupInfo.ToFull(*outName, *isDevel)

	prevSize, err := bake.PathSize(setup.OutPkgPath)
	if err != nil {
		abort(err, 2)
	}
	if err := bake.CreatePkgZip(setup); err != nil {
		abort(err, 3)
	}

	nextSize, err := bake.PathSize(setup.OutPkgPath)
	if err != nil {
		abort(err, 4)
	}

	diffSize := int(nextSize) - int(prevSize)
	info("pkg: %s:\n", setup.OutPkgPath)
	info("pkg: %d -> %d (%+d)\n", prevSize, nextSize, diffSize)
	info("-- pkg OK --\n")

	if err := bake.Build(setup, *isBuildAll); err != nil {
		abort(err, 5)
	}

	if *isDevel {
		if err := bake.Run(setup, true); err != nil {
			abort(err, 6)
		}
	}
	if *isRelease {
		if err := bake.CreateRelZip(setup); err != nil {
			abort(err, 7)
		}
		info("-- release OK --\n")
	}
}
