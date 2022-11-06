import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const uniqueid: Column<CoordintaorListEntry> = {
  Header: 'UniqueID',
  accessor: (row) => row.uniqueId,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
