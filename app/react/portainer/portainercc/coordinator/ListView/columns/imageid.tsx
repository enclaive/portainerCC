import { Column } from 'react-table';


import { CoordinatorListEntry } from '../../types';

export const imageid: Column<CoordinatorListEntry> = {
  Header: 'ImageID',
  accessor: (row) => row.imageId,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
