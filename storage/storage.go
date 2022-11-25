// TODO: Create API for other services to use

package storage

import (
	"database/sql"

	"github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/andreaswachs/lazyworkflows/model/response"
)

type Action uint8

const (
	StoreWorkflow Action = iota
	UpdateWorkflow
	RemoveWorkflow
)

// Message conatins the data for each message to the storage service
type Message struct {
	Action          Action
	Workflow        response.Workflow
	WorkflowId      string // Used for updating and removing workflows
	responseChannel chan interface{}
}

// Start starts the storage service
func Start(msgs chan Message) error {
	// TODO: would it be better to save the DB to a file between sessions?
	db, err := sql.Open("sqlite", ":memory:")

	// Whoa cowboy, this is opposite of what we usually do!
	// If there are no errors, we want to continue to listen for messages to the storage service
	if err == nil {
		go listen(db, msgs)
	}

	return err
}

// Listen listens for messages and responds accordingly
func listen(db *sql.DB, msgs chan Message) {
	for {
		msg := <-msgs
		switch msg.Action {
		case StoreWorkflow:
			insertWorkflow(db, msg.Workflow)
			continue
		case UpdateWorkflow:
			updateWorkflow(db, msg.Workflow)
			continue
		case RemoveWorkflow:
			continue

		}
	}
}

func insertRepo(db *sql.DB, repo appconfig.Repo) error {
	insertRepo := `INSERT INTO repository (name, owner, token) VALUES (?, ?, ?)`
	_, err := db.Exec(insertRepo, repo.Repo, repo.Owner, repo.Token)

	return err
}

func insertWorkflow(db *sql.DB, workflow response.Workflow) error {
	insertWorkflow := `INSERT INTO workflow (id, node_id, name, path, state, created_at, updated_at, url, html_url, badge_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(insertWorkflow, workflow.Id, workflow.NodeId, workflow.Name, workflow.Path, workflow.State, workflow.CreatedAt, workflow.Url, workflow.UpdatedAt, workflow.HtmlUrl, workflow.BadgeUrl)

	return err
}

func updateWorkflow(db *sql.DB, workflow response.Workflow) error {
	updateWorkflow := `UPDATE workflow SET id=?, node_id=?, name=?, path=?, state=?, created_at=?, updated_at=?, url=?, html_url=?, badge_url=? WHERE id=?`
	_, err := db.Exec(updateWorkflow, workflow.Id, workflow.NodeId, workflow.Name, workflow.Path, workflow.State, workflow.CreatedAt, workflow.Url, workflow.UpdatedAt, workflow.HtmlUrl, workflow.BadgeUrl, workflow.Id)

	return err
}

func deleteWorkflow(db *sql.DB, workflowId string) error {
	deleteWorkflow := `DELETE FROM workflow WHERE id=?`
	_, err := db.Exec(deleteWorkflow, workflowId)

	return err
}

func getWorkflowsFromRepo(db *sql.DB, repo appconfig.Repo) ([]response.Workflow, error) {
	getWorkflows := `SELECT * FROM workflows WHERE repo=?`
	rows, err := db.Query(getWorkflows, repo.Repo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []response.Workflow
	for rows.Next() {
		var workflow response.Workflow
		err = rows.Scan(&workflow.Id, &workflow.NodeId, &workflow.Name, &workflow.Path, &workflow.State, &workflow.CreatedAt, &workflow.Url, &workflow.UpdatedAt, &workflow.HtmlUrl, &workflow.BadgeUrl)
		if err != nil {
			return nil, err
		}

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

func getWorkflow(db *sql.DB, workflowId string) (response.Workflow, error) {
	var workflow response.Workflow
	getWorkflow := `SELECT * FROM workflows WHERE id=?`
	err := db.QueryRow(getWorkflow, workflowId).Scan(&workflow.Id, &workflow.NodeId, &workflow.Name, &workflow.Path, &workflow.State, &workflow.CreatedAt, &workflow.Url, &workflow.UpdatedAt, &workflow.HtmlUrl, &workflow.BadgeUrl)
	if err != nil {
		return workflow, err
	}

	return workflow, nil
}

func getRepos(db *sql.DB) ([]appconfig.Repo, error) {
	getRepos := `SELECT * FROM repository`
	rows, err := db.Query(getRepos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []appconfig.Repo
	for rows.Next() {
		var repo appconfig.Repo
		err = rows.Scan(&repo.Repo, &repo.Owner, &repo.Token)
		if err != nil {
			return nil, err
		}

		repos = append(repos, repo)
	}

	return repos, nil
}

func createDb(db *sql.DB) error {
	createTables := `CREATE TABLE "workflows" (
		"id"	TEXT NOT NULL,
		"node_id"	TEXT NOT NULL,
		"name"	TEXT NOT NULL,
		"path"	TEXT NOT NULL,
		"state"	TEXT NOT NULL,
		"created_at"	TEXT NOT NULL,
		"updated_at"	TEXT NOT NULL,
		"url"	TEXT NOT NULL,
		"html_url"	TEXT NOT NULL,
		"badge_url"	TEXT NOT NULL,
		PRIMARY KEY("id")
	);
	
	CREATE TABLE "repository" (
	"name"	TEXT NOT NULL,
	"owner"	TEXT NOT NULL,
	"token"	TEXT NOT NULL,
	PRIMARY KEY("name")
	);`
	_, err := db.Exec(createTables)

	return err
}
