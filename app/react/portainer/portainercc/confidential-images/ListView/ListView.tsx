
import { PageHeader } from '@@/PageHeader';

const storageKey = 'portainercc-confimages';

export function ConfidentialImagesListView() {



    let title = "Confidential images";

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

        </>
    );
}
