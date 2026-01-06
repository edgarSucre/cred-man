package repository

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/bank"
	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/internal/domain/credit"
	"github.com/edgarSucre/mye"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type creditRepository struct {
	q Querier
}

func NewCreditRepository(pool *pgxpool.Pool) *creditRepository {
	return &creditRepository{q: New(pool)}
}

func (repo *creditRepository) CreateCredit(
	ctx context.Context,
	in credit.Credit,
) (credit.Credit, error) {
	q := repo.querier(ctx)

	model, err := q.CreateCredit(ctx, CreateCreditParams{
		BankID:     in.BankID().UUID(),
		ClientID:   in.ClientID().UUID(),
		MinPayment: in.MinPayment(),
		MaxPayment: in.MaxPayment(),
		CreditType: CreditType(in.CreditType().String()),
		TermMonths: int16(in.TermMonths()),
		Status:     CreditStatus(in.Status().String()),
	})

	if err != nil {
		code, slug := CodeAndSlug(err)

		if code == mye.CodeInvalid {
			err = mye.Wrap(err, code, slug, "bank_insert_failed").
				WithUserMsg("can't create credit, make sure credit values are valid")

			return credit.Credit{}, err
		}

		if code == mye.CodeTimeout {
			err = mye.Wrap(err, code, slug, "bank_insert_failed").
				WithUserMsg("Credit creation is taking a bit too log due to high traffic. Please try again in a few seconds.")
		}

		return credit.Credit{}, mye.Wrap(err, code, slug, "failed to create credit")
	}

	newCredit, err := model.ToDomain()
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditRepository.CreateCredit > %w", err)
	}

	return newCredit, err
}

func (m Credit) ToDomain() (credit.Credit, error) {
	id, err := credit.NewIDFromUUID(m.ID)
	if err != nil {
		err = mye.Wrap(
			err,
			mye.CodeInternal,
			ErrDataIntegrity,
			"corrupted ID in the database",
		).WithAttribute("ID", m.ID).WithAttribute("table", "credits")

		return credit.Credit{}, err
	}

	bankID, err := bank.NewID(m.BankID.String())
	if err != nil {
		err = mye.Wrap(
			err,
			mye.CodeInternal,
			ErrDataIntegrity,
			"corrupted bank ID in the database",
		).WithAttribute("bank_id", m.BankID).WithAttribute("table", "credits")

		return credit.Credit{}, err
	}

	clientID, err := client.NewID(m.ClientID.String())
	if err != nil {
		err = mye.Wrap(
			err,
			mye.CodeInternal,
			ErrDataIntegrity,
			"corrupted client ID in the database",
		).WithAttribute("client_id", m.ClientID).WithAttribute("table", "credits")

		return credit.Credit{}, err
	}

	creditType, err := credit.CreditTypeFromString(string(m.CreditType))
	if err != nil {
		err = mye.Wrap(
			err,
			mye.CodeInternal,
			ErrDataIntegrity,
			"corrupted credit_type in the database",
		).WithAttribute("credit_type", m.CreditType).WithAttribute("table", "credits")

		return credit.Credit{}, err
	}

	creditStatus, err := credit.CreditStatusFromString(string(m.Status))
	if err != nil {
		err = mye.Wrap(
			err,
			mye.CodeInternal,
			ErrDataIntegrity,
			"corrupted credit_status in the database",
		).WithAttribute("credit_status", m.Status).WithAttribute("table", "credits")

		return credit.Credit{}, err
	}

	return credit.Rehydrate(credit.RehydrateOpts{
		BankID:     bankID,
		ClientID:   clientID,
		CreatedAt:  m.CreatedAt,
		CreditType: creditType,
		ID:         id,
		MaxPayment: m.MaxPayment,
		MinPayment: m.MinPayment,
		Status:     creditStatus,
		TermMonths: int(m.TermMonths),
	}), nil
}

func (repo *creditRepository) GetCredit(ctx context.Context, id credit.ID) (credit.Credit, error) {
	c, err := repo.q.GetCredit(ctx, id.UUID())
	if err != nil {
		code, slug := CodeAndSlug(err)

		if code == mye.CodeNotFound {
			return credit.Credit{}, mye.Wrap(err, code, slug, "credit not found").
				WithAttribute("id", id.String()).
				WithUserMsg("we couldn't find that credit")
		}

		if code == mye.CodeTimeout {
			return credit.Credit{}, mye.Wrap(err, code, slug, "getCredit time out").
				WithUserMsg("The credit search is taking a bit too long due to high traffic. Please try again in a few seconds")
		}

		return credit.Credit{}, mye.Wrap(err, code, slug, "failed to retrieve credit")
	}

	domainCredit, err := c.ToDomain()
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditRepository.GetCredit > %w", err)
	}

	return domainCredit, nil
}

func (repo creditRepository) GetAggregate(
	ctx context.Context,
	creditID credit.ID,
	clientID client.ID,
) (*credit.CreditAggregate, error) {
	credits, err := repo.q.GetClientCredits(ctx, clientID.UUID())
	if err != nil {
		code, slug := CodeAndSlug(err)

		if code == mye.CodeTimeout {
			return nil, mye.Wrap(err, code, slug, "getClientCredits time out").
				WithUserMsg("The credits search is taking a bit too long due to high traffic. Please try again in a few seconds")
		}

		return nil, mye.Wrap(err, code, slug, "getClientCredits failure")
	}

	dCredits, err := creditsToDomain(credits)
	if err != nil {
		return nil, fmt.Errorf("creditsToDomain > %w", err)
	}

	var creditToProcess credit.Credit

	for _, c := range dCredits {
		if c.ID().IsEqual(creditID) {
			creditToProcess = c
		}
	}

	if creditToProcess.ID().IsEmpty() {
		return nil, mye.New(
			mye.CodeNotFound,
			"credit_to_process_not_found",
			"can't find the credit to process in the list of client's credits",
		)
	}

	return credit.RehydrateAggregate(
		dCredits,
		creditToProcess.CreditType(),
		creditID,
		creditToProcess.Status(),
	), nil
}

func creditsToDomain(credits []Credit) ([]credit.Credit, error) {
	dCredits := make([]credit.Credit, len(credits))

	for i, c := range credits {
		dCredit, err := c.ToDomain()
		if err != nil {
			return nil, err
		}

		dCredits[i] = dCredit
	}

	return dCredits, nil
}

func (repo creditRepository) ProcessCredit(
	ctx context.Context,
	creditAggregate credit.CreditAggregate,
) error {
	q := repo.querier(ctx)

	status := creditAggregate.Status().String()

	params := UpdateCreditStatusParams{
		Status: CreditStatus(status),
		ID:     creditAggregate.ID().UUID(),
	}

	if err := q.UpdateCreditStatus(ctx, params); err != nil {
		return fmt.Errorf("repo.UpdateCreditStatus > %w", err)
	}

	return nil
}

func (repo creditRepository) querier(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return New(tx)
	}

	return repo.q
}
