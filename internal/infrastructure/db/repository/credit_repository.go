package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/bank"
	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/internal/domain/credit"
	"github.com/edgarSucre/crm/pkg/terror"
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
		return credit.Credit{}, fmt.Errorf("creditRepository.CreateCredit > %w", err)
	}

	newCredit, err := model.ToDomain()
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditRepository.CreateCredit > %w", err)
	}

	return newCredit, nil
}

func (m Credit) ToDomain() (credit.Credit, error) {
	id, err := credit.NewIDFromUUID(m.ID)
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditModel.ToDomain > %w", terror.ToInternal(err))
	}

	bankID, err := bank.NewID(m.BankID.String())
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditModel.ToDomain > %w", terror.ToInternal(err))
	}

	clientID, err := client.NewID(m.ClientID.String())
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditModel.ToDomain > %w", terror.ToInternal(err))
	}

	creditType, err := credit.CreditTypeFromString(string(m.CreditType))
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditModel.ToDomain > %w", terror.ToInternal(err))
	}

	creditStatus, err := credit.CreditStatusFromString(string(m.Status))
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditModel.ToDomain > %w", terror.ToInternal(err))
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

var ErrNoCreditFound = terror.NotFound.New("not-found", "no credit found")

func (repo *creditRepository) GetCredit(ctx context.Context, id credit.ID) (credit.Credit, error) {
	c, err := repo.q.GetCredit(ctx, id.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			err = ErrNoCreditFound
		}

		return credit.Credit{}, fmt.Errorf("creditRepository.GetCredit > %w", err)
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
		return nil, fmt.Errorf("repo.GetClientCredits > %w", err)
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
		return nil, ErrNoCreditFound
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
