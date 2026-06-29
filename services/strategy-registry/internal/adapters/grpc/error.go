package strategygrpc

import (
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	"github.com/stratflow-labs/stratflow/internal/grpcserver"

	"google.golang.org/grpc/codes"
)

func mapError(err error) error {
	if appCode := apperr.Code(err); appCode != "" {
		return mapAppError(err)
	}

	return grpcserver.NewStatusError(
		codes.Internal,
		"strategyRegistry.internal",
		http.StatusText(http.StatusInternalServerError),
	)
}

func mapAppError(err error) error {
	code := normalizeReason(apperr.Code(err))
	message := normalizeMessage(code, apperr.Message(err))

	switch apperr.KindOf(err) {
	case apperr.KindNotFound:
		return grpcserver.NewStatusError(codes.NotFound, code, message)
	case apperr.KindAlreadyExists:
		return grpcserver.NewStatusError(codes.AlreadyExists, code, message)
	case apperr.KindInvalidArgument:
		return grpcserver.NewStatusError(codes.InvalidArgument, code, message)
	case apperr.KindPermissionDenied:
		return grpcserver.NewStatusError(codes.PermissionDenied, code, message)
	case apperr.KindUnauthenticated:
		return grpcserver.NewStatusError(codes.Unauthenticated, code, message)
	default:
		return grpcserver.NewStatusError(codes.Internal, "strategyRegistry.internal", http.StatusText(http.StatusInternalServerError))
	}
}

func normalizeReason(code string) string {
	switch code {
	case "strategy.pageTooLarge":
		return "strategy.pageOutOfRange"
	case "attribute.pageTooLarge":
		return "attribute.pageOutOfRange"
	case "attributeValue.pageTooLarge":
		return "attributeValue.pageOutOfRange"
	case "strategy.clone.itemsEmpty":
		return "strategy.cloneItemsEmpty"
	case "strategy.clone.duplicateSlug":
		return "strategy.cloneDuplicateSlug"
	case "attributeValue.relation.combinationNotFound":
		return "attributeValue.relationCombinationNotFound"
	case "attributeValue.relation.duplicate":
		return "attributeValue.relationDuplicate"
	case "attributeValue.relation.selfReference":
		return "attributeValue.relationSelfReference"
	case "strategy.batchAction.empty":
		return "strategyGraph.batchActionEmpty"
	default:
		return code
	}
}

func normalizeMessage(code, message string) string {
	switch code {
	case "strategy.pageOutOfRange", "attribute.pageOutOfRange", "attributeValue.pageOutOfRange":
		return "page is out of range"
	default:
		return message
	}
}

func invalidArgument(reason, message string) error {
	return grpcserver.NewStatusError(codes.InvalidArgument, reason, message)
}
