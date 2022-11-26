import { Column } from 'react-table';


import { CoordinatorDeploymentEntry } from '../../types';

export const id: Column<CoordinatorDeploymentEntry> = {
  Header: 'ID',
  accessor: (row) => row.id,
  disableFilters: true,
  canHide: false,
  sortType: 'number',
};
