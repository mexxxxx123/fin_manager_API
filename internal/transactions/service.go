package transactions

import (
	"context"
	"fin_manager_API/m/middleware"
)

type TransactionService struct {
	*TransactionRepository
}

func NewTrasactionService(repo *TransactionRepository) *TransactionService {
	return &TransactionService{
		TransactionRepository: repo,
	}
}

func (service *TransactionService) GetAllTransaction(ctx context.Context, limit int, offset int) (*[]Transaction, int64, string) {
	userId, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, 0, "Cant get userId"
	}
	transactions := service.TransactionRepository.GetAllTrans(limit, offset, userId)
	count := service.TransactionRepository.TotalCount(userId)
	return &transactions, count, ""
}

func (service *TransactionService) GetTransaction(ctx context.Context, transID uint) (*Transaction, string) {
	userId, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, "Cant get userId"
	}
	trans, err := service.TransactionRepository.GetTransByTransID(transID, userId)
	if err != nil {
		return nil, "Transaction not found or acces denied"
	}
	return trans, ""

}

func (service *TransactionService) DeleteTransByID(ctx context.Context, transID64 uint64) string {
	userId, err := middleware.GetUserID(ctx)
	if err != nil {
		return "Cant get userId"
	}
	trans, err := service.TransactionRepository.GetTransByTransID(uint(transID64), userId)
	if err != nil {
		return "Transaction not found or acces denied"
	}
	err = service.TransactionRepository.Delete(trans, uint(transID64))
	if err != nil {
		return "Delete failed"
	}
	return ""

}
