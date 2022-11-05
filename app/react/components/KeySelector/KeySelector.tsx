import { KeyEntry, KeyId } from '@/react/portainer/portainercc/keymanagement/types';
import { PortainerSelect } from '@@/form-components/PortainerSelect'


interface Props {
  name?: string;
  value: KeyId;
  onChange(value: KeyId): void;
  keys: KeyEntry[];
  dataCy?: string;
  inputId?: string;
  placeholder?: string;
}

export function KeySelector({
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
      options={options}
      value={value}
      onChange={(value) => {
        if (value)
          onChange(value)
        else
          onChange(0)
      }}
      data-cy={dataCy}
      inputId={inputId}
      placeholder={placeholder}
    />
  );
}
