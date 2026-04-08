package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func fileSize(path string) (uint, error) {
	fp, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 0, nil
		}

		return 0, err
	}

	info, err := fp.Stat()
	if err != nil {
		return 0, err
	}

	return uint(info.Size()), nil
}

func createPkgZip(w *zip.Writer, pkgName string, root fs.FS) error {
	return fs.WalkDir(
		root, ".",
		func(name string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if name == "." {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return err
			}
			if !d.IsDir() && !info.Mode().IsRegular() {
				return nil
			}

			h, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			h.Name = pkgName + "/" + name
			if d.IsDir() {
				h.Name += "/"
			}

			h.Method = zip.Deflate
			fw, err := w.CreateHeader(h)
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			f, err := root.Open(name)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(fw, f)
			return err
		})
}

func createZip(resPath, pkgName, outPath string) error {
	outFp, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cannot create: %w", err)
	}
	defer outFp.Close()

	inPath := filepath.Join(resPath, pkgName)
	inDir := os.DirFS(inPath)

	w := zip.NewWriter(outFp)
	if err := createPkgZip(w, pkgName, inDir); err != nil {
		return fmt.Errorf("cannot inflate: %w", err)
	}
	defer w.Close()

	return nil
}

func runCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}
	defer cmd.Wait()

	sc := bufio.NewScanner(pipe)
	for sc.Scan() {
		line := sc.Text()
		fmt.Printf("%s\n", line)
	}

	return nil
}

func runDevel(what, exe string) error {
	outPath := exe + ".dev.exe"
	fmt.Printf("\nbuild >>>\n")
	if err := runCmd("go", "build", "-o", outPath, what); err != nil {
		fmt.Printf("\n<<< build FAIL\n")
		return err
	}
	fmt.Printf("\n<<< build OK\n")
	fmt.Printf("\ndevel >>>\n")
	exePath := fmt.Sprintf(".%c%s", os.PathSeparator, outPath)
	if err := runCmd(exePath); err != nil {
		fmt.Printf("\n<<< devel FAIL\n")
		return err
	}
	fmt.Printf("\n<<< devel OK\n")
	return nil
}

func main() {
	resPath := flag.String("res", "", "resource directory path")
	pkgName := flag.String("pkg", "base", "package name inside of the resource directory")
	flag.Parse()

	if *resPath == "" {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "abort: bad -res path\n")
		os.Exit(1)
	}

	outPath := filepath.Join(*resPath, *pkgName+".zip")

	prevSize, err := fileSize(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "abort: %v", err)
		os.Exit(2)
	}

	if err := createZip(*resPath, *pkgName, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "abort: %v", err)
		os.Exit(3)
	}

	nextSize, err := fileSize(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "abort: %v", err)
		os.Exit(4)
	}

	diffSize := int(nextSize) - int(prevSize)
	fmt.Fprintf(os.Stderr, "%s:\n", outPath)
	fmt.Fprintf(os.Stderr, "%d -> %d (%+d)\n", prevSize, nextSize, diffSize)

	if err := runDevel(".", "rkt"); err != nil {
		fmt.Fprintf(os.Stderr, "abort: %v", err)
		os.Exit(5)
	}
}
