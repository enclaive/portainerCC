import { Column } from 'react-table';
import { KeyEntry } from '../../types';


export const teamAccess: Column<KeyEntry> = {
  Header: 'Access',
  accessor: (row) => "row.TeamAccessPolicies",
  disableFilters: true,
  Filter: () => null,
  canHide: false,
};
