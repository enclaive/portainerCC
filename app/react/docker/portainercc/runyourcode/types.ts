export type CoordinatorDeploymentID = number;

export type ServiceAdd = {
    EnvironmentID: number,
    Name: string,
    Username: string,
    Password: string,
}

export type ServiceDeploy = {
    EnvironmentID: number,
    Name: string,
    ImageID: string,
}

export type ConfidentialTemplate = {
    Id: number,
    Image: string,
    LogoURL: string,
    TemplateName: string,
    Inputs: Array<Input>;
}

export type Input = {
    Label: string,
    Default: string,
    Type: string,
    SecretName?: string,
    ReplacePattern?: string,
    PortContainer?: string,
    PortType?: string
}