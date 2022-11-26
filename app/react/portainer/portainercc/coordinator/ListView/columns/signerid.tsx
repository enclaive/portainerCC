import { Column } from 'react-table';


import { CoordinatorListEntry } from '../../types';

export const signerid: Column<CoordinatorListEntry> = {
  Header: 'SignerID',
  accessor: (row) => row.signerId,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
