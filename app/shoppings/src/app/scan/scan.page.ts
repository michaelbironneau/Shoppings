import { Component, OnDestroy, OnInit } from '@angular/core';
import { createWorker, Word } from 'tesseract.js';
import * as commonFoods from './common-foods.json';

@Component({
  selector: 'app-scan',
  templateUrl: './scan.page.html',
  styleUrls: ['./scan.page.scss'],
})
export class ScanPage implements OnInit, OnDestroy {
  content = null;
  foodCache: Set<string> = new Set();
  listItems: Set<string> = new Set();
  worker = createWorker({
    logger: (m) => console.log(m), // Add logger here
  });
  constructor() {
    // Disable ESLint as for some reason forEach fails
    // eslint-disable-next-line guard-for-in
    for (const ix in commonFoods) {
      if (commonFoods[ix] === 'blackberries') {
        console.log('have blackberries');
      }
      this.foodCache.add(commonFoods[ix]);
    }
  }

  ngOnInit() {
    this.loadWorker().then(() => {
      console.log('Worker loaded');
    });
  }

  ngOnDestroy() {
    this.destroyWorker().then(() => {
      console.log('Tesseract worker destroyed');
    });
  }

  async loadWorker() {
    await this.worker.load();
    await this.worker.loadLanguage('eng');
    await this.worker.initialize('eng');
  }

  async destroyWorker() {
    await this.worker.terminate();
  }

  scoreLine(words: Word[]): number {
    let score = 0;
    words.forEach((word: Word) => {
      const wordNoPunct = word.text
        .toLowerCase()
        .replace(/[.,\/#!$%\^&\*;:{}=\-_`~()]/g, '');
      if (this.foodCache.has(wordNoPunct)) {
        this.listItems.add(wordNoPunct);
        score++;
      }
    });
    return score;
  }

  onScan(filename: string) {
    console.log('Scanning', filename);
    this.listItems = new Set();
    this.worker.recognize(`../assets/test-images/${filename}`).then((data) => {
      console.log(data);
      const scores = data.data.lines.map(
        (line) => this.scoreLine(line.words) / line.words.length
      );
      console.log('Scores', scores);
      console.log(Array.from(this.listItems));
    });
  }
}
