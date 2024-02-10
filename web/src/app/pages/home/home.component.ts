import { Component } from '@angular/core';
import { MinionComponent } from '../../shared/components/minion';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'md-home',
  standalone: true,
  imports: [MatButtonModule, MinionComponent],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent {}
