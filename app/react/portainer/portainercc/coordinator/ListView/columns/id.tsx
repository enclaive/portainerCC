import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const id: Column<CoordintaorListEntry> = {
  Header: 'ID',
  accessor: (row) => row.id,
  disableFilters: true,
  canHide: false,
  sortType: 'number',
};
