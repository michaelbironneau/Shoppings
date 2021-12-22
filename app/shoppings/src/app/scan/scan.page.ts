import { Component, OnDestroy, OnInit } from '@angular/core';
import * as Tesseract from 'tesseract.js';
import { createWorker } from 'tesseract.js';

@Component({
  selector: 'app-scan',
  templateUrl: './scan.page.html',
  styleUrls: ['./scan.page.scss'],
})
export class ScanPage implements OnInit, OnDestroy {
  content = null;
  worker = createWorker({
    logger: (m) => console.log(m), // Add logger here
  });
  constructor() {}

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

  onScan(filename: string) {
    console.log('Scanning', filename);
    this.worker.recognize(`../assets/test-images/${filename}`).then((data) => {
      console.log(data);
      this.content = data.data.text;
    });
  }
}
