import { Component } from '@angular/core';
import { List } from '../shared/models/list';

@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrls: ['home.page.scss'],
})
export class HomePage {
  lists: List[] = [
    {
      id: 'asdf',
      name: 'Dec 12th List',
      archived: false,
    },
    {
      id: 'aaaa',
      name: 'Dec 7th List',
      archived: false,
    },
  ];

  constructor() {}
}
