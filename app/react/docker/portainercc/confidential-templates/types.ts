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