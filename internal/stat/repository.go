package stat

import (
	"fin_manager_API/m/pkg/db"
	"fmt"
	"time"
)

type StatRepository struct {
	db *db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		db: db,
	}
}

type UserStats struct {
	UserID           uint    `json:"user_id"`
	Balance          float64 `json:"balance"`
	TransactionCount int64   `json:"transaction_count"`
	IncomeTotal      float64 `json:"income_total"`
	ExpenseTotal     float64 `json:"expense_total"`
}

func (repo *StatRepository) GetUserStat(userID uint) (*UserStats, error) {
	var stats UserStats

	err := repo.db.DB.Raw(`
        SELECT 
            user_id,
            COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income_total,
            COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense_total,
 			COALESCE(SUM(
        			CASE 
            			WHEN type = 'income' THEN amount
            			WHEN type = 'expense' THEN -amount  
            			ELSE 0 
        			END
    			), 0) as balance,
            COUNT(*) as transaction_count
        FROM transactions 
        WHERE user_id = ? AND deleted_at IS NULL
        GROUP BY user_id
    `, userID).Scan(&stats).Error

	return &stats, err
}

func (repo *StatRepository) GetUserStatByPeriod(userID uint, from, to time.Time, by string) ([]UserStats, error) {
	var stats []UserStats

	format := "'YYYY-MM-DD'"
	if by == SortByMonth {
		format = "'YYYY-MM'"
	}

	query := fmt.Sprintf(`
        SELECT 
            user_id,
            to_char(created_at, %s) as period,
            COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income_total,
            COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense_total,
            COALESCE(SUM(
                CASE 
                    WHEN type = 'income' THEN amount
                    WHEN type = 'expense' THEN -amount  
                    ELSE 0 
                END
            ), 0) as balance,
            COUNT(*) as transaction_count
        FROM transactions 
        WHERE user_id = ? 
            AND created_at >= ? 
            AND created_at <= ?
            AND deleted_at IS NULL
        GROUP BY user_id, to_char(created_at, %s)  -- ← ДОБАВИЛ!
        ORDER BY period
    `, format, format) // format используется два раза

	err := repo.db.Raw(query, userID, from, to).Scan(&stats).Error
	return stats, err
}

// func (repo *StatRepository) GetUserStatByMonth(userID uint, from, to time.Time, by string) (*UserStats, error) {
// 	var stats UserStats
// 	var selectQuery string
// 	switch by {
// 	case SortByDay:
// 		selectQuery = "to_char(created_at, 'YYYY-MM-DD') as period "
// 	case SortByMonth:
// 		selectQuery = "to_char(created_at, 'YYYY-MM') as period "
// 	}
//
// 	err := repo.db.DB.Select(selectQuery).
// 		Raw(`
//         SELECT
//             user_id,
//             COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income_total,
//             COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense_total,
//  			COALESCE(SUM(
//         			CASE
//             			WHEN type = 'income' THEN amount
//             			WHEN type = 'expense' THEN -amount
//             			ELSE 0
//         			END
//     			), 0) as balance,
//             COUNT(*) as transaction_count
//         FROM transactions
//         WHERE user_id = ? AND deleted_at IS NULL
//     `, userID, from, to).
// 		Group("period").
// 		Order("period").
// 		Scan(&stats).Error
//
// 	return &stats, err
// }
