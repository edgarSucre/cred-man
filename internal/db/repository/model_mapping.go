package repository

import "github.com/edgarSucre/crm/pkg"

func (cl *CreateClientParams) FromDomain(param pkg.CreateClientParams) {
	cl.Birthdate = param.Birthdate
	cl.Country = param.Country
	cl.Email = param.Email
	cl.FullName = param.FullName
}

func (cl *Client) ToDomain() *pkg.Client {
	return &pkg.Client{
		Birthdate: *cl.Birthdate,
		Country:   *cl.Country,
		Email:     cl.Email,
		FullName:  cl.FullName,
		ID:        cl.ID,
		CreatedAt: cl.CreatedAt,
	}
}
