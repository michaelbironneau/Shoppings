/* eslint-disable no-console */
import { Observable, of } from 'rxjs';
import { Injectable } from '@angular/core';
import { BaseService } from './base.service';
import { map } from 'rxjs/operators';
import { seedListItemsA, seedListItemsB } from '../seed-data/list-items';
import { ListItem } from '../models/list-item';
import { ListUpdate } from '../models/list-update';
import { autocompleteItems } from '../seed-data/autocomplete';
import { environment } from '../../../environments/environment';
import { Item } from '../models/item';
import { SyncService } from './sync.service';
import { HttpClient } from '@angular/common/http';

if (!environment.api) {
  localStorage.setItem('list-aaaa', JSON.stringify(seedListItemsA));
  localStorage.setItem('list-asdf', JSON.stringify(seedListItemsB));
}

@Injectable({
  providedIn: 'root',
})
export class ListItemService extends BaseService {
  constructor(private sync: SyncService, private http: HttpClient) {
    super();
  }

  trySyncList(listID: string): Observable<boolean> {
    if (!environment.api) {
      return of(true);
    }
    if (environment.api && !this.haveNetworkConnectivity) {
      return of(false); // can't sync list if we don't have a network
    }
    this.http.get<ListUpdate>(`${environment.api}/lists/${listID}/items`).pipe(
      map((items: ListUpdate) => {
        localStorage.setItem(`list-${listID}`, JSON.stringify(items.updates));
        localStorage.setItem(
          `list-${listID}-updated`,
          items.updatedAt.toString()
        );
        return true;
      })
    );
  }

  trySyncAutocompleteList(): Observable<boolean> {
    if (!environment.api) {
      localStorage.setItem('autocomplete', JSON.stringify(autocompleteItems));
      return of(true);
    }
    if (!this.haveNetworkConnectivity) {
      return of(false);
    }
    this.http.get<Item[]>(`${environment.api}/items`).pipe(
      map((items: Item[]) => {
        localStorage.setItem('autocomplete', JSON.stringify(items));
        return true;
      })
    );
  }

  //  get all list items
  getAll(listID: string): Observable<ListItem[]> {
    return this.trySyncList(listID).pipe(
      map((status: boolean) => {
        if (!status) {
          console.warn(`Failed to sync list ${listID}`);
        }
        const listStr = localStorage.getItem(`list-${listID}`);
        if (listStr == null) {
          return [];
        }
        return JSON.parse(listStr);
      })
    );
  }

  searchAutocompleteAPI(needle: string): Observable<ListItem[]> {
    return this.http
      .get<Item[]>(`${environment.api}/item-search/${needle}`)
      .pipe(
        map((items: Item[]) =>
          items.map((item: Item) => ({
            name: item.name,
            id: item.id,
            checked: false,
            listId: null,
            quantity: 1,
          }))
        )
      );
  }

  searchAutocomplete(needle: string): Observable<ListItem[]> {
    if (environment.api && this.haveNetworkConnectivity) {
      // prefer live search if available
      return this.searchAutocompleteAPI(needle);
    }
    let searchSpace: Item[] = [];

    const storedResults = localStorage.getItem('autocomplete');
    if (storedResults) {
      searchSpace = JSON.parse(storedResults);
    }

    console.debug(`Searching for ${needle}`);
    const hits = searchSpace.filter(
      (item: Item) =>
        item.name.toLowerCase().indexOf(needle.toLowerCase()) !== -1
    );
    const itemsToReturn = hits.map((hit: Item) => ({
      name: hit.name,
      id: hit.id,
      checked: false,
      listId: null,
      quantity: 1,
    }));
    console.debug(`Returning ${itemsToReturn.length} hits`);
    return of(itemsToReturn);
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
    if (!environment.api) {
      return of(true);
    }
    if (environment.api && this.haveNetworkConnectivity) {
      let update: ListUpdate;
      update.updates = [items[existingItem]];
      return this.http
        .patch(`${environment.api}/lists/${listID}/updates`, update)
        .pipe(map(() => true));
    }
    // no network connectivity - enqueue
    this.sync.enqueueUpdate(items[existingItem]);
    return of(true);
  }

  _updateItem(items: ListItem[], update: ListItem) {
    const itemIndex = items.findIndex(
      (item) =>
        (update.name && item.name && update.name === item.name) ||
        (update.id && update.id && update.id === item.id)
    );
    if (itemIndex === -1) {
      items.push(update);
      return;
    }
    items[itemIndex] = update;
  }

  syncListItems(listID: string): Observable<ListItem[]> {
    // sync locally and return any new items from the API
    if (!environment.api || !this.haveNetworkConnectivity) {
      return of([]);
    }
    const storageKey = `list-${listID}`;
    const updateKey = `list-${listID}-updated`;
    this.http
      .get(`${environment.api}/lists/${listID}/updates/${updateKey}`)
      .pipe(
        map((update: ListUpdate) => {
          if (update.updatedAt === 0) {
            return [];
          }
          const newUpdateValue = update.updatedAt.toString();
          localStorage.setItem(updateKey, newUpdateValue);
          const listStr = localStorage.getItem(storageKey);
          if (listStr === null) {
            localStorage.setItem(storageKey, JSON.stringify(update.updates)); // new list, or no items written yet
            return update.updates;
          }
          const listItems: ListItem[] = JSON.parse(listStr);
          update.updates.forEach((item: ListItem) => {
            // eslint-disable-next-line no-underscore-dangle
            this._updateItem(listItems, item);
          });
          localStorage.setItem(storageKey, JSON.stringify(listItems));
          return update.updates;
        })
      );
  }

  // apply update to the list and return if successful
  applyUpdate(listID: string, update: ListUpdate): Observable<boolean> {
    const storageKey = `list-${listID}`;
    const updateKey = `list-${listID}-updated`;
    const listStr = localStorage.getItem(storageKey);
    if (!listStr) {
      // we should do a full update if we have network connectivity,
      // but for now let's just update based on these items.
      console.warn(`Updating unknown list ${listID}`);
      localStorage.setItem(storageKey, JSON.stringify(update.updates));
      localStorage.setItem(updateKey, update.updatedAt.toString());
    }

    let listItems: ListItem[] = JSON.parse(listStr);
    update.updates.forEach((item: ListItem) => {
      // eslint-disable-next-line no-underscore-dangle
      this._updateItem(listItems, item);
    });
    listItems = listItems.filter((item: ListItem) => item.quantity > 0); // delete if quantity == 0
    localStorage.setItem(storageKey, JSON.stringify(listItems));
    if (environment.api && this.haveNetworkConnectivity) {
      return this.http
        .patch(`${environment.api}/lists/${listID}/updates`, update)
        .pipe(map(() => true));
    } else if (environment.api && !this.haveNetworkConnectivity) {
      update.updates.forEach((item: ListItem) => {
        this.sync.enqueueUpdate(item);
      });
    }
    return of(true);
  }
}
