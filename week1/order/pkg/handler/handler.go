package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	orderv1 "github.com/mbakhodurov/homeworks2/week1/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week1/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/mbakhodurov/homeworks2/week1/shared/pkg/proto/payment/v1"
)

// OrderStatus — статус заказа
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// PaymentMethod — способ оплаты заказа
type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Order представляет заказ на постройку космического корабля
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// orderStore — хранилище заказов (in-memory)
type orderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов
func NewOrderStore() *orderStore {
	return &orderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// handler реализует интерфейс orderv1.Handler, сгенерированный ogen
type handler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *orderStore
}

// NewHandler создаёт новый обработчик заказов
func NewHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *orderStore,
) *handler {
	return &handler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика
func SetupServer(h *handler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации)
// GET /api/v1/orders/{order_uuid}.
func (h *handler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

func (h *handler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	// 1. Валидация: hull_uuid и engine_uuid обязательны
	if req.HullUUID == uuid.Nil || req.EngineUUID == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid и engine_uuid обязательны",
		}, nil
	}

	// Собираем все UUID запрошенных деталей
	uuids := []string{req.HullUUID.String(), req.EngineUUID.String()}
	if shieldUUID, ok := req.ShieldUUID.Get(); ok {
		if shieldUUID == uuid.Nil {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "shield_uuid имеет неверный формат",
			}, nil
		}
		uuids = append(uuids, shieldUUID.String())
	}
	if weaponUUID, ok := req.WeaponUUID.Get(); ok {
		if weaponUUID == uuid.Nil {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "weapon_uuid имеет неверный формат",
			}, nil
		}
		uuids = append(uuids, weaponUUID.String())
	}

	// 2. Получаем все детали одним запросом
	parts, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{Uuids: uuids},
	})
	if err != nil {
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "ошибка получения деталей из склада",
		}, nil
	}
	if len(parts.Part) != len(uuids) {
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "Component not found",
		}, nil
	}

	// 3–5. Проверяем наличие и остаток каждой детали, считаем цену
	var totalPrice int64
	for _, part := range parts.Part {
		if part.Info.StockQuantity <= 0 {
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: "Component out of stock",
			}, nil
		}
		totalPrice += part.Info.Price
	}

	// 6–7. Создаём заказ и сохраняем
	orderUUID := uuid.New()
	order := Order{
		OrderUUID:  orderUUID,
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		TotalPrice: totalPrice,
		Status:     OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}
	if shieldUUID, ok := req.ShieldUUID.Get(); ok {
		order.ShieldUUID = &shieldUUID
	}
	if weaponUUID, ok := req.WeaponUUID.Get(); ok {
		order.WeaponUUID = &weaponUUID
	}

	h.store.mu.Lock()
	h.store.orders[orderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay
func (h *handler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	if order.Status != OrderStatusPendingPayment {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен или отменён",
		}, nil
	}

	payResp, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		UserUuid:      params.OrderUUID.String(),
		PaymentMethod: toProtoPaymentMethod(req.PaymentMethod),
	})
	if err != nil {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "ошибка обработки платежа",
		}, nil
	}

	transactionUUID, err := uuid.Parse(payResp.TransactionUuid)
	if err != nil {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "неверный UUID транзакции от платёжного сервиса",
		}, nil
	}

	pm := PaymentMethod(req.PaymentMethod)
	order.Status = OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &pm

	h.store.mu.Lock()
	h.store.orders[params.OrderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.PayOrderResponse{TransactionUUID: transactionUUID}, nil
}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel
func (h *handler) CancelOrder(_ context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	if order.Status != OrderStatusPendingPayment {
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "нельзя отменить заказ в статусе " + string(order.Status),
		}, nil
	}

	order.Status = OrderStatusCancelled

	h.store.mu.Lock()
	h.store.orders[params.OrderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.CancelOrderResponse{}, nil
}

// GetAllOrders реализует операцию getAllOrders
// GET /api/v1/orders
func (h *handler) GetAllOrders(_ context.Context) (orderv1.GetAllOrdersRes, error) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	orders := make([]orderv1.OrderDto, 0, len(h.store.orders))
	for _, order := range h.store.orders {
		var shieldUUID orderv1.OptNilUUID
		if order.ShieldUUID != nil {
			shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
		}
		var weaponUUID orderv1.OptNilUUID
		if order.WeaponUUID != nil {
			weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
		}
		var transactionUUID orderv1.OptNilUUID
		if order.TransactionUUID != nil {
			transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
		}
		var paymentMethod orderv1.OptNilPaymentMethod
		if order.PaymentMethod != nil {
			paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
		}
		orders = append(orders, orderv1.OrderDto{
			OrderUUID:       order.OrderUUID,
			HullUUID:        order.HullUUID,
			EngineUUID:      order.EngineUUID,
			ShieldUUID:      shieldUUID,
			WeaponUUID:      weaponUUID,
			TotalPrice:      order.TotalPrice,
			TransactionUUID: transactionUUID,
			PaymentMethod:   paymentMethod,
			Status:          orderv1.OrderStatus(order.Status),
			CreatedAt:       order.CreatedAt,
		})
	}

	return &orderv1.GetAllOrdersResponse{
		Orders:     orders,
		TotalCount: int64(len(orders)),
	}, nil
}

func toProtoPaymentMethod(m orderv1.PaymentMethod) paymentv1.PaymentMethod {
	switch m {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
