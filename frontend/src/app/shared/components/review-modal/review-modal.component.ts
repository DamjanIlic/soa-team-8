import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ReviewService, ReviewRequest } from '../../../core/services/review.service';
import { Tour } from '../../../core/models/tour.model';

@Component({
  selector: 'app-review-modal',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div *ngIf="isVisible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-bold">Leave a Review</h2>
          <button (click)="closeModal()" class="text-gray-500 hover:text-gray-700">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <div class="mb-4">
          <h3 class="font-medium text-gray-800">{{ tour?.name }}</h3>
        </div>

        <form (ngSubmit)="submitReview()" #reviewForm="ngForm">
          <!-- Rating -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Rating</label>
            <div class="flex space-x-1">
              <button
                type="button"
                *ngFor="let star of [1, 2, 3, 4, 5]"
                (click)="setRating(star)"
                [class]="'p-1 ' + (star <= rating ? 'text-yellow-400' : 'text-gray-300')"
              >
                <svg class="w-6 h-6 fill-current" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Visit Date -->
          <div class="mb-4">
            <label for="visitDate" class="block text-sm font-medium text-gray-700 mb-2">Visit Date</label>
            <input
              type="date"
              id="visitDate"
              [(ngModel)]="visitDate"
              name="visitDate"
              required
              [max]="maxDate"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- Comment -->
          <div class="mb-4">
            <label for="comment" class="block text-sm font-medium text-gray-700 mb-2">Comment</label>
            <textarea
              id="comment"
              [(ngModel)]="comment"
              name="comment"
              rows="4"
              placeholder="Share your experience..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            ></textarea>
          </div>

          <!-- Images -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Images (URLs)</label>
            <div class="space-y-2">
              <input
                *ngFor="let url of imageUrls; let i = index; trackBy: trackByIndex"
                type="url"
                [(ngModel)]="imageUrls[i]"
                [name]="'imageUrl' + i"
                placeholder="https://example.com/image.jpg"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div class="flex space-x-2 mt-2">
              <button
                type="button"
                (click)="addImageUrl()"
                class="px-3 py-1 text-sm bg-gray-200 hover:bg-gray-300 rounded-md">
                + Add Image URL
              </button>
              <button
                *ngIf="imageUrls.length > 1"
                type="button"
                (click)="removeLastImageUrl()"
                class="px-3 py-1 text-sm bg-red-200 hover:bg-red-300 rounded-md">
                Remove Last
              </button>
            </div>
          </div>

          <!-- Submit Buttons -->
          <div class="flex justify-end space-x-2">
            <button
              type="button"
              (click)="closeModal()"
              class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              [disabled]="rating === 0 || !visitDate || loading"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {{ loading ? 'Submitting...' : 'Submit Review' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  `
})
export class ReviewModalComponent {
  @Input() isVisible = false;
  @Input() tour: Tour | null = null;
  @Output() close = new EventEmitter<void>();
  @Output() reviewSubmitted = new EventEmitter<void>();

  closeModal(): void {
    this.close.emit();
  }

  rating = 0;
  comment = '';
  visitDate = '';
  imageUrls: string[] = [''];
  loading = false;

  constructor(private reviewService: ReviewService) {}

  setRating(star: number): void {
    this.rating = star;
  }

  get maxDate(): string {
    return new Date().toISOString().split('T')[0];
  }

  submitReview(): void {
    if (!this.tour || this.rating === 0 || !this.visitDate) {
      return;
    }

    // validacija da datum nije u buducnosti
    const selectedDate = new Date(this.visitDate);
    const today = new Date();
    today.setHours(23, 59, 59, 999); // ukljucuje ceo danasnji dan
    
    if (selectedDate > today) {
      alert('Visit date cannot be in the future');
      return;
    }

    this.loading = true;

    //filter out empty urls
    const validImageUrls = this.imageUrls.filter(url => url.trim() !== '');

    const reviewRequest: ReviewRequest = {
      rating: this.rating,
      comment: this.comment,
      visited_at: this.visitDate,
      images: validImageUrls
    };

    this.reviewService.createReview(this.tour.id, reviewRequest).subscribe({
      next: () => {
        this.reviewSubmitted.emit();
        this.resetForm();
        this.close.emit();
        this.loading = false;
      },
      error: (err) => {
        console.error('Error submitting review:', err);
        alert('Failed to submit review. Please try again.');
        this.loading = false;
      }
    });
  }

  private resetForm(): void {
    this.rating = 0;
    this.comment = '';
    this.visitDate = '';
  }

  addImageUrl(): void {
    this.imageUrls.push('');
  }

  removeLastImageUrl(): void {
    if (this.imageUrls.length > 1) {
      this.imageUrls.pop();
    }
  }

  trackByIndex(index: number): number {
    return index;
  }

}