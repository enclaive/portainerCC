import { PageHeader } from '@@/PageHeader';


import { useEnvironmentId } from '@/portainer/hooks/useEnvironmentId';
import { useConfidentialTemplates } from '../queries';
import { ConfidentialTemplateEntryView } from './ConfidentialTemplateEntryView';
import { useEncryptedVolumes } from '@/react/docker/volumes/queries';

export function ConfidentialTemplatesView() {
    let title = "Confidential Templates";

    let templateQuery = useConfidentialTemplates();
    let env = useEnvironmentId();

    const envId = Number(env);
    let encVolumeQuery = useEncryptedVolumes(envId);

    if (!templateQuery.data || !env || !encVolumeQuery.data) {
        return null;
    }

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {templateQuery.data.map((entry) => <ConfidentialTemplateEntryView key={entry.Id} template={entry} envId={envId} encryptedVolumes={encVolumeQuery.data.Volumes} />)}

        </>
    );

}
