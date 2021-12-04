import { Component, OnInit } from '@angular/core';
import { ListItem } from '../shared/models/list-item';

@Component({
  selector: 'app-shop',
  templateUrl: './shop.page.html',
  styleUrls: ['./shop.page.scss'],
})
export class ShopPage implements OnInit {
  items: ListItem[] = [
    {
      listId: 'asdf',
      name: 'Onion',
      quantity: 3,
      checked: false,
    },
    {
      listId: 'asdf',
      name: 'Toothpaste',
      quantity: 1,
      checked: false,
    },
  ];
  constructor() {}

  ngOnInit() {}
}
