package validates

type (
	OrderValidateImpl struct{}
)

type OrderValidate interface {
}

func NewOrderValidate() OrderValidate {
	return &OrderValidateImpl{}
}
