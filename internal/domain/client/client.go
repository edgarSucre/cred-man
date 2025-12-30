package client

import (
	"time"
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
	var cl Client

	if !b.IsValid() {
		return cl, ErrInvalidBirthdate
	}

	if email.IsValid() {
		return cl, ErrInvalidEmail
	}

	if len(fullName) == 0 {
		return cl, ErrInvalidClientFullName
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
