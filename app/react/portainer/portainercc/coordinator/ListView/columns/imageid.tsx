import { Column } from 'react-table';


import { CoordintaorListEntry } from '../../types';

export const imageid: Column<CoordintaorListEntry> = {
  Header: 'ImageID',
  accessor: (row) => row.imageId,
  disableFilters: true,
  canHide: false,
  sortType: 'string',
};
