
import { PageHeader } from '@@/PageHeader';

import { CreateKeyForm } from './CreateKeyForm/CreateKeyForm';

import { useKeyTypeParam } from './useKeyTypeParam';
import { useTeams } from '../../../users/teams/queries';
import { useUsers } from '@/portainer/users/queries';
import { useUser } from '@/portainer/hooks/useUser';

export function KeysView() {
    const { isAdmin } = useUser();


    const type = useKeyTypeParam();
    const usersQuery = useUsers(false);
    const teamsQuery = useTeams(!isAdmin, 0, { enabled: !!usersQuery.data });
    let title = "";

    if (type == "signing") {
        title = "SGX Signing Keys"
    } else if (type == "pf") {
        title = "Gramine Protected File Keys"
    } else {
        throw Error("invalid key type")
    }

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {teamsQuery.data && (
                <CreateKeyForm keytype={type} teams={teamsQuery.data} />
            )}
        </>
    );
}
