// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Permission represents a set of permissions
type Permission struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

// Repository represents a repository
type Repository struct {
	ID            int64       `json:"id"`
	Owner         *User       `json:"owner"`
	Name          string      `json:"name"`
	FullName      string      `json:"full_name"`
	Description   string      `json:"description"`
	Empty         bool        `json:"empty"`
	Private       bool        `json:"private"`
	Fork          bool        `json:"fork"`
	Parent        *Repository `json:"parent"`
	Mirror        bool        `json:"mirror"`
	Size          int         `json:"size"`
	HTMLURL       string      `json:"html_url"`
	SSHURL        string      `json:"ssh_url"`
	CloneURL      string      `json:"clone_url"`
	Website       string      `json:"website"`
	Stars         int         `json:"stars_count"`
	Forks         int         `json:"forks_count"`
	Watchers      int         `json:"watchers_count"`
	OpenIssues    int         `json:"open_issues_count"`
	DefaultBranch string      `json:"default_branch"`
	// swagger:strfmt date-time
	Created       time.Time   `json:"created_at"`
	// swagger:strfmt date-time
	Updated       time.Time   `json:"updated_at"`
	Permissions   *Permission `json:"permissions,omitempty"`
}

// ListMyRepos lists all repositories for the authenticated user that has access to.
func (c *Client) ListMyRepos() ([]*Repository, error) {
	repos := make([]*Repository, 0, 10)
	return repos, c.getParsedResponse("GET", "/user/repos", nil, nil, &repos)
}

// ListUserRepos list all repositories of one user by user's name
func (c *Client) ListUserRepos(user string) ([]*Repository, error) {
	repos := make([]*Repository, 0, 10)
	return repos, c.getParsedResponse("GET", fmt.Sprintf("/users/%s/repos", user), nil, nil, &repos)
}

// ListOrgRepos list all repositories of one organization by organization's name
func (c *Client) ListOrgRepos(org string) ([]*Repository, error) {
	repos := make([]*Repository, 0, 10)
	return repos, c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/repos", org), nil, nil, &repos)
}

// CreateRepoOption options when creating repository
// swagger:model
type CreateRepoOption struct {
	// Name of the repository to create
	//
	// required: true
	// unique: true
	Name string `json:"name" binding:"Required;AlphaDashDot;MaxSize(100)"`
	// Description of the repository to create
	Description string `json:"description" binding:"MaxSize(255)"`
	// Whether the repository is private
	Private bool `json:"private"`
	// Whether the repository should be auto-intialized?
	AutoInit bool `json:"auto_init"`
	// Gitignores to use
	Gitignores string `json:"gitignores"`
	// License to use
	License string `json:"license"`
	// Readme of the repository to create
	Readme string `json:"readme"`
}

// CreateRepo creates a repository for authenticated user.
func (c *Client) CreateRepo(opt CreateRepoOption) (*Repository, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("POST", "/user/repos", jsonHeader, bytes.NewReader(body), repo)
}

// CreateOrgRepo creates an organization repository for authenticated user.
func (c *Client) CreateOrgRepo(org string, opt CreateRepoOption) (*Repository, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("POST", fmt.Sprintf("/org/%s/repos", org), jsonHeader, bytes.NewReader(body), repo)
}

// GetRepo returns information of a repository of given owner.
func (c *Client) GetRepo(owner, reponame string) (*Repository, error) {
	repo := new(Repository)
	return repo, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s", owner, reponame), nil, nil, repo)
}

// DeleteRepo deletes a repository of user or organization.
func (c *Client) DeleteRepo(owner, repo string) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s", owner, repo), nil, nil)
	return err
}

// MigrateRepoOption options for migrating a repository from an external service
type MigrateRepoOption struct {
	// required: true
	CloneAddr string `json:"clone_addr" binding:"Required"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	// required: true
	UID int `json:"uid" binding:"Required"`
	// required: true
	RepoName string `json:"repo_name" binding:"Required"`
	Mirror bool `json:"mirror"`
	Private bool `json:"private"`
	Description string `json:"description"`
}

// MigrateRepo migrates a repository from other Git hosting sources for the
// authenticated user.
//
// To migrate a repository for a organization, the authenticated user must be a
// owner of the specified organization.
func (c *Client) MigrateRepo(opt MigrateRepoOption) (*Repository, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("POST", "/repos/migrate", jsonHeader, bytes.NewReader(body), repo)
}

// MirrorSync adds a mirrored repository to the mirror sync queue.
func (c *Client) MirrorSync(owner, repo string) error {
	_, err := c.getResponse("POST", fmt.Sprintf("/repos/%s/%s/mirror-sync", owner, repo), nil, nil)
	return err
}
