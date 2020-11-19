package netrc

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	bindingFilepath = "/platform/bindings/netrc/secret/.netrc"
	netrcFile       = ".netrc"
)

func Setup() (bool, error) {
	done, err := fromBinding()
	if err != nil {
		return false, fmt.Errorf("from binding: %w", err)
	}
	if done {
		return true, nil
	}

	done, err = fromEnvironment()
	if err != nil {
		return false, fmt.Errorf("from binding: %w", err)
	}
	if done {
		return true, nil
	}

	return false, nil
}

func fromBinding() (bool, error) {
	if _, err := os.Stat(bindingFilepath); os.IsNotExist(err) {
		return false, nil
	}

	if err := cp(bindingFilepath, netrcFilepath()); err != nil {
		return false, fmt.Errorf("cp: %w", err)
	}

	return true, nil
}

func fromEnvironment() (bool, error) {
	var (
		machine  = os.Getenv("NETRC_MACHINE")
		login    = os.Getenv("NETRC_LOGIN")
		password = os.Getenv("NETRC_PASSWORD")
	)

	if machine == "" || login == "" || password == "" {
		return false, nil
	}

	content := strings.Join([]string{
		"machine " + machine,
		"login " + login,
		"password" + password,
	}, "\n")

	fpath := netrcFilepath()
	if err := ioutil.WriteFile(fpath, []byte(content), 0644); err != nil {
		return false, fmt.Errorf("write file '%s': %w", fpath, err)
	}

	return true, nil
}

func netrcFilepath() string {
	return filepath.Join(os.Getenv("HOME"), ".netrc")
}

func cp(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}
