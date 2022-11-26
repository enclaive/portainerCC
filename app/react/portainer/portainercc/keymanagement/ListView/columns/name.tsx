import { CellProps, Column } from 'react-table';

import { Link } from '@@/Link';

import { KeyEntry } from '../../types';

export const name: Column<KeyEntry> = {
  Header: 'Description',
  accessor: 'Description',
  id: 'Description',
  Cell: NameCell,
  disableFilters: true,
  Filter: () => null,
  canHide: false,
  sortType: 'string',
};

export function NameCell({ value: name, row }: CellProps<KeyEntry>) {
  return (
    <Link to="#" params={{ id: row.original.Id }}>
      {name}
    </Link>
  );
}
