import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { MatMenuModule } from '@angular/material/menu';
import { NavigationLink } from '../../types/navigation-link';
import { AuthService } from '../../services';

@Component({
    selector: 'top-nav',
    imports: [
        CommonModule,
        RouterModule,
        MatToolbarModule,
        MatTabsModule,
        MatButtonModule,
        MatIconModule,
        MatMenuModule,
    ],
    templateUrl: './top-nav.component.html',
    styleUrl: './top-nav.component.scss'
})
export class TopNavComponent {
  private _adminNav: NavigationLink[] = [
    { routeName: 'Home', routeLink: 'home', routeIcon: 'home' },
    { routeName: 'User', routeLink: 'user', routeIcon: 'person' },
    { routeName: 'Reports', routeLink: 'report', routeIcon: 'show_chart' },
    { routeName: 'Members', routeLink: 'member', routeIcon: 'people' },
    { routeName: 'Resources', routeLink: 'resource', routeIcon: 'devices' },
  ];

  navLinks: NavigationLink[] = [
    { routeName: 'Home', routeLink: 'home', routeIcon: 'home' },
    { routeName: 'User', routeLink: 'user', routeIcon: 'person' },
  ];

  username: string = '';

  constructor(private readonly authService: AuthService) {
    this.runBeforeInit();
  }

  toggleTheme(): void {
    document.body.classList.toggle('dark-theme');
  }

  logout(): void {
    this.authService.logout().subscribe((_) => {
      window.location.reload();
    });
  }

  private runBeforeInit(): void {
    this.username = this.authService.user$.getValue().email;
    const isAdmin: boolean = this.authService.user$.getValue().isAdmin;

    if (isAdmin) {
      this.navLinks = this._adminNav;
    }
  }
}
