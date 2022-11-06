import clsx from 'clsx';

import { Icon } from '@/react/components/Icon';

// interface Props {
//   containers: DockerContainer[];
// }

export function useCoordinatorStatusComponent() {
  return <CoordinatorStatus/>;
}

export function CoordinatorStatus() {
  return (
    <div className="pull-right">
      <div>
        <div className="vertical-center space-right pr-5">
          <Icon
            icon="power"
            className={clsx('icon icon-sm icon-success')}
            feather
          />
          1 running
        </div>
        <div className="vertical-center space-right">
          <Icon
            icon="power"
            className={clsx('icon icon-sm icon-danger')}
            feather
          />
          1 stopped
        </div>
      </div>
      <div>
        <div className="vertical-center space-right pr-5">
          <Icon
            icon="heart"
            className={clsx('icon icon-sm icon-success')}
            feather
          />
          1 healthy
        </div>
        <div className="vertical-center space-right">
          <Icon
            icon="heart"
            className={clsx('icon icon-sm icon-danger')}
            feather
          />
           1 unhealthy
        </div>
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
