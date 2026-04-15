package bake

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

type SetupInfo struct {
	Project    string // project name
	ResDirPath string // resource directory path
	PkgName    string // package name inside of the resource directory
	VerName    string // release version name
	Platform   string // target platform
}
type SetupFull struct {
	SetupInfo
	OutExePath string // output executable path
	OutPkgPath string // output package path
	OutRelPath string // output release path
}

func (i *SetupInfo) ToFull(outName string, isDevel bool) *SetupFull {
	f := new(SetupFull)
	f.SetupInfo = *i
	f.OutExePath = GetOutPath(outName, isDevel)
	f.OutRelPath = fmt.Sprintf("%s%s_%s.zip", outName, i.VerName, i.Platform)
	f.OutPkgPath = filepath.Join(i.ResDirPath, i.PkgName+".zip")
	return f
}

func NewSetupInfo(project, resDirPath, pkgName, verName, platform string) *SetupInfo {
	s := new(SetupInfo)
	s.Project = project
	s.ResDirPath = resDirPath
	s.PkgName = pkgName
	s.VerName = verName
	s.Platform = platform
	return s
}

func PathSize(path string) (uint, error) {
	fp, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
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

func writePkgZip(w *zip.Writer, setup *SetupInfo, root fs.FS) error {
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

			h.Name = setup.PkgName + "/" + name
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

func CreatePkgZip(setup *SetupFull) error {
	outFp, err := os.Create(setup.OutPkgPath)
	if err != nil {
		return fmt.Errorf("cannot create: %w", err)
	}
	defer outFp.Close()

	inPath := filepath.Join(setup.ResDirPath, setup.PkgName)
	inDir := os.DirFS(inPath)

	w := zip.NewWriter(outFp)
	if err := writePkgZip(w, &setup.SetupInfo, inDir); err != nil {
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
		fmt.Printf("%s %s\n", tag, line)
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

func Build(setup *SetupFull, isFull bool) error {
	arg := []string{"build", "-v"}
	if isFull {
		arg = append(arg, "-a")
	}

	arg = append(arg, "-o", setup.OutExePath, ".")
	if err := runExec("build: ", "go", arg...); err != nil {
		fmt.Println("-- build FAIL --")
		return err
	}

	fmt.Println("-- build OK --")
	return nil
}

type ZipFile struct {
	Name     string
	SaveAs   string
	Optional bool
}

func (z *ZipFile) add(w *zip.Writer) (int64, error) {
	fr, err := os.Open(z.Name)
	if err != nil {
		if z.Optional && errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	defer fr.Close()

	outName := z.Name
	if z.SaveAs != "" {
		outName = z.SaveAs
	}

	fw, err := w.Create(outName)
	if err != nil {
		return 0, err
	}

	written, err := io.Copy(fw, fr)
	if err != nil {
		return 0, err
	}

	return written, nil
}

func writeRelZip(w *zip.Writer, files []*ZipFile) error {
	total := int64(0)
	for _, file := range files {
		written, err := file.add(w)
		if err != nil {
			return err
		}
		total += written
	}

	return nil
}

func CreateRelZip(setup *SetupFull) error {
	outFp, err := os.Create(setup.OutRelPath)
	if err != nil {
		return fmt.Errorf("cannot create: %w", err)
	}
	defer outFp.Close()

	files := []*ZipFile{
		{setup.OutPkgPath, "", false},
		{setup.OutExePath, "", false},
		{"README.md", "", false},
		{"README.html", "", true}}

	w := zip.NewWriter(outFp)
	if err := writeRelZip(w, files); err != nil {
		return fmt.Errorf("cannot inflate: %w", err)
	}
	defer w.Close()

	return nil
}

func Run(setup *SetupFull, isDevel bool) error {
	exePath := fmt.Sprintf(".%c%s", os.PathSeparator, setup.OutExePath)
	if err := runExec("devel: ", exePath); err != nil {
		fmt.Println("-- devel FAIL --")
		return err
	}

	fmt.Println("-- devel OK --")
	return nil
}
