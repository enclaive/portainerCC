export interface MarbleManifest {
  Packages: Record<string, Package>,
  Marbles: Record<string, any>
}

export interface Package {
  UniqueID?: string,
  SignerID?: string,
  ProductID?: number,
  SecurityVersion?: number
  Debug?: boolean,
}

export interface FormValues {
  coordinatorImageId: number;
  verify: boolean;
}
