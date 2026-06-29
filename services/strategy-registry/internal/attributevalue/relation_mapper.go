package attributevalue

import attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

func relationInputsFromDomain(relations []attributeValuedomain.AttributeValueRelation) []AttributeValueRelationInput {
	if len(relations) == 0 {
		return []AttributeValueRelationInput{}
	}

	items := make([]AttributeValueRelationInput, len(relations))
	for i := range relations {
		items[i] = AttributeValueRelationInput{
			ToAttributeID: relations[i].ToAttributeID,
			ToValueID:     relations[i].ToValueID,
		}
	}

	return items
}
