import { Team } from '../../users/teams/types';

export type KeyId = number;

export type KeyEntry = {
  Id: KeyId;
  KeyType?: string;
  Description: string;
  TeamAccessPolicies: any;
  Export?: string;
  AllTeams: Team[];
};
