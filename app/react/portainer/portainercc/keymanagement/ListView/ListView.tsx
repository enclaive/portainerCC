
import { PageHeader } from '@@/PageHeader';
import { Key } from 'react-feather';

import { CreateKeyForm } from './CreateKeyForm/CreateKeyForm';
import { Datatable } from '@@/datatables';
import { columns } from './columns';
import { createStore } from '@/react/portainer/environments/update-schedules/ListView/datatable-store';

import { useKeyTypeParam } from './useKeyTypeParam';
import { useTeams } from '@/react/portainer/users/teams/queries';
import { useUsers } from '@/portainer/users/queries';
import { useUser } from '@/portainer/hooks/useUser';
import { useKeys } from '../queries';

const storageKey = 'portainercc-keys';
const useStore = createStore(storageKey);

export function ListView() {
    const { isAdmin } = useUser();

    const store = useStore();

    const type = useKeyTypeParam();
    const usersQuery = useUsers(false);
    const teamsQuery = useTeams(!isAdmin, 0, { enabled: !!usersQuery.data });
    const keysQuery = useKeys(type)
    let title = "";

    if (type == "SIGNING") {
        title = "SGX Signing Keys"
    } else if (type == "FILE_ENC") {
        title = "Gramine Protected File Keys"
    } else {
        throw Error("invalid key type")
    }

    if(!keysQuery.data){
        return null;
    }
    console.log(keysQuery)

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {teamsQuery.data && (
                <CreateKeyForm keytype={type} teams={teamsQuery.data} />
            )}

            <Datatable
                columns={columns}
                titleOptions={{
                    title: title,
                    icon: Key,
                }}
                dataset={keysQuery.data}
                settingsStore={store}
                storageKey={storageKey}
                emptyContentLabel="No keys found"
                // isLoading={listQuery.isLoading}
                totalCount={keysQuery.data.length}
                // renderTableActions={(selectedRows) => (
                //     <TableActions selectedRows={selectedRows} />
                // )}
            />
        </>
    );
}
