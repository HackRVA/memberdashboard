import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { AuthService } from '@md-shared/services/auth.service';
import { MinionComponent } from '@md-shared/components/minion';
import { HomeService } from './services';

@Component({
    selector: 'md-home',
    imports: [MatButtonModule, MatSnackBarModule, MinionComponent],
    providers: [HomeService],
    templateUrl: './home.component.html',
    styleUrl: './home.component.scss'
})
export class HomeComponent {
  isAdmin: boolean = false;
  isProcessing: boolean = false;
  constructor(
    private readonly authService: AuthService,
    private readonly homeService: HomeService,
    private readonly snackBar: MatSnackBar
  ) {
    this.isAdmin = this.authService.user$.getValue().isAdmin;
  }

  openDoor(): void {
    this.isProcessing = true;
    this.homeService.openResource({ name: 'frontdoor' }).subscribe({
      next: () => {
        this.snackBar.open('Front door is open (probably)', '', {
          duration: 3000,
        });
      },
      error: () => {
        this.snackBar.open('Hrmmmmmm, it failed', '', { duration: 3000 });
      },
      complete: () => {
        this.isProcessing = false;
      },
    });
  }
}
