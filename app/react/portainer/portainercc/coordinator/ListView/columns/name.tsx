import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const name: Column<CoordintaorListEntry> = {
  Header: 'Name',
  accessor: (row) => row.name,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
