package handler

import (
	"e-wallet/internal/dto/request"
	"e-wallet/internal/infrastructure/logger"
	"e-wallet/internal/usecase"
	apperrors "e-wallet/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	checkUseCase        *usecase.WalletCheckUseCase
	depositUseCase      *usecase.WalletDepositUseCase
	balanceUseCase      *usecase.WalletBalanceUseCase
	monthlyStatsUseCase *usecase.WalletMonthlyStatsUseCase
}

func NewWalletHandler(
	checkUseCase *usecase.WalletCheckUseCase,
	depositUseCase *usecase.WalletDepositUseCase,
	balanceUseCase *usecase.WalletBalanceUseCase,
	monthlyStatsUseCase *usecase.WalletMonthlyStatsUseCase,
) *WalletHandler {
	return &WalletHandler{
		checkUseCase:        checkUseCase,
		depositUseCase:      depositUseCase,
		balanceUseCase:      balanceUseCase,
		monthlyStatsUseCase: monthlyStatsUseCase,
	}
}

// CheckWallet godoc
// @Summary Check wallet existence
// @Description Checks if a wallet account exists by account_id
// @Tags Wallet
// @Accept json
// @Produce json
// @Param X-UserId header string true "User ID"
// @Param X-Digest header string true "HMAC-SHA1 digest"
// @Param request body request.CheckWalletRequest true "Check wallet request"
// @Success 200 {object} response.CheckWalletResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Security HMACAuth
// @Security HMACDigest
// @Router /wallet/check [post]
func (h *WalletHandler) CheckWallet(c *gin.Context) {
	ip := c.ClientIP()
	logger.Info.Printf("[handler.CheckWallet]: Client with IP %s requested wallet check (request ID: %s)", ip, c.GetString("request_id"))

	var req request.CheckWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, apperrors.ErrInvalidRequest)
		logger.Error.Printf("[handler.CheckWallet]: Failed to bind request: %v", err)
		return
	}

	resp, err := h.checkUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err)
		return
	}
	logger.Info.Printf("[CheckWallet]: Client with IP %s successfully checked wallet (request_id=%s)", ip, c.GetString("request_id"))

	c.JSON(http.StatusOK, resp)
}

// Deposit godoc
// @Summary Deposit to wallet
// @Description Deposits money to a wallet account. Amount is in dirams (1 TJS = 100 dirams). Validates limits: 10,000 TJS for unidentified, 100,000 TJS for identified wallets.
// @Tags Wallet
// @Accept json
// @Produce json
// @Param X-UserId header string true "User ID"
// @Param X-Digest header string true "HMAC-SHA1 digest"
// @Param request body request.DepositRequest true "Deposit request"
// @Success 200 {object} response.DepositResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Security HMACAuth
// @Security HMACDigest
// @Router /wallet/deposit [post]
func (h *WalletHandler) Deposit(c *gin.Context) {
	ip := c.ClientIP()
	logger.Info.Printf("[handler.Deposit]: Client with IP %s requested wallet deposit (request ID: %s)", ip, c.GetString("request_id"))

	var req request.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, apperrors.ErrInvalidRequest)
		logger.Error.Printf("[handler.Deposit]: Failed to bind request: %v", err)
		return
	}

	resp, err := h.depositUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err)
		return
	}
	logger.Info.Printf("[Deposit]: Client with IP %s successfully deposited money (request_id=%s)", ip, c.GetString("request_id"))

	c.JSON(http.StatusOK, resp)
}

// GetBalance godoc
// @Summary Get wallet balance
// @Description Returns current wallet balance in dirams (1 TJS = 100 dirams)
// @Tags Wallet
// @Accept json
// @Produce json
// @Param X-UserId header string true "User ID"
// @Param X-Digest header string true "HMAC-SHA1 digest"
// @Param request body request.GetBalanceRequest true "Get balance request"
// @Success 200 {object} response.GetBalanceResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Security HMACAuth
// @Security HMACDigest
// @Router /wallet/balance [post]
func (h *WalletHandler) GetBalance(c *gin.Context) {
	ip := c.ClientIP()
	logger.Info.Printf("[handler.GetBalance]: Client with IP %s requested wallet balance (request ID: %s)", ip, c.GetString("request_id"))

	var req request.GetBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, apperrors.ErrInvalidRequest)
		logger.Error.Printf("[handler.GetBalance]: Failed to bind request: %v", err)
		return
	}

	resp, err := h.balanceUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err)
		return
	}
	logger.Info.Printf("[GetBalance]: Client with IP %s successfully retrieved wallet balance (request_id=%s)", ip, c.GetString("request_id"))

	c.JSON(http.StatusOK, resp)
}

// GetMonthlyStats godoc
// @Summary Get monthly statistics
// @Description Returns total count and amount of deposits for the current month
// @Tags Wallet
// @Accept json
// @Produce json
// @Param X-UserId header string true "User ID"
// @Param X-Digest header string true "HMAC-SHA1 digest"
// @Param request body request.GetMonthlyStatsRequest true "Get monthly stats request"
// @Success 200 {object} response.MonthlyStatsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Security HMACAuth
// @Security HMACDigest
// @Router /wallet/monthly-stats [post]
func (h *WalletHandler) GetMonthlyStats(c *gin.Context) {
	ip := c.ClientIP()
	logger.Info.Printf("[handler.GetMonthlyStats]: Client with IP %s requested monthly statistics (request ID: %s)", ip, c.GetString("request_id"))

	var req request.GetMonthlyStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, apperrors.ErrInvalidRequest)
		logger.Error.Printf("[handler.GetMonthlyStats]: Failed to bind request: %v", err)
		return
	}

	resp, err := h.monthlyStatsUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err)
		return
	}
	logger.Info.Printf("[GetMonthlyStats]: Client with IP %s successfully retrieved monthly statistics (request_id=%s)", ip, c.GetString("request_id"))

	c.JSON(http.StatusOK, resp)
}
