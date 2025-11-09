package mapper

import (
	"e-wallet/internal/domain/entity"
	"e-wallet/internal/domain/valueobject"
	"e-wallet/internal/infrastructure/database/models"
)

type WalletMapper struct{}

func NewWalletMapper() *WalletMapper {
	return &WalletMapper{}
}

func (m *WalletMapper) ToDomain(dbWallet *models.Wallet) (*entity.Wallet, error) {
	accountID, err := valueobject.NewAccountID(dbWallet.AccountID)
	if err != nil {
		return nil, err
	}

	walletType, err := valueobject.NewWalletType(dbWallet.Type)
	if err != nil {
		return nil, err
	}

	balance, err := valueobject.NewMoneyFromMinor(dbWallet.Balance)
	if err != nil {
		return nil, err
	}

	return &entity.Wallet{
		ID:        dbWallet.ID,
		AccountID: accountID,
		Type:      walletType,
		Balance:   balance,
		CreatedAt: dbWallet.CreatedAt,
		UpdatedAt: dbWallet.UpdatedAt,
	}, nil
}

func (m *WalletMapper) ToModel(wallet *entity.Wallet) *models.Wallet {
	return &models.Wallet{
		ID:        wallet.ID,
		AccountID: wallet.AccountID.Value(),
		Type:      wallet.Type.String(),
		Balance:   wallet.Balance.Amount(),
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
