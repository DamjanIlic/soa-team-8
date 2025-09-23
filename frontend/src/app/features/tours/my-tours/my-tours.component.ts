import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Tour } from '../../../core/models/tour.model';
import { PurchaseService } from '../../../core/services/purchase.service';
import { TourService } from '../../../core/services/tour.service';
import { ReviewModalComponent } from '../../../shared/components/review-modal/review-modal.component';

@Component({
  selector: 'app-my-tours',
  standalone: true,
  imports: [CommonModule, ReviewModalComponent],
  templateUrl: './my-tours.component.html',
  styleUrls: ['./my-tours.component.css']
})
export class MyToursComponent implements OnInit {
  purchasedTours: Tour[] = [];
  loading = true;
  error = '';
  showReviewModal = false;
  selectedTour: Tour | null = null;

  constructor(
    private tourService: TourService,
    private purchaseService: PurchaseService
  ) {}

  ngOnInit(): void {
    this.loadPurchasedTours();
  }

  loadPurchasedTours(): void {
    this.loading = true;
    this.purchaseService.getPurchasedTokens()
      .subscribe({
        next: (tokens) => {
          if (tokens.length === 0) {
            this.purchasedTours = [];
            this.loading = false;
            return;
          }

          // Get tour details for each purchased token
          const tourIds = tokens.map(token => token.tour_id);
          this.tourService.getAllTours()
            .subscribe({
              next: (allTours) => {
                this.purchasedTours = allTours.filter(tour => 
                  tourIds.includes(tour.id)
                );
                this.loading = false;
              },
              error: (err) => {
                console.error('Error loading tour details:', err);
                this.error = 'Failed to load tour details';
                this.loading = false;
              }
            });
        },
        error: (err) => {
          console.error('Error loading purchased tours:', err);
          this.error = 'Failed to load purchased tours';
          this.loading = false;
        }
      });
  }

  leaveReview(tour: Tour): void {
    this.selectedTour = tour;
    this.showReviewModal = true;
  }

  closeReviewModal(): void {
    this.showReviewModal = false;
    this.selectedTour = null;
  }

  onReviewSubmitted(): void {
    // Reload tours to update is_reviewed status
    this.loadPurchasedTours();
  }

  formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }

  getDifficultyColor(difficulty: string): string {
    switch (difficulty.toLowerCase()) {
      case 'easy': return 'bg-green-100 text-green-800';
      case 'medium': return 'bg-yellow-100 text-yellow-800';
      case 'hard': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }

  parseTags(tags: string): string[] {
    return tags ? tags.split(',').map(tag => tag.trim()) : [];
  }
}