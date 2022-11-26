import { Column } from 'react-table';


import { CoordinatorDeploymentEntry } from '../../types';

export const environment: Column<CoordinatorDeploymentEntry> = {
  Header: 'Environment',
  accessor: (row) => row.endpointId,
  disableFilters: true,
  canHide: false,
  sortType: 'number',
};
