// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlapi

import (
	domains "github.com/mercari/mtc2018-web/server/domains"
)

// セッションに対していいね！する時の引数です。
// uuidはいいね！した人のUUIDです。
type CreateLikeInput struct {
	ClientMutationID *string `json:"clientMutationId"`
	UUID             string  `json:"uuid"`
	SessionID        int     `json:"sessionId"`
}

type CreateLikePayload struct {
	ClientMutationID *string      `json:"clientMutationId"`
	Like             domains.Like `json:"like"`
}

type LikeEvent struct {
	Session domains.Session `json:"session"`
	Likes   int             `json:"likes"`
}

type NewsConnection struct {
	PageInfo PageInfo       `json:"pageInfo"`
	Edges    []NewsEdge     `json:"edges"`
	Nodes    []domains.News `json:"nodes"`
}

type NewsEdge struct {
	Cursor *string      `json:"cursor"`
	Node   domains.News `json:"node"`
}

type Node interface{}

type PageInfo struct {
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
}

type SessionConnection struct {
	PageInfo PageInfo          `json:"pageInfo"`
	Edges    []SessionEdge     `json:"edges"`
	Nodes    []domains.Session `json:"nodes"`
}

type SessionEdge struct {
	Cursor *string         `json:"cursor"`
	Node   domains.Session `json:"node"`
}

type SessionListInput struct {
	Tmp *string `json:"tmp"`
}
