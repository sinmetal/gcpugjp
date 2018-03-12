export interface OrganizationListAPIResponse {
    list : Organization[];
    hasNext : boolean;
    cursor : string;
}

export interface Organization {
    key: string;
    name: string;
    url: string;
    logoUrl: string;
    order: number;
    createdAt: Date;
    updatedAt: Date;
}
