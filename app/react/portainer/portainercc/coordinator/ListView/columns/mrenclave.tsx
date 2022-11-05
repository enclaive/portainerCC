import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const mrenclave: Column<CoordintaorListEntry> = {
  Header: 'MRENCLAVE',
  accessor: (row) => row.mrenclave,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
