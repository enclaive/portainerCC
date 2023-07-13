import { useQuery } from "react-query";
import { getConfidentialTemplates } from "./confidential-templates.service";
import { ConfidentialTemplate } from "./types";


export function useConfidentialTemplates<T = ConfidentialTemplate[]>(
  enabled = true,
  select: (data: ConfidentialTemplate[]) => T = (data) => data as unknown as T
) {
  const templates = useQuery(
    ['confidential-templates'],
    () => getConfidentialTemplates(),
    {
      meta: {
        error: { title: 'Failure', message: 'Unable to load confidential templates' },
      },
      enabled,
      select,
    }
  );

  return templates;
}