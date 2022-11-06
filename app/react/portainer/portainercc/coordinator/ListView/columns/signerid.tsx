import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const signerid: Column<CoordintaorListEntry> = {
  Header: 'SignerID',
  accessor: (row) => row.signerId,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
