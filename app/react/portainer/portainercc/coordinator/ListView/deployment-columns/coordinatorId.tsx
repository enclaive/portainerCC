import { Column } from 'react-table';


import { CoordinatorDeploymentEntry } from '../../types';

export const coordinatorId: Column<CoordinatorDeploymentEntry> = {
  Header: 'Coordinator ID',
  accessor: (row) => row.coordinatorId,
  disableFilters: true,
  canHide: false,
  sortType: 'number',
};
