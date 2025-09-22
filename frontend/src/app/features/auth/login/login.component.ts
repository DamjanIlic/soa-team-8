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
        // Koristi AuthService umesto TokenStorage
        this.authService.setAuthStatus(response.token); // ili response.access_token, zavisi šta backend šalje
        
        this.successMessage = 'Login successful!';
        this.errorMessage = '';
        
        // Navigate gde hoćeš
        this.router.navigate(['/profile']); // ili '/dashboard'
      },
      error: (err) => {
        this.errorMessage = err.error?.message || 'Login failed. Check your credentials.';
        this.successMessage = '';
      }
    });
  }
}