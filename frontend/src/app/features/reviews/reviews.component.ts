import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-reviews',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-6">
      <!-- Admin vidi sve reviews -->
      <div *ngIf="userRole === 'admin'">
        <h1 class="text-2xl font-bold mb-4">All Reviews</h1>
        <!-- Lista svih reviews -->
      </div>

      <!-- Guide vidi reviews za svoje ture -->
      <div *ngIf="userRole === 'guide'">
        <h1 class="text-2xl font-bold mb-4">Reviews for My Tours</h1>
        <!-- Lista reviews za guide-ove ture -->
      </div>

      <!-- Tourist vidi svoje reviews -->
      <div *ngIf="userRole === 'tourist'">
        <h1 class="text-2xl font-bold mb-4">My Reviews</h1>
        <!-- Lista reviews koje je tourist napisao -->
      </div>
    </div>
  `
})
export class ReviewsComponent implements OnInit {
  userRole: string = '';

  ngOnInit(): void {
    this.userRole = this.getUserRoleFromToken();
  }

  getUserRoleFromToken(): string {
    const token = localStorage.getItem('access_token');
    if (token) {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.role || payload.user_type || 'user';
    }
    return '';
  }
}