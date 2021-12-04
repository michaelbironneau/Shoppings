export interface ListItem {
  id?: string;
  listId: string;
  name: string;
  quantity: number;
  checked: boolean;
  storeOrder?: number;
}
