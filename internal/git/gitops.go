package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type GitOps struct {
	RepoPath string
	RepoURL  string
}

func (g *GitOps) Clone() error {
	if _, err := os.Stat(g.RepoPath); os.IsNotExist(err) {
		err = os.MkdirAll(g.RepoPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	// Clone the repository
	_, err := git.PlainClone(g.RepoPath, false, &git.CloneOptions{
		URL: g.RepoURL,
	})
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		return err
	}

	return nil
}

func (g *GitOps) Pull() error {
	// Open the repository
	repo, err := git.PlainOpen(g.RepoPath)
	if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Pull the latest changes from the origin remote
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		return err
	}

	return nil
}

// commit changes
func (g *GitOps) Commit(msg string) error {

	// Open the repository
	repo, err := git.PlainOpen(g.RepoPath)
	if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Add all files to the staging area
	_, err = w.Add(".")
	if err != nil {
		return err
	}

	// Commit the changes to the repository
	_, err = w.Commit(msg, &git.CommitOptions{})
	if err != nil {
		return err
	}

	return nil
}

// push changes
func (g *GitOps) Push() error {
	// Open the repository
	repo, err := git.PlainOpen(g.RepoPath)
	if err != nil {
		return err
	}

	// Push the changes to the origin remote
	err = repo.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	return nil
}
