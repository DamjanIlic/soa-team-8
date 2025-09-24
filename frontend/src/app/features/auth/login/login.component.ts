import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../../core/services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './login.component.html',
})
export class LoginComponent {
  loginForm: FormGroup;
  successMessage = '';
  errorMessage = '';

  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private router: Router
    // Ukloni TokenStorage - ne treba više
  ) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required]
    });
  }

  onSubmit(): void {
    if (this.loginForm.invalid) return;

    this.authService.login(this.loginForm.value).subscribe({
      next: (response) => {
        console.log('Login response:', response);

        if (!response.token) {
          // ako je blokiran ili nije uspeo login
          this.errorMessage = response.message || 'Login failed.';
          this.successMessage = '';
          return;
        }

        // ako token postoji uloguj korisnika
        this.authService.setAuthStatus(response.token);
        this.successMessage = 'Login successful!';
        this.errorMessage = '';
        this.router.navigate(['/profile']);
      },
      error: (err) => {
        console.log('Login error:', err);
        this.errorMessage = err.error?.message || err.error || 'Login failed. Check your credentials.';
        this.successMessage = '';
      }
    });
  }

}