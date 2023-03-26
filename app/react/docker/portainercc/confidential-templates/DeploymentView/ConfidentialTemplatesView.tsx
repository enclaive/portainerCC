import { PageHeader } from '@@/PageHeader';


import { useEnvironmentId } from '@/portainer/hooks/useEnvironmentId';
import { useConfidentialTemplates } from '../queries';
import { ConfidentialTemplateEntryView } from './ConfidentialTemplateEntryView';
import { useEncryptedVolumes } from '@/react/docker/volumes/queries';

export function ConfidentialTemplatesView() {

    let templateQuery = useConfidentialTemplates();

    if (!templateQuery.data) {
        return null;
    }


    let title = "Confidential Templates";

    let env = useEnvironmentId();
    if (!env) {
        return null;
    }

    const envId = Number(env);

    //encrypted volumes
    let encVolumeQuery = useEncryptedVolumes(envId);
    if (!encVolumeQuery.data) {
        return null;
    }


    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {templateQuery.data.map((entry) => <ConfidentialTemplateEntryView key={entry.Id} template={entry} envId={envId} encryptedVolumes={encVolumeQuery.data.Volumes} />)}


        </>
    );

}
