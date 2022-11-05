import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const timestamp: Column<CoordintaorListEntry> = {
  Header: 'Created at',
  accessor: (row) => row.timestamp,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
