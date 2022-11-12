import { CellProps, Column } from 'react-table';

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
  const containerId = container.Id;
  // const gpusQuery = useContainerGpus(environmentId, containerId);

  //get image from secure images => true
  //if running => true

  if(Math.random() > 0.7){
    return <>R - A</>;
  }

}
