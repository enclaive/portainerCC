import { Environment } from '../environments/types';

export interface EndpointProvider {
  setEndpointID(id: Environment['Id']): void;
  setEndpointPublicURL(url?: string): void;
  setCurrentEndpoint(endpoint: Environment | undefined): void;
}

export interface StateManager {
  updateEndpointState(endpoint: Environment): Promise<void>;
}
