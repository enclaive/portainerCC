
import { PageHeader } from '@@/PageHeader';
import { BuildCoordinatorForm } from './BuildCoordinatorForm/BuildCoordinatorForm';
import { useKeys } from '../../keymanagement/queries';
import { Datatable } from '@@/datatables';
import { columns } from './columns';


export function CoordinatorImagesListView() {

    const keysQuery = useKeys('SIGNING')
    const coordintaorQuery = null;

    let title = "Coordinator images";


    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {keysQuery.data && (
                <BuildCoordinatorForm keys={keysQuery.data} />
            )}

            {/* <Datatable
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
            /> */}
        </>
    );
}
