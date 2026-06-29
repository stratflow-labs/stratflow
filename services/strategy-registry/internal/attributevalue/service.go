package attributevalue

import tx "github.com/stratflow-labs/stratflow/internal/foundation/tx"

type Service struct {
	attributeValueRepo AttributeValueRepository
	txManager          tx.Manager
	clock              Clock
}

func NewService(
	attributeValueRepo AttributeValueRepository,
	txManager tx.Manager,
	clock Clock,
) *Service {
	return &Service{
		attributeValueRepo: attributeValueRepo,
		txManager:          txManager,
		clock:              clock,
	}
}
