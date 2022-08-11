import { Environment } from '../environments/types';

export interface StateManager {
  updateEndpointState(endpoint: Environment): Promise<void>;
}
