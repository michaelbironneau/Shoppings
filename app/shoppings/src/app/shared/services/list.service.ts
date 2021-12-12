import { Observable, of } from 'rxjs';
import { Injectable } from '@angular/core';
import { BaseService } from './base.service';
import { List } from '../models/list';
import { map } from 'rxjs/operators';
import { seedLists } from '../seed-data/lists';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root',
})
export class ListService extends BaseService {
  trySyncLists(): Observable<boolean> {
    if (!environment.api) {
      localStorage.setItem('lists', JSON.stringify(seedLists));
      return of(true);
    }
    if (!this.haveNetworkConnectivity) {
      return of(false);
    }
    this.http.get<List[]>(`${environment.api}/lists`).pipe(
      map((lists: List[]) => {
        localStorage.setItem('lists', JSON.stringify(lists));
        return true;
      })
    );
  }

  //  get all lists
  getAll(): Observable<List[]> {
    return this.trySyncLists().pipe(
      map(() => {
        const listStr = localStorage.getItem('lists');
        if (listStr == null) {
          return [];
        }
        const lists: List[] = JSON.parse(listStr);
        return lists;
      })
    );
  }

  // add list, returning ID
  add(list: List): Observable<string> {
    return this.getAll().pipe(
      map((lists: List[]) => {
        const newID: string = new Date().getTime().toString();
        const newList = { ...list, id: newID };
        lists.push(newList);
        localStorage.setItem('lists', JSON.stringify(lists));
        return newID;
      })
    );
  }

  // archive the given list by ID
  archive(listID: string): Observable<boolean> {
    // remove any list items from local storage
    localStorage.removeItem(`list-${listID}`);
    return this.getAll().pipe(
      map((lists: List[]) => {
        const ix = lists.findIndex((list) => list.id === listID);
        if (ix === -1) {
          return false;
        }
        lists[ix].archived = true;
        lists.splice(ix, 1);
        localStorage.setItem('lists', JSON.stringify(lists));
        return true;
      })
    );
  }
}
