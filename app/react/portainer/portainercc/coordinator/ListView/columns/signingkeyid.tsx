import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const signingkeyid: Column<CoordintaorListEntry> = {
  Header: 'SigningKeyID',
  accessor: (row) => row.signingKeyId,
  disableFilters: true,
  canHide: false,
  sortType: 'number',
};
