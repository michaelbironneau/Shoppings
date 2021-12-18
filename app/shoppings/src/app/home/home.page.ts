/* eslint-disable no-console */
import { Component, OnInit } from '@angular/core';
import { List } from '../shared/models/list';
import { ListService } from '../shared/services/list.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrls: ['home.page.scss'],
})
export class HomePage implements OnInit {
  lists: List[] = [];

  constructor(private listService: ListService, private router: Router) {}

  ngOnInit() {
    this.refresh();
  }

  refresh() {
    this.listService.getAll().subscribe((lists) => {
      this.lists = lists;
    });
  }

  onArchive(listID: string) {
    console.debug(`Archiving list ${listID}`);
    this.listService.archive(listID).subscribe((success: boolean) => {
      console.debug(`Success? ${success}`);
      this.refresh();
    });
  }

  onOpen(list: List) {
    this.router.navigate(['/list', list.id]);
  }

  uniqueName(): string {
    // come up with a unique name of form 'Dec 12 List #2'
    const startingPoint =
      new Date().toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
      }) + ' List';
    let current = startingPoint;
    let i = 1;
    while (i++) {
      const ix = this.lists.findIndex((list) => list.name === current);
      if (ix === -1) {
        break;
      }
      current = startingPoint + ' #' + i.toString();
    }
    return current;
  }

  onAdd() {
    console.debug('Adding new list');
    const newList: List = {
      id: '',
      name: this.uniqueName(),
      archived: false,
      summary: '(empty)',
    };
    this.listService.add(newList).subscribe((newID: string) => {
      console.debug(`New list ID: ${newID}`);
      this.refresh();
    });
  }
}
