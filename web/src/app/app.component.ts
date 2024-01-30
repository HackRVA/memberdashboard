import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MDLayoutComponent } from './shared/components/md-layout';
import { TopNavComponent } from './shared/components/top-nav';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, MDLayoutComponent, TopNavComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {}
