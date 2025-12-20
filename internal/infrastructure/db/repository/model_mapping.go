package repository

import "github.com/edgarSucre/crm/pkg/domain"

func (cl *CreateClientParams) FromDomain(param domain.CreateClientParams) {
	cl.Birthdate = param.Birthdate
	cl.Country = param.Country
	cl.Email = param.Email
	cl.FullName = param.FullName
}

func (cl *Client) ToDomain() *domain.Client {
	return &domain.Client{
		Birthdate: *cl.Birthdate,
		Country:   *cl.Country,
		Email:     cl.Email,
		FullName:  cl.FullName,
		ID:        cl.ID,
		CreatedAt: cl.CreatedAt,
	}
}

func (bp *CreateBankParams) FromDomain(param domain.CreateBankParams) {
	bp.Name = param.Name
	bp.Type = BankType(param.Type)
}

func (bp Bank) ToDomain() *domain.Bank {
	return &domain.Bank{
		ID:   bp.ID,
		Name: bp.Name,
		Type: domain.BankType(bp.Type),
	}
}

func (cp *CreateCreditParams) FromDomain(param domain.CreateCreditParams) {
	cp.BankID = param.BankID
	cp.ClientID = param.ClientID
	cp.CreditType = CreditType(param.CreditType)
}

func (cp *Credit) ToDomain() *domain.Credit {
	return &domain.Credit{
		BankID:     cp.BankID,
		ClientID:   cp.ClientID,
		CreatedAt:  cp.CreatedAt,
		CreditType: domain.CreditType(cp.CreditType),
		ID:         cp.ID,
		MaxPayment: cp.MaxPayment,
		MinPayment: cp.MinPayment,
		Status:     domain.CreditStatus(cp.Status),
		TermMonths: int(cp.TermMonths),
	}
}
