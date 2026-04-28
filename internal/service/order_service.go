package service

import (
	"styleai-backend/internal/common"
	"styleai-backend/internal/models"
	"styleai-backend/internal/repository"

	"gorm.io/gorm"
)

type OrderService struct {
	cartRepo    *repository.CartRepository
	orderRepo   *repository.OrderRepository
	variantRepo *repository.VariantRepository
	db          *gorm.DB
}

func NewOrderService(
	cartRepo *repository.CartRepository,
	orderRepo *repository.OrderRepository,
	variantRepo *repository.VariantRepository,
	db *gorm.DB,
) *OrderService {
	return &OrderService{cartRepo: cartRepo, orderRepo: orderRepo, variantRepo: variantRepo, db: db}
}

func (s *OrderService) Checkout(userID uint) (*models.Order, error) {
	var order *models.Order

	err := s.db.Transaction(func(tx *gorm.DB) error {

		cartRepo := s.cartRepo.WithTx(tx)
		orderRepo := s.orderRepo.WithTx(tx)
		variantRepo := s.variantRepo.WithTx(tx)

		cart, err := cartRepo.GetCartWithItems(userID)
		if err != nil {
			return err
		}

		if len(cart.Items) == 0 {
			return common.ErrCartEmpty
		}

		var total float64

		order = &models.Order{
			UserID: userID,
			Status: models.OrderStatusPending,
		}

		if err := orderRepo.Create(order); err != nil {
			return err
		}

		for _, item := range cart.Items {
			if err := variantRepo.ReduceStock(item.VariantID, item.Quantity); err != nil {
				return err
			}

			orderItem := models.OrderItem{
				OrderID:   order.ID,
				VariantID: item.VariantID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}

			if err := orderRepo.CreateItem(&orderItem); err != nil {
				return err
			}

			total += item.Price * float64(item.Quantity)
		}

		order.TotalPrice = total
		if err := orderRepo.Update(order); err != nil {
			return err
		}

		if err := cartRepo.ClearCart(cart.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrdersByUser(userID uint) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *OrderService) GetOrderByID(userID, orderID uint) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	if order.UserID != userID {
		return nil, common.ErrUnauthorized
	}

	return order, nil
}
