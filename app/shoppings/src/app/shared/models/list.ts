export interface List {
  id: string;
  storeId?: string;
  name: string;
  summary?: string;
  archived: boolean;
}

export interface NewResource {
  id: string;
}
