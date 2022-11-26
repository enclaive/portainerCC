import { useQuery } from "react-query";
import { getCoordinatorDeployments, getCoordinatorDeploymentForEnv, deployCoordinator } from "./coordinator.service";
import { CoordinatorDeployment } from "./types";


export function useCoordinatorDeployments<T = CoordinatorDeployment[]>(
  enabled = true,
  select: (data: CoordinatorDeployment[]) => T = (data) => data as unknown as T
) {
  const deployments = useQuery(
    ['coordinator-deployments'],
    () => getCoordinatorDeployments(),
    {
      meta: {
        error: { title: 'Failure', message: 'Unable to load coordinator deployments' },
      },
      enabled,
      select,
    }
  );

  return deployments;
}


export function useCoordinatorDeploymentForEnv<T = CoordinatorDeployment | undefined>(
  envId: number,
  enabled = true,
  select: (data: CoordinatorDeployment | undefined) => T = (data) => data as unknown as T
) {
  const deployment = useQuery(
    ['coordinator-deployment'],
    () => getCoordinatorDeploymentForEnv(envId),
    {
      meta: {
        error: { title: 'Failure', message: 'Unable to load coordinator deployment' },
      },
      enabled,
      select,
    }
  );

  return deployment;
}

export function useDeployCoordinator(
  envId: number,
  coordinatorId: number,
  enabled = false,
) {
  const deployment = useQuery(
    ['coordinator-deployment-deploy'],
    () => deployCoordinator(envId, coordinatorId),
    {
      meta: {
        error: { title: 'Failure', message: 'Unable to deploy coordinator' },
      },
      enabled,
    }
  );

  return deployment;
}