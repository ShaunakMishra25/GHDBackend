package types

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

func (r UserRole) IsValid() bool {
	return r == RoleCustomer || r == RoleAdmin
}

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderConfirmed  OrderStatus = "confirmed"
	OrderAccepted   OrderStatus = "accepted"
	OrderPreparing  OrderStatus = "preparing"
	OrderDispatched OrderStatus = "dispatched"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderPending, OrderConfirmed, OrderAccepted, OrderPreparing, OrderDispatched, OrderDelivered, OrderCancelled:
		return true
	}
	return false
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	transitions := map[OrderStatus][]OrderStatus{
		OrderPending:    {OrderConfirmed, OrderCancelled},
		OrderConfirmed:  {OrderAccepted, OrderCancelled},
		OrderAccepted:   {OrderPreparing, OrderCancelled},
		OrderPreparing:  {OrderDispatched},
		OrderDispatched: {OrderDelivered},
		OrderDelivered:  {},
		OrderCancelled:  {},
	}

	allowed, ok := transitions[s]
	if !ok {
		return false
	}
	for _, status := range allowed {
		if status == next {
			return true
		}
	}
	return false
}

type PaymentStatus string

const (
	PaymentCreated   PaymentStatus = "created"
	PaymentAttempted PaymentStatus = "attempted"
	PaymentCaptured  PaymentStatus = "captured"
	PaymentFailed    PaymentStatus = "failed"
)

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentCreated, PaymentAttempted, PaymentCaptured, PaymentFailed:
		return true
	}
	return false
}

type DeliveryChargeConfig struct {
	FlatAmount      float64 `json:"flat_amount"`
	MaxOrderAmount  float64 `json:"max_order_amount"`
	PerItemAmount   float64 `json:"per_item_amount"`
}
