import { CoordinatorListEntry, CoordinatorImageId } from '@/react/portainer/portainercc/coordinator/types';
import { PortainerSelect } from '@@/form-components/PortainerSelect'


interface Props {
  name?: string;
  value: CoordinatorImageId;
  onChange(value: CoordinatorImageId): void;
  images: CoordinatorListEntry[];
  dataCy?: string;
  inputId?: string;
  placeholder?: string;
}

export function CoordinatorImageSelector({
  name,
  value,
  onChange,
  images,
  dataCy,
  inputId,
  placeholder,
}: Props) {
  const options = images.map((image) => ({ label: image.name, value: image.id }));

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
