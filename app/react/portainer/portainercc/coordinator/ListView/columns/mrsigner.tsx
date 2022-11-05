import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const mrsigner: Column<CoordintaorListEntry> = {
  Header: 'MRSIGNER',
  accessor: (row) => row.mrsigner,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
