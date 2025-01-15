package service

import (
	"delivery/internal/models"
	repository "delivery/internal/repositories"
	"errors"
)

type OrderService interface {
	PlaceOrder(userID int, cartItems []models.Cart) (int, error)
}

type orderService struct {
	itemRepo  repository.ItemRepository
	orderRepo repository.OrderRepository
}

func NewOrderService(itemRepo repository.ItemRepository, orderRepo repository.OrderRepository) OrderService {
	return &orderService{
		itemRepo:  itemRepo,
		orderRepo: orderRepo,
	}
}

func (s *orderService) PlaceOrder(userID int, cartItems []models.Cart) (int, error) {
	var totalAmount float64
	for _, cartItem := range cartItems {
		item, err := s.itemRepo.GetItemByID(int(cartItem.ItemID))
		if err != nil {
			return 0, err
		}
		if cartItem.Quantity > item.Stock {
			return 0, errors.New("not enough stock for item: " + item.Name)
		}
		totalAmount += float64(cartItem.Quantity) * item.Price
	}

	order := models.Order{
		UserID:      uint(userID),
		TotalAmount: totalAmount,
	}

	orderID, err := s.orderRepo.SaveOrder(order, cartItems)
	if err != nil {
		return 0, err
	}

	for _, cartItem := range cartItems {
		if err := s.itemRepo.UpdateItemStock(int(cartItem.ItemID), cartItem.Quantity); err != nil {
			return 0, err
		}
	}

	return orderID, nil
}
