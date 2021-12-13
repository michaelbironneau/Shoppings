import { ListItem } from './list-item';

export const PRIORITY_UPDATE_ARCHIVE = 'archive';
export const PRIOTITY_UPDATE_NEW = 'new';

export interface ListUpdate {
  updatedAt: number;
  updates: ListItem[];
}
