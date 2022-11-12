
import { PageHeader } from '@@/PageHeader';
import { BuildCoordinatorForm } from './BuildCoordinatorForm/BuildCoordinatorForm';
import { useKeys } from '../../keymanagement/queries';
import { Datatable } from '@@/datatables';
import { columns } from './columns';
import { Codesandbox } from 'react-feather';
import { createStore } from '@/react/portainer/environments/update-schedules/ListView/datatable-store';
import { useCoordinatorImages } from '../queries';

const storageKey = 'portainercc-coordinators';
const useStore = createStore(storageKey);

export function CoordinatorImagesListView() {

    const store = useStore();

    const keysQuery = useKeys('SIGNING')
    const coordintaorQuery = useCoordinatorImages();

    let title = "Coordinator images";

    if (!coordintaorQuery.data) {
        return null;
    }

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {keysQuery.data && (
                <BuildCoordinatorForm keys={keysQuery.data} />
            )}

            
            <Datatable
                columns={columns}
                titleOptions={{
                    title: title,
                    icon: Codesandbox,
                }}
                disableSelect
                dataset={coordintaorQuery.data}
                settingsStore={store}
                storageKey={storageKey}
                emptyContentLabel="No keys found"
                // isLoading={listQuery.isLoading}
                totalCount={coordintaorQuery.data.length}
            // renderTableActions={(selectedRows) => (
            //     <TableActions selectedRows={selectedRows} />
            // )}
            />
        </>
    );
}
