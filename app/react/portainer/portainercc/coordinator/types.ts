export type CoordinatorImageId = number;

export type CoordinatorListEntry = {
    id: number,
    name: string,
    imageId: string,
    signingKeyId: number,
    uniqueId: string,
    signerId: string,
}

export type CoordinatorDeploymentEntry = {
    id: number,
    coordinatorId: number,
    endpointId: number,
    verified: boolean,
}