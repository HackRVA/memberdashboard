import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { LayoutComponent } from './shared/components/layout';
import { TopNavComponent } from './shared/components/top-nav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { AuthService } from './shared/services';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, LayoutComponent, TopNavComponent, MatToolbarModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent implements OnInit {
  isLogin: boolean = false;

  constructor(private readonly authService: AuthService) {}

  ngOnInit(): void {
    this.isLogin = this.authService.user$.getValue().isLogin;
  }
}
