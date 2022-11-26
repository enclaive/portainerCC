export type CoordinatorDeploymentID = number;

export type CoordinatorDeployment = {
    id: CoordinatorDeploymentID,
    coordinatorId: number,
    endpointId: number,
    rootCert: {
        Type: string,
        Headers: string,
        Bytes: string,
    },
    userCert: {
        Type: string,
        Headers: string,
        Bytes: string,
    },
    userPrivKey: {
        Type: string,
        Headers: string,
        Bytes: string,
    },
    manifest: any,
    verified: boolean
}