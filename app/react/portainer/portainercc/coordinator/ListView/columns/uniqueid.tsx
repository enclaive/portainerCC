import { Column } from 'react-table';


import { CoordinatorListEntry } from '../../types';

export const uniqueid: Column<CoordinatorListEntry> = {
  Header: 'UniqueID',
  accessor: (row) => row.uniqueId,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
