import { Component } from '@angular/core';
import { NgFor, NgIf } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [NgFor, NgIf, FormsModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  url: string = '';
  downloading: string[] = [];
  downloaded: string[] = [];

  startDownload() {
    if (!this.url.trim()) return;

    // Simulate download
    const title = `Downloading: ${this.url}`;
    this.downloading.push(title);

    const index = this.downloading.length - 1;

    // Simulate completion after 3 seconds
    setTimeout(() => {
      const finished = this.downloading[index];
      this.downloading.splice(index, 1);
      this.downloaded.unshift(finished.replace('Downloading:', 'Downloaded:'));
    }, 3000);

    this.url = '';
  }
}
