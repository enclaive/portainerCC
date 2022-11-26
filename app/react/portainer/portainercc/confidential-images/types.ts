export type ConfidentialImageId = number;

export type ConfidentialImage = {
    id: ConfidentialImageId,
    imageid: string,
    mrsigner: string,
    mrenclave: string,
    timestamp: Date,
}