package service

import (
	"styleai-backend/internal/common"
	"styleai-backend/internal/models"
	"styleai-backend/internal/repository"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	variantRepo *repository.VariantRepository
}

func NewCartService(cartRepo *repository.CartRepository, variantRepo *repository.VariantRepository) *CartService {
	return &CartService{cartRepo: cartRepo, variantRepo: variantRepo}
}

func (s *CartService) AddItem(userID, variantID uint, quantity int) error {
	if quantity <= 0 {
		return common.ErrInvalidQuantity
	}

	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	variant, err := s.variantRepo.FindByID(variantID)
	if err != nil {
		return err
	}

	if quantity > variant.Stock {
		return common.ErrInsufficientStock
	}

	item, err := s.cartRepo.FindItem(cart.ID, variantID)
	if err != nil {
		return err
	}

	if item != nil {
		newQty := item.Quantity + quantity

		if newQty > variant.Stock {
			return common.ErrInsufficientStock
		}

		item.Quantity = newQty
		return s.cartRepo.UpdateItem(item)
	}

	newItem := models.CartItem{
		CartID:    cart.ID,
		VariantID: variantID,
		Quantity:  quantity,
		Price:     variant.Price,
	}

	return s.cartRepo.CreateItem(&newItem)
}

func (s *CartService) GetCart(userID uint) (*models.Cart, float64, error) {
	cart, err := s.cartRepo.GetCartWithItems(userID)
	if err != nil {
		return nil, 0, err
	}

	var total float64
	for _, item := range cart.Items {
		total += item.Price * float64(item.Quantity)
	}

	return cart, total, nil
}

func (s *CartService) UpdateItem(userID, itemID uint, quantity int) error {
	if quantity < 0 {
		return common.ErrInvalidQuantity
	}

	item, err := s.cartRepo.GetItemByID(itemID)
	if err != nil {
		return common.ErrCartItemNotFound
	}

	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	if item.CartID != cart.ID {
		return common.ErrUnauthorized
	}

	if quantity == 0 {
		return s.cartRepo.DeleteItem(itemID)
	}

	variant, err := s.variantRepo.FindByID(item.VariantID)
	if err != nil {
		return err
	}

	if quantity > variant.Stock {
		return common.ErrInsufficientStock
	}

	item.Quantity = quantity
	return s.cartRepo.UpdateItem(item)
}

func (s *CartService) RemoveItem(userID, itemID uint) error {
	item, err := s.cartRepo.GetItemByID(itemID)
	if err != nil {
		return common.ErrCartItemNotFound
	}

	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	if item.CartID != cart.ID {
		return common.ErrUnauthorized
	}

	return s.cartRepo.DeleteItem(itemID)
}
