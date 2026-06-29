import { createClient } from "@connectrpc/connect";

import {
  IDENTITY_CONNECT_BASE_URL,
  STRATEGY_REGISTRY_CONNECT_BASE_URL,
} from "@/shared/config/env";
import { IdentityService } from "@/shared/api/gen/identity/proto/v1/service_pb";
import { StrategyRegistryService } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_service_pb";

import { createBrowserTransport } from "./transport";

export const identityTransport = createBrowserTransport(
  IDENTITY_CONNECT_BASE_URL,
);
export const strategyRegistryTransport = createBrowserTransport(
  STRATEGY_REGISTRY_CONNECT_BASE_URL,
);

export const identityClient = createClient(IdentityService, identityTransport);
export const strategyRegistryClient = createClient(
  StrategyRegistryService,
  strategyRegistryTransport,
);
