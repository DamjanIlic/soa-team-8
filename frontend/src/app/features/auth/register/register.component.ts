import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { AuthService } from '../../../core/services/auth.service';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

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
    private router: Router
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
        localStorage.setItem('token', res.token);
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

        fetch('http://localhost:8000/api/stakeholders', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${res.token}`
          },
          body: JSON.stringify(stakeholderData)
        })
        .then(response => {
          if (!response.ok) throw new Error('Failed to create stakeholder profile');
          return response.json();
        })
        .then(data => {
          console.log('Stakeholder profile created:', data);
        })
        .catch(err => {
          console.error('Error creating stakeholder:', err);
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