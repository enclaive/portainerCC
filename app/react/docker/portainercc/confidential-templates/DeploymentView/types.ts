export interface FormValues {
  Name: string;
  Username: string;
  Password: string;
  ImageID: string,
}

export interface DeployConfidentialTemplateFormValues {
  Id: number;
  EnvId: number;
  Name: string;
  Values: Object;
}
