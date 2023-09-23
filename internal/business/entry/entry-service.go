package entry

import (
	"context"
	"log"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
	"paywise/internal/models"
)

type entryService struct {
	entryRepo core.EntryRepo
}

type EntryServiceConfig struct {
	EntryRepo core.EntryRepo
}

func New(esc *EntryServiceConfig) core.EntryService {
	return &entryService{
		entryRepo: esc.EntryRepo,
	}
}

func (es *entryService) Create(ctx context.Context, reqDto *dtos.CreateEntryReq) (*models.Entry, error) {
	entry, err := es.entryRepo.Insert(ctx, &models.Entry{
		AccountID: reqDto.AccountID,
		Amount:    reqDto.Amount,
	})
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return entry, nil
}

func (es *entryService) GetAll(ctx context.Context, reqDto *dtos.GetAllEntriesOfAccountReq) ([]*models.Entry, error) {
	entries, err := es.entryRepo.Get(ctx, reqDto.AccountID)
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return entries, nil
}

func (es *entryService) GetbyID(ctx context.Context, reqDto *dtos.GetEntryByIdReq) (*models.Entry, error) {
	entry, err := es.entryRepo.GetbyID(ctx, reqDto.AccountID, reqDto.EntryID)
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return entry, nil
}

func (es *entryService) GetPage(ctx context.Context, reqDto *dtos.GetEntriesInPage) ([]*models.Entry, error) {
	entries, err := es.entryRepo.GetPage(ctx, reqDto.AccountID, reqDto.Limit, (reqDto.Offset-1)*reqDto.Limit)
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return entries, nil
}
