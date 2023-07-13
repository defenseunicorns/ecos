package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	goyaml "github.com/goccy/go-yaml"
	"github.com/mholt/archiver/v3"
	"github.com/otiai10/copy"
)

func Archive(sourceDir string, archiveName string) error {
	archiveSrc := []string{sourceDir + string(os.PathSeparator)}

	fmt.Printf("Creating archive %s\n", archiveName)

	if err := archiver.Archive(archiveSrc, archiveName); err != nil {
		return fmt.Errorf("Unable to create archive: %w", err)
	}

	_, err := os.Stat(archiveName)
	if err != nil {
		return fmt.Errorf("Unable to read the package archive: %w", err)
	}

	return nil
}

func Copy(source string, destination string) error {
	return filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		fmt.Printf("Copying %s\n", path)

		return copy.Copy(
			source,
			destination,
			copy.Options{
				Skip: func(info os.FileInfo, src, dest string) (bool, error) {
					return strings.HasSuffix(src, ".git-like"), nil
				},
				OnSymlink: func(src string) copy.SymlinkAction {
					return copy.Skip
				},
				PermissionControl: copy.AddPermission(0644),
			},
		)
	})
}

func CreateDirectory(path string, mode os.FileMode) error {
	if InvalidPath(path) {
		return os.MkdirAll(path, mode)
	}
	return nil
}

func ExecCommand(command string, envVars []string, args ...string) error {
	var (
		stdoutBuf, stderrBuf bytes.Buffer
		wg                   sync.WaitGroup
		stdoutErr, stderrErr error
	)

	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), envVars...)

	cmdStdout, _ := cmd.StdoutPipe()
	cmdStderr, _ := cmd.StderrPipe()

	stdoutWriters := []io.Writer{
		&stdoutBuf,
		os.Stdout,
	}

	stderrWriters := []io.Writer{
		&stderrBuf,
		os.Stderr,
	}

	stdout := io.MultiWriter(stdoutWriters...)
	stderr := io.MultiWriter(stderrWriters...)

	fmt.Printf("Executing command: %s %s\n", command, strings.Join(args, " "))

	if err := cmd.Start(); err != nil {
		return err
	}

	wg.Add(2)

	go func() {
		_, stdoutErr = io.Copy(stdout, cmdStdout)
		wg.Done()
	}()

	go func() {
		_, stderrErr = io.Copy(stderr, cmdStderr)
		wg.Done()
	}()

	wg.Wait()

	if stdoutErr != nil {
		return stdoutErr
	}

	if stderrErr != nil {
		return stderrErr
	}

	fmt.Println()

	return nil
}

func InvalidPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsPermission(err) && err != nil
}

func MakeTempDir(tmpDir string) (string, error) {
	// Create the base tmp directory if it is specified.
	if tmpDir != "" {
		if err := CreateDirectory(tmpDir, 0700); err != nil {
			return "", err
		}
	}

	tmp, err := os.MkdirTemp(tmpDir, "ecos-")
	fmt.Printf("Using temp path: '%s'\n", tmp)
	return tmp, err
}

func PrintYaml(srcSpec any) error {
	content, err := goyaml.Marshal(srcSpec)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n\n", content)
	return nil
}

func ReadYaml(path string, spec any) error {
	fmt.Printf("Loading ecos spec %s\n", path)

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := goyaml.Unmarshal(file, spec); err != nil {
		return err
	}

	PrintYaml(spec)
	return nil
}

func WriteFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create the file at %s to write the contents: %w", path, err)
	}

	_, err = f.Write(data)
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("unable to write the file at %s contents:%w", path, err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("error saving file %s: %w", path, err)
	}

	return nil
}

func WriteYaml(path string, srcSpec any, perm fs.FileMode) error {
	fmt.Printf("Writing spec to %s\n", path)
	PrintYaml(srcSpec)

	// Save the parsed output to the config path given
	content, err := goyaml.Marshal(srcSpec)
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, perm)
}
