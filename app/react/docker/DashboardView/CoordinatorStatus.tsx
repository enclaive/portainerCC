import clsx from 'clsx';

import { Icon } from '@/react/components/Icon';

interface Props {
  verified: boolean;
}


export function useCoordinatorStatusComponent(verified: boolean) {
  return <CoordinatorStatus verified={verified} />;
}

export function CoordinatorStatus({ verified }: Props) {
  return (
    <div className="pull-right">
      <div>
        {verified &&
          <div className="vertical-center space-right">
            <Icon
              icon="lock"
              className={clsx('icon icon-sm icon-success')}
              feather
            />
            Verified
          </div>
        }
        {!verified &&
          <div className="vertical-center space-right">
            <Icon
              icon="lock"
              className={clsx('icon icon-sm icon-danger')}
              feather
            />
            Not verified
          </div>
        }
      </div>
    </div>
  );
}

// function runningContainersFilter(containers: DockerContainer[]) {
//   return containers.filter((container) => container.Status === 'running')
//     .length;
// }
// function stoppedContainersFilter(containers: DockerContainer[]) {
//   return containers.filter((container) => container.Status === 'exited').length;
// }
// function healthyContainersFilter(containers: DockerContainer[]) {
//   return containers.filter((container) => container.Status === 'healthy')
//     .length;
// }
// function unhealthyContainersFilter(containers: DockerContainer[]) {
//   return containers.filter((container) => container.Status === 'unhealthy')
//     .length;
// }
