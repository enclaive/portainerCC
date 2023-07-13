export interface FormValues {
  Type: string;
  EnvId: number;
  SigningKeyId: number;
  Name: string;
  Ports: Array<{Type: string, Host: string, Container: string}>;
  Repository: string;
  BuildArgs: string;
  RunArgs: string;
}