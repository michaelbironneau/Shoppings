/* eslint-disable no-console */
import { Observable, of } from 'rxjs';
import { Injectable } from '@angular/core';
import { BaseService } from './base.service';
import { List } from '../models/list';
import { map } from 'rxjs/operators';
import { seedListItemsA, seedListItemsB } from '../seed-data/list-items';
import { ListItem } from '../models/list-item';
import { ListUpdate } from '../models/list-update';

// run once
localStorage.setItem('list-aaaa', JSON.stringify(seedListItemsA));
localStorage.setItem('list-asdf', JSON.stringify(seedListItemsB));

@Injectable({
  providedIn: 'root',
})
export class ListItemService extends BaseService {
  //  get all list items
  getAll(listID: string): Observable<ListItem[]> {
    const listStr = localStorage.getItem(`list-${listID}`);
    if (listStr == null) {
      return of([]);
    }
    return of(JSON.parse(listStr));
  }

  completeItem(listID: string, item: ListItem): Observable<boolean> {
    const listStr = localStorage.getItem(`list-${listID}`);
    if (listStr == null) {
      console.warn(
        `Tried to complete item ${item.name} from unknown list ${listID}`
      );
      return of(false);
    }
    const items: ListItem[] = JSON.parse(listStr);
    const existingItem = items.findIndex((it) => item.name === it.name);
    if (existingItem === -1) {
      console.warn(`Tried to complete item which did not exist ${item.name}`);
      return of(false);
    }
    if (items[existingItem].checked) {
      console.debug(`Item ${item.name} was already complete`);
      return of(true);
    }
    console.debug(`Checked item ${item.name} on list ${listID}`);
    items[existingItem].checked = true;
    localStorage.setItem(`list-${listID}`, JSON.stringify(items));
    return of(true);
  }

  // apply update to the list and return if successful
  applyUpdate(listID: string, update: ListUpdate): Observable<boolean> {
    const listStr = localStorage.getItem(`list-${listID}`);
    if (listStr == null && update.quantityDiff <= 0) {
      console.warn(
        `Tried to remove item ${update.name} from unknown list ${listID}`
      );
      return of(false);
    }
    if (listStr == null && update.quantityDiff > 0) {
      console.debug(`Creating new item list for ${listID}`);
      console.debug(`Adding ${update.quantityDiff} of ${update.name}`);
      const newItems: ListItem[] = [
        {
          listId: listID,
          name: update.name,
          quantity: update.quantityDiff,
          checked: false,
        },
      ];
      localStorage.setItem(`list-${listID}`, JSON.stringify(newItems));
      return of(true);
    }
    const items: ListItem[] = JSON.parse(listStr);
    const existingItem = items.findIndex((item) => item.name === update.name);

    if (existingItem === -1 && update.quantityDiff > 0) {
      console.debug(`Adding new item to list ${listID}: ${update.name}`);
      items.push({
        listId: listID,
        name: update.name,
        quantity: update.quantityDiff,
        checked: false,
      });
      return of(true);
    }

    if (existingItem === -1 && update.quantityDiff <= 0) {
      console.debug(
        `Update for item ${update.name} on list ${listID} had no effect, with quantity diff ${update.quantityDiff}`
      );
      return of(true);
    }

    items[existingItem].quantity += update.quantityDiff;
    console.debug(
      `Updated item ${update.name} on list ${listID} to quantity ${items[existingItem].quantity}`
    );

    if (items[existingItem].quantity <= 0) {
      console.debug(`Removed item ${update.name} on list ${listID}`);
      items.splice(existingItem, 1);
    }
    localStorage.setItem(`list-${listID}`, JSON.stringify(items));

    return of(true);
  }
}
