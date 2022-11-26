import { Column } from 'react-table';


import { CoordinatorListEntry } from '../../types';

export const signingkeyid: Column<CoordinatorListEntry> = {
  Header: 'SigningKeyID',
  accessor: (row) => row.signingKeyId,
  disableFilters: true,
  canHide: false,
  sortType: 'number',
};
