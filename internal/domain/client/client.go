package client

import (
	"time"

	"github.com/edgarSucre/mye"
)

type Client struct {
	id        ID
	birthdate Birthdate
	country   *string
	email     Email
	fullName  string
	createdAt time.Time
}

func New(b Birthdate, country *string, email Email, fullName string) (Client, error) {
	err := mye.New(mye.CodeInvalid, "client_creation_failed", "validation failed").
		WithUserMsg("client creation validation failed")

	if !b.IsValid() {
		err.WithField("birthdate", "birthdate must be a valid date")
	}

	if email.IsValid() {
		err.WithField("email", "client email should be a valid email address")
	}

	if len(fullName) == 0 {
		err.WithField("email", "client full name can't be empty")
	}

	if err.HasFields() {
		return Client{}, err
	}

	return Client{
		birthdate: b,
		country:   country,
		email:     email,
		fullName:  fullName,
	}, nil
}

func (client *Client) ID() ID {
	return client.id
}

func (client *Client) Birthdate() Birthdate {
	return client.birthdate
}

func (client *Client) Country() *string {
	return client.country
}

func (client *Client) Email() Email {
	return client.email
}

func (client *Client) FullName() string {
	return client.fullName
}

func (client *Client) CreatedAt() time.Time {
	return client.createdAt
}

func Rehydrate(
	birthdate Birthdate,
	country *string,
	createdAt time.Time,
	email Email,
	fullName string,
	id ID,
) Client {
	return Client{
		birthdate: birthdate,
		country:   country,
		createdAt: createdAt,
		email:     email,
		fullName:  fullName,
		id:        id,
	}
}
