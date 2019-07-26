package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Library struct {
	url     string
	dir     string
	name    string
	hosting *GitLabHosting
}

func NewLibrary(libsDir, name string) *Library {
	dir := filepath.Join(libsDir, name)
	url := "https://" + name
	return &Library{url: url, dir: dir, name: name, hosting: NewGitLabHosting(name)}
}

func (l *Library) Prepare() error {
	if _, err := os.Stat(l.dir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(l.dir, 0700)
			if err != nil {
				return err
			}
			return l.cloneRepository()
		} else {
			return err
		}
	}
	if err := l.fetchRepository(); err != nil {
		return err
	}

	return nil
}

func (l *Library) UpdateDependency(dependency string, version string) (string, error) {
	dependent, err := l.isDependent(dependency)
	if err != nil {
		return "", err
	}

	if dependent {
		fmt.Println(l.name, "depending", dependency)
	} else {
		return "", nil
	}

	branchName, err := l.makeChanges(dependency, version)
	if err != nil {
		return "", err
	}

	title := l.changesMessage(dependency, version)
	if url, err := l.hosting.CreateMR(title, branchName, "master"); err != nil {
		return "", err
	} else {
		return url, nil
	}
}

func (l *Library) isDependent(dependency string) (bool, error) {
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	cmd.Dir = l.dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var module ModulePublic
		if err := decoder.Decode(&module); err != nil {
			if err == io.EOF {
				break
			}
			return false, err
		}

		if module.Path == dependency {
			return !module.Main && !module.Indirect, nil
		}
	}

	return false, nil
}

func (l *Library) makeChanges(dependency, version string) (string, error) {
	branchName := l.branchName(dependency, version)
	if err := l.createBranch(branchName); err != nil {
		return "", err
	}

	if err := l.updateVersion(dependency, version); err != nil {
		return "", err
	}

	changesMessage := l.changesMessage(dependency, version)
	if err := l.commitChanges(changesMessage); err != nil {
		return "", err
	}

	if err := l.pushBranch(branchName); err != nil {
		return "", err
	}

	return branchName, nil
}

func (l *Library) branchName(dependency, version string) string {
	now := time.Now()
	return fmt.Sprintf("update_%s_to_%s_%s", dependency, version, now.Format("20060102150405"))
}

func (l *Library) changesMessage(dependency, version string) string {
	return fmt.Sprintf("update %s to %s", dependency, version)
}

func (l *Library) execCommand(libraryDir bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if libraryDir {
		cmd.Dir = l.dir
	}
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Println(string(out))
	}
	return err
}

func (l *Library) cloneRepository() error {
	fmt.Println("cloning:", l.name)
	return l.execCommand(false, "git", "clone", l.url, l.dir)
}

func (l *Library) fetchRepository() error {
	fmt.Println("fetching:", l.name)
	return l.execCommand(true, "git", "fetch")
}

func (l *Library) createBranch(branch string) error {
	return l.execCommand(true, "git", "checkout", "-b", branch, "origin/master")
}

func (l *Library) pushBranch(branch string) error {
	return l.execCommand(true, "git", "push", "-u", "origin", branch)
}

func (l *Library) commitChanges(message string) error {
	err := l.execCommand(true, "git", "add", "go.mod", "go.sum")
	if err != nil {
		return err
	}
	return l.execCommand(true, "git", "commit", "-m", message)
}

func (l *Library) updateVersion(dependency, version string) error {
	return l.execCommand(true, "go", "get", "-d", fmt.Sprintf("%s@%s", dependency, version))
}
