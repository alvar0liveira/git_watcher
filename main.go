package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/joho/godotenv"
)

const (
	RAPIDCOMMIT int = iota
	NORMALCOMMIT
	SLOWCOMMIT
)

func main() {

	var (
		path         string
		cadence      int
		author_name  string
		author_email string
	)

	err := godotenv.Load()

	if err != nil {
		log.Fatal(`
		No .env file found

		Please create a .env file with the following content:
		author_name="Your Name"
		author_email="Your Email"
		`)
	}

	author_name = os.Getenv("author_name")
	author_email = os.Getenv("author_email")

	flag.StringVar(&path, "path", ".", "Path to the git repository")
	flag.IntVar(&cadence, "cadence", 1, "Cadence of the commits in minutes. Can be 0 (0.5m), 1 (1m) or 2 (2m).")
	flag.Parse()

	var cadenceDuration time.Duration

	switch cadence {
	case RAPIDCOMMIT:
		cadenceDuration = 30 * time.Second
	case NORMALCOMMIT:
		cadenceDuration = 1 * time.Minute
	case SLOWCOMMIT:
		cadenceDuration = 2 * time.Minute
	default:
		cadenceDuration = 1 * time.Minute
	}

	isRepo, err := isGitRepo(path)
	if err != nil {
		println(err)
	}

	repo := &git.Repository{}

	if isRepo {
		println("It's a git repo")
		repo, err = git.PlainOpen(path)

		if err != nil {
			println(err)
			return
		}
	} else {
		println("It's not a git repo")

		repo, err = InitRepo(path)

		if err != nil {
			println(err)
			return
		}

	}

	quit := make(chan struct{})

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println(err)
				return
			}
			if char == 'q' {
				close(quit)
				return

			}

		}
	}()

	ticker := time.NewTicker(cadenceDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			println("Commiting changes")
			CommitChanges(repo, &git.CommitOptions{
				Author: &object.Signature{
					Name:  author_name,
					Email: author_email,
					When:  time.Now(),
				},
				Committer: &object.Signature{
					Name:  author_name,
					Email: author_email,
					When:  time.Now(),
				},
			})

		case <-quit:
			println("Quitting")
			return
		}
	}
}

func isGitRepo(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	_, err = git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {
		return false, err
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func InitRepo(path string) (*git.Repository, error) {
	println("Creating a git repo")
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func CommitChanges(repo *git.Repository, commitOptions *git.CommitOptions) (commit *object.Commit, err error) {
	workTree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	_, err = workTree.Add(".")

	if err != nil {
		return nil, err
	}

	status, err := workTree.Status()

	if err != nil {
		return nil, err
	}

	if status.IsClean() {
		println("Nothing to commit")
		return nil, err
	}

	commitHash, err := workTree.Commit(time.Now().Format(time.RFC3339), commitOptions)

	if err != nil {
		return nil, err
	}

	commit, err = repo.CommitObject(commitHash)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
