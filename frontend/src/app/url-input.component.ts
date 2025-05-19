import { Component, EventEmitter, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'url-input',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './url-input.component.html',
  styleUrls: ['./url-input.component.scss'],
})
export class UrlInputComponent {
  url: string = '';

  @Output()
  download = new EventEmitter<string>();

  submit() {
    if (!this.url.trim()) return;
    this.download.emit(this.url);
    this.url = '';
  }
}
