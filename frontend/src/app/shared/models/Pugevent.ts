export interface PugeventListAPIResponse {
  list : Pugevent[];
  hasNext : boolean;
  cursor : string;
}

export interface Pugevent {
  id: string;
  organizationId: string;
  title: string;
  description: string;
  url: string;
  startAt: Date;
  createdAt: Date;
  updatedAt: Date;
}
