package dev_release

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func PathSize(path string) (uint, error) {
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

func CreateZip(resPath, pkgName, outPath string) error {
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

func runExec(tag, name string, arg ...string) error {
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
		fmt.Printf("%s> %s\n", tag, line)
	}

	return nil
}

func GetOutPath(outName string, isDevel bool) string {
	outPathInfix := ""
	if isDevel {
		outPathInfix = ".dev"
	}

	return outName + outPathInfix + BinExt
}

func Build(what, outName string, isDevel, isFull bool) error {
	arg := []string{"build", "-v"}
	if isFull {
		arg = append(arg, "-a")
	}

	arg = append(arg, "-o", GetOutPath(outName, isDevel), what)

	if err := runExec("build", "go", arg...); err != nil {
		fmt.Println("-- build FAIL --")
		return err
	}
	fmt.Println("-- build OK --")

	return nil
}

/*
	fmt.Println("\n\n\n\n\x1b[7mpress [r] to restart\x1b[0m\x1b[4A")
	b := make([]byte, 1)
	n, _ := os.Stdin.Read(b)
	if n > 1 {
		fmt.Println("commence")
	}
*/

func Run(outName string, isDevel bool) error {
	exePath := fmt.Sprintf(
		".%c%s", os.PathSeparator, GetOutPath(outName, isDevel))
	if err := runExec("devel", exePath); err != nil {
		fmt.Println("-- devel FAIL --")
		return err
	}
	fmt.Println("-- devel OK --")
	return nil
}
