export type KeyId = number;

export type KeyEntry = {
  Id: KeyId;
  KeyType?: string;
  Description: string;
  TeamAccessPolicies: any;
};