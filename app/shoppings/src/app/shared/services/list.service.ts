import { Observable, of } from 'rxjs';
import { Injectable } from '@angular/core';
import { List, NewResource } from '../models/list';
import { map } from 'rxjs/operators';
import { seedLists } from '../seed-data/lists';
import { environment } from 'src/environments/environment';
import { HttpClient } from '@angular/common/http';
import { SyncService } from './sync.service';

@Injectable({
  providedIn: 'root',
})
export class ListService {
  constructor(private http: HttpClient, private sync: SyncService) {}
  trySyncLists(): Observable<boolean> {
    if (!environment.api) {
      localStorage.setItem('lists', JSON.stringify(seedLists));
      return of(true);
    }
    if (!this.sync.haveNetworkConnectivity.getValue()) {
      console.log('No network - not syncing lists');
      return of(false);
    }
    return this.http.get<List[]>(`${environment.api}/lists`).pipe(
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
    if (environment.api && !this.sync.haveNetworkConnectivity.getValue()) {
      throw Error('You must be connected to the internet to add lists');
    } else if (
      environment.api &&
      this.sync.haveNetworkConnectivity.getValue()
    ) {
      return this.http.post<NewResource>(`${environment.api}/lists`, list).pipe(
        map((newR: NewResource) => {
          list.id = newR.id;
          const listsStr = localStorage.getItem('lists');
          let lists: List[];
          if (!listsStr) {
            lists = [list];
            localStorage.setItem('lists', JSON.stringify([list]));
            return newR.id;
          } else {
            lists = JSON.parse(listsStr);
          }
          lists.push(list);
          localStorage.setItem('lists', JSON.stringify(lists));
          return newR.id;
        })
      );
    }

    // local "demo" version
    return this.getAll().pipe(
      map((lists: List[]) => {
        let newID: string;
        if (!environment.api) {
          newID = new Date().getTime().toString();
        }
        const newList = { ...list, id: newID };
        lists.push(newList);
        localStorage.setItem('lists', JSON.stringify(lists));
        return newID;
      })
    );
  }

  archiveLocal(listID: string): boolean {
    // remove any list items from local storage
    localStorage.removeItem(`list-${listID}`);
    const listsStr = localStorage.getItem('lists');
    let lists: List[];
    if (!listsStr) {
      lists = [];
    } else {
      lists = JSON.parse(listsStr);
    }

    const ix = lists.findIndex((list) => list.id === listID);
    if (ix === -1) {
      return false;
    }
    lists[ix].archived = true;
    lists.splice(ix, 1);
    localStorage.setItem('lists', JSON.stringify(lists));
    return true;
  }

  // archive the given list by ID
  archive(listID: string): Observable<boolean> {
    if (environment.api && !this.sync.haveNetworkConnectivity.getValue()) {
      throw Error('You must be connected to the internet to archive a list');
    } else if (
      environment.api &&
      this.sync.haveNetworkConnectivity.getValue()
    ) {
      return this.http
        .post(`${environment.api}/lists/${listID}/archive`, null)
        .pipe(map(() => this.archiveLocal(listID)));
    }

    // local version
    return of(this.archiveLocal(listID));
  }
}
