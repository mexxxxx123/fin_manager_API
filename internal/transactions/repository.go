package transactions

import "fin_manager_API/m/pkg/db"

type TransactionRepository struct {
	Database *db.Db
}

func NewTransactionRepository(database *db.Db) *TransactionRepository {
	return &TransactionRepository{
		Database: database,
	}
}

func (repo TransactionRepository) Create(trans *Transaction) (*Transaction, error) {
	result := repo.Database.DB.Create(trans)
	if result.Error != nil {
		return nil, result.Error
	}
	return trans, nil
}

func (repo TransactionRepository) Delete(trans *Transaction, id uint) error {
	result := repo.Database.
		Table("transactions").
		Where("deleted_at is null").
		Delete(&Transaction{}, id)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *TransactionRepository) TotalCount(userid uint) int64 {
	var count int64
	repo.Database.
		Table("transactions").
		Where("user_id = ?", userid).
		Where("deleted_at is null").
		Count(&count)

	return count

}

func (repo *TransactionRepository) GetAllTrans(limit int, offset int, userid uint) []Transaction {
	var transns []Transaction
	repo.Database.
		Table("transactions").
		Where("deleted_at is null").
		Where("user_id = ?", userid).
		Limit(limit).
		Offset(offset).
		Scan(&transns)

	return transns
}

func (repo *TransactionRepository) GetTransByTransID(id, userID uint) (*Transaction, error) {
	var trans Transaction
	result := repo.Database.Where("deleted_at is null").First(&trans, "ID = ? AND user_id = ?", id, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &trans, nil
}
