import { useCallback, useEffect, useState } from 'react';
import { useRouter } from '@uirouter/react';

import { Widget, WidgetBody } from '@/portainer/components/widget';
import { FormControl } from '@/portainer/components/form-components/FormControl';
import { Button } from '@/portainer/components/Button';
import { FormSectionTitle } from '@/portainer/components/form-components/FormSectionTitle';
import { TextTip } from '@/portainer/components/Tip/TextTip';
import { Code } from '@/portainer/components/Code';
import { CopyButton } from '@/portainer/components/Button/CopyButton';
import { Input } from '@/portainer/components/form-components/Input';
import { useUser } from '@/portainer/hooks/useUser';

import { useCreateAccessTokenMutation } from '../queries';

import styles from './CreateAccessToken.module.css';

export function CreateAccessToken() {
  const router = useRouter();
  const [description, setDescription] = useState('');
  const [errorText, setErrorText] = useState('');
  const [accessToken, setAccessToken] = useState('');

  const { user } = useUser();

  useEffect(() => {
    if (description.length === 0) {
      setErrorText('this field is required');
    } else setErrorText('');
  }, [description]);

  const mutation = useCreateAccessTokenMutation();

  const handleSubmit = useCallback(() => {
    if (!user) {
      throw new Error('User is not authenticated');
    }

    mutation.mutate(
      { id: user.Id, description },
      {
        onSuccess(accessToken) {
          setAccessToken(accessToken);
        },
      }
    );
  }, [mutation, description, user]);

  return (
    <Widget>
      <WidgetBody>
        <div>
          <FormControl inputId="input" label="Description" errors={errorText}>
            <Input
              id="input"
              onChange={(e) => setDescription(e.target.value)}
              value={description}
            />
          </FormControl>
          <Button
            disabled={!!errorText || !!accessToken || mutation.isLoading}
            onClick={() => handleSubmit()}
            className={styles.addButton}
          >
            Add access token
          </Button>
        </div>
        {accessToken && (
          <>
            <FormSectionTitle>New access token</FormSectionTitle>
            <TextTip>
              Please copy the new access token. You won&#39;t be able to view
              the token again.
            </TextTip>
            <Code>{accessToken}</Code>
            <CopyButton copyText={accessToken} className={styles.copyButton}>
              Copy access token
            </CopyButton>
            <hr />
            <Button
              type="button"
              onClick={() => router.stateService.go('portainer.account')}
            >
              Done
            </Button>
          </>
        )}
      </WidgetBody>
    </Widget>
  );
}
