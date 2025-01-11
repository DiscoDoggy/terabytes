package main

import (
	"net/http"
	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
)

type BlogPostContentPayload struct {
	ContentType 	string 	`json:"content_type"`
	ContentData 	string	`json:"content_data"`
	ContentOrder 	int		`json:"content_order"`
}

type CreateBlogPostPayload struct {
	Title 				string `json:"title"`
	Description 		string `json:"description"`
	Content 			[]BlogPostContentPayload `json:"content"` //TODO: CHANGE TO A MORE COMPLEX DATA STRUCT LATER... LATER TO ADAPT TO HANDLE ORDERING NESTED CONTENT AND DIFFERENT TYPES OF MEDIA
	Tags 				[]string `json:"tags"`
}

func (app *application) createBlogHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateBlogPostPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	blogPostContent := make([]store.BlogPostContent, len(payload.Content))
	for i:=0; i < len(blogPostContent); i++ {
		blogPostContent[i].ContentType = payload.Content[i].ContentType
		blogPostContent[i].ContentData = payload.Content[i].ContentData
		blogPostContent[i].ContentOrder = payload.Content[i].ContentOrder
	}

	blogPost := &store.BlogPost{
		UserId: "9efacfaf-2893-4665-b223-0ba333e04137",
		Title: payload.Title,
		Description: payload.Description,
		Content: blogPostContent,
		Tags: payload.Tags,
	}

	ctx := r.Context()
	
	err = app.store.Posts.Create(ctx, blogPost)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}