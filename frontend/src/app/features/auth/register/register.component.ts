import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { AuthService } from '../../../core/services/auth.service';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { TokenStorage } from '../../../core/interceptors/token.service';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './register.component.html'
})
export class RegisterComponent {
  registerForm: FormGroup;
  successMessage = '';
  errorMessage = '';

  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private router: Router,
    private tokenStorage: TokenStorage,
    private http: HttpClient
  ) {
    this.registerForm = this.fb.group({
      username: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required],
      role: ['', Validators.required],
      name: [''],
      surname: [''],
      profile_image: [''],
      biography: [''],
      motto: ['']
    });
  }

  onSubmit(): void {
    if (this.registerForm.invalid) return;

    this.authService.register(this.registerForm.value).subscribe({
      next: (res) => {
        this.tokenStorage.saveAccessToken(res.token); 
        this.successMessage = 'The user is registered!';
        this.errorMessage = '';

        const stakeholderData = {
        user_id: res.userId,  
        name: this.registerForm.value.name,
        surname: this.registerForm.value.surname,
        profile_image: this.registerForm.value.profile_image,
        biography: this.registerForm.value.biography,
        motto: this.registerForm.value.motto
        };

        this.http.post('http://localhost:8000/api/stakeholders', stakeholderData)
          .subscribe({
            next: (data) => {
              console.log('Stakeholder profile created:', data);
            },
            error: (err) => {
              console.error('Error creating stakeholder:', err);
            }
          });

        this.registerForm.reset();
        this.router.navigate(['/dashboard']);
      },
      error: (err) => {
        this.errorMessage = err.error?.message || 'An error occurred during registration.';
        this.successMessage = '';
      }
    });
  }
}