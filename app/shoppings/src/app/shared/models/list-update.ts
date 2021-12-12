import { ListItem } from './list-item';

export interface ListUpdate {
  updatedAt: number;
  updates: ListItem[];
}
