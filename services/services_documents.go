package services

import (
	"time"
)

//Document document model
type Document struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        int
	ProjectID int
	Title     string
	SortOrder int
	Body      string
}

//CreateDocument --
func CreateDocument(projectID int, title string, body string, sortOrder int) (document *Document, err error) {
	document = &Document{
		ProjectID: projectID,
		Title:     title,
		Body:      body,
		SortOrder: sortOrder,
	}

	err = DB.QueryRow("insert into documents (project_id, title, body, sort_order) values($1,$2,$3,$4) returning created_at,updated_at,id", projectID, title, body, sortOrder).Scan(&document.CreatedAt, &document.UpdatedAt, &document.ID)
	return
}

//LoadDocuments loads all the documents of a project
func LoadDocuments(projectID int) ([]*Document, error) {
	rows, err := DB.Query("select created_at, updated_at, id, title, sort_order, body from documents where project_id=$1 order by sort_order", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*Document, 0)
	for rows.Next() {
		doc := new(Document)
		doc.ProjectID = projectID
		err = rows.Scan(&doc.CreatedAt, &doc.UpdatedAt, &doc.ID, &doc.Title, &doc.SortOrder, &doc.Body)
		if err != nil {
			return nil, err
		}
		result = append(result, doc)
	}
	return result, nil
}

//LoadDocument --
func LoadDocument(id int) (*Document, error) {
	doc := new(Document)
	doc.ID = id
	err := DB.QueryRow("select created_at, updated_at, project_id, title, sort_order, body from documents where id=$1", id).Scan(&doc.CreatedAt, &doc.UpdatedAt, &doc.ProjectID, &doc.Title, &doc.SortOrder, &doc.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

//UpdateDocument --
func UpdateDocument(id int, title string, body string, sortOrder int) error {
	_, err := DB.Exec("update documents set title=$1, body=$2, sort_order=$3 where id=$4", title, body, sortOrder, id)
	return err
}

//DeleteDocument --
func DeleteDocument(id int) error {
	_, err := DB.Exec("delete from documents where id=$1", id)
	return err
}
