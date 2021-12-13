import { List } from './list';
import { ListItem } from './list-item';

export const PRIORITY_UPDATE_ARCHIVE = 'archive';
export const PRIOTITY_UPDATE_NEW = 'new';

export interface ListUpdate {
  updatedAt: number;
  updates: ListItem[];
}

export interface PriorityQueueItem {
  updateType: string;
  list: List;
}
