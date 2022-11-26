import { Column } from 'react-table';


import { CoordinatorListEntry } from '../../types';

export const name: Column<CoordinatorListEntry> = {
  Header: 'Name',
  accessor: (row) => row.name,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
