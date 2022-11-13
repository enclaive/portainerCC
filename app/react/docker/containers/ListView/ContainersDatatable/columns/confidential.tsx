import { CellProps, Column } from 'react-table';
import { Icon } from '@@/Icon';
import { Link } from '@@/Link';
import { useRowContext } from '../RowContext';

import type { DockerContainer } from '@/react/docker/containers/types';

export const confidential: Column<DockerContainer> = {
  Header: 'Confidential',
  accessor: 'Confidential',
  id: 'confidential',
  Cell: ConfCell,
  disableFilters: true,
  canHide: true,
  Filter: () => null,
};


function ConfCell({
  row: { original: container },
}: CellProps<DockerContainer>) {
  const { environment } = useRowContext();
  const containerId = container.Id;
  // const gpusQuery = useContainerGpus(environmentId, containerId);

  //get image from secure images => true
  //if running => true

  if (Math.random() > 0.7) {
    return <>
      <Link
        to="portainer.confimages"
        params={{ val: "test" }}
        title="UNIQUEID: AF9902FB, SIGNERID: ABFA9992"
      >
        <Icon icon="shield" feather size='md' mode='success' />
      </Link>
      <Link
        to="docker.coordinator"
        params={{ endpointId: environment.Id }}
        title="Coordinator"
      >
        <Icon icon="lock" feather size='md' mode='success' />
      </Link>
    </>;
  } else {
    return <></>;
  }

}
