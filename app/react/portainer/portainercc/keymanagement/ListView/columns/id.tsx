import { Column } from 'react-table';
import { KeyEntry } from '../../types';


export const id: Column<KeyEntry> = {
  Header: 'Id',
  accessor: (row) => row.Id,
  disableFilters: true,
  Filter: () => null,
  canHide: false,
};
