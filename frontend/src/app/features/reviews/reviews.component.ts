import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReviewService, Review } from '../../core/services/review.service';

@Component({
  selector: 'app-reviews',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="container mx-auto px-4 py-8">
      <h1 class="text-3xl font-bold text-gray-800 mb-6">
        <span *ngIf="userRole === 'admin'">All Reviews</span>
        <span *ngIf="userRole === 'guide'">Reviews for My Tours</span>
        <span *ngIf="userRole === 'tourist'">My Reviews</span>
      </h1>

      <!-- Loading State -->
      <div *ngIf="loading" class="text-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
        <p class="mt-2 text-gray-600">Loading reviews...</p>
      </div>

      <!-- Error State -->
      <div *ngIf="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
        <p class="text-red-700">{{ error }}</p>
      </div>

      <!-- Reviews Grid -->
      <div class="grid gap-4" *ngIf="!loading && !error && reviews.length > 0">
        <div *ngFor="let review of reviews" class="bg-white rounded-lg shadow-md p-6 border-l-4 border-blue-500">
          <div class="flex justify-between items-start mb-2">
            <div class="flex items-center">
              <div class="flex text-yellow-400 mr-2">
                <svg *ngFor="let star of [1,2,3,4,5]" 
                     [class]="'w-4 h-4 fill-current ' + (star <= review.rating ? 'text-yellow-400' : 'text-gray-300')"
                     viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                </svg>
              </div>
              <span class="font-medium">{{ review.rating }}/5</span>
            </div>
            <span class="text-sm text-gray-500">{{ formatDate(review.created_at) }}</span>
          </div>
          <p class="text-gray-600 mb-2">{{ review.comment }}</p>
          <div *ngIf="review.images && review.images.length > 0" class="mt-3">
            <div class="grid grid-cols-2 gap-2 max-w-md">
              <img
                *ngFor="let imageUrl of review.images"
                [src]="imageUrl"
                [alt]="'Review image'"
                (error)="onImageError($event)"
                class="w-full h-20 object-cover rounded-md border hover:scale-105 transition-transform cursor-pointer">
            </div>
          </div>
          <div class="text-sm text-gray-500">
            Visited: {{ formatDate(review.visited_at) }}
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div *ngIf="!loading && !error && reviews.length === 0" class="text-center py-16">
        <svg class="mx-auto h-16 w-16 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012 2v8a2 2 0 01-2 2h-1l-4 4z" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900">No reviews yet</h3>
        <p class="mt-2 text-gray-500">Reviews will appear here once they are submitted.</p>
      </div>
    </div>
  `
})
export class ReviewsComponent implements OnInit {
  userRole: string = '';
  reviews: Review[] = [];
  loading = true;
  error = '';

  constructor(private reviewService: ReviewService) {}

  ngOnInit(): void {
    this.userRole = this.getUserRoleFromToken();
    this.loadReviews();
  }

  getUserRoleFromToken(): string {
    const token = localStorage.getItem('access_token');
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return payload.role || payload.user_type || 'user';
      } catch (error) {
        return '';
      }
    }
    return '';
  }

  loadReviews(): void {
    this.loading = true;
    let reviewObservable;

    switch (this.userRole) {
      case 'admin':
        reviewObservable = this.reviewService.getAllReviews();
        break;
      case 'guide':
        reviewObservable = this.reviewService.getGuideReviews();
        break;
      case 'tourist':
        reviewObservable = this.reviewService.getMyReviews();
        break;
      default:
        this.error = 'Unauthorized';
        this.loading = false;
        return;
    }

    reviewObservable.subscribe({
      next: (reviews) => {
        this.reviews = reviews;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading reviews:', err);
        this.error = 'Failed to load reviews';
        this.loading = false;
      }
    });
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  onImageError(event: any): void {
    event.target.style.display = 'none';
  }
}