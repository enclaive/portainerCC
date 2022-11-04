import { KeyEntry, KeyId } from '@/react/portainer/portainercc/keymanagement/types';
import { PortainerSelect } from '@@/form-components/PortainerSelect'


interface Props {
  name?: string;
  value: KeyId[] | readonly KeyId[];
  onChange(value: readonly KeyId[]): void;
  keys: KeyEntry[];
  dataCy?: string;
  inputId?: string;
  placeholder?: string;
}

export function TeamsSelector({
  name,
  value,
  onChange,
  keys,
  dataCy,
  inputId,
  placeholder,
}: Props) {
  const options = keys.map((key) => ({ label: key.Id + ": " + key.Description, value: key.Id }));

  return (
    <PortainerSelect<number>
      name={name}
      isMulti
      options={options}
      value={value}
      onChange={(value) => onChange(value)}
      data-cy={dataCy}
      inputId={inputId}
      placeholder={placeholder}
    />
  );
}
