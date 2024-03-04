import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { AuthService, LocalStorageService } from '@md-shared/services';
import { AuthResponse } from '@md-shared/types';
import { passwordMatchValidator } from './validator';

@Component({
  selector: 'md-login',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatSnackBarModule,
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
})
export class LoginComponent {
  showLoginForm: boolean = true;
  hasFailed: boolean = false;

  loginFormGroup: FormGroup = new FormGroup({
    email: new FormControl<string>(null, [
      Validators.required,
      Validators.email,
    ]),
    password: new FormControl<string>(null, [Validators.required]),
  });

  registerFormGroup: FormGroup = new FormGroup(
    {
      email: new FormControl<string>(null, [
        Validators.required,
        Validators.email,
      ]),
      password: new FormControl<string>(null, [Validators.required]),
      retypePassword: new FormControl<string>(null, [Validators.required]),
    },
    { validators: passwordMatchValidator }
  );

  constructor(
    private readonly authService: AuthService,
    private readonly localStorageService: LocalStorageService,
    private readonly snackBar: MatSnackBar
  ) {}

  toggleForm(): void {
    this.showLoginForm = !this.showLoginForm;
  }

  login(): void {
    this.authService.login(this.loginFormGroup.value).subscribe({
      next: (response: AuthResponse) => {
        if (response) {
          this.localStorageService.upsert<string>(
            this.authService.getSessionKey(),
            response.token
          );
          window.location.reload();
        }
      },
      error: (_) => {
        this.loginFormGroup.setErrors({
          apiError: 'Email and/or password are incorrect.',
        });
      },
    });
  }

  register(): void {
    this.authService.register(this.registerFormGroup.value).subscribe({
      next: () => {
        this.snackBar.open("You're all set!", '', { duration: 3000 });
      },
      error: () => {
        this.snackBar.open(
          'Hrmmm, are you sure everything in the form is correct?',
          '',
          { duration: 3000 }
        );
      },
    });
  }
}
