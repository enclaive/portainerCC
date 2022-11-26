import { Column } from 'react-table';


import { CoordinatorListEntry } from '../../types';

export const id: Column<CoordinatorListEntry> = {
  Header: 'ID',
  accessor: (row) => row.id,
  disableFilters: true,
  canHide: false,
  sortType: 'number',
};
