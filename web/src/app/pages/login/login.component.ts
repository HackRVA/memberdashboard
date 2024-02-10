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
import { Observable } from 'rxjs';
import { AuthService, LocalStorageService } from '../../shared/services';
import { AuthResponse } from '../../shared/types';

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
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
})
export class LoginComponent {
  showLoginForm: boolean = true;

  loginFormGroup: FormGroup = new FormGroup({
    email: new FormControl<string>(null, [
      Validators.required,
      Validators.email,
    ]),
    password: new FormControl<string>(null, [Validators.required]),
  });

  registerFormGroup: FormGroup = new FormGroup({
    email: new FormControl<string>(null, [
      Validators.required,
      Validators.email,
    ]),
    password: new FormControl<string>(null, [Validators.required]),
    retypePassword: new FormControl<string>(null, [Validators.required]),
  });

  constructor(
    private readonly authService: AuthService,
    private readonly localStorageService: LocalStorageService
  ) {}

  toggleForm(): void {
    this.showLoginForm = !this.showLoginForm;
  }

  submit(): void {
    const auth$: Observable<AuthResponse | void> = this.showLoginForm
      ? this.authService.login(this.loginFormGroup.value)
      : this.authService.register(this.registerFormGroup.value);

    auth$.subscribe({
      next: (response: AuthResponse | void) => {
        if (response) {
          this.localStorageService.upsert<string>(
            this.authService.getSessionKey(),
            response.token
          );
          window.location.reload();
        }
      },
    });
  }
}
