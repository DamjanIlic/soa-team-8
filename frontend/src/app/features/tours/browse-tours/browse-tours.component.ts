import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Tour } from '../../../core/models/tour.model';
import { CartService } from '../../../core/services/cart.service';
import { PurchaseService } from '../../../core/services/purchase.service';
import { TourService } from '../../../core/services/tour.service';

@Component({
  selector: 'app-browse-tours',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './browse-tours.component.html',
  styleUrls: ['./browse-tours.component.css']
})
export class BrowseToursComponent implements OnInit {
  tours: Tour[] = [];
  publishedTours: Tour[] = [];
  purchasedTourIds: string[] = [];
  loading = true;
  error = '';

  constructor(
    private tourService: TourService,
    private cartService: CartService,
    private purchaseService: PurchaseService
  ) {}

  ngOnInit(): void {
    this.loadTours();
    this.loadPurchasedTours();
  }

  loadTours(): void {
    this.loading = true;
    this.tourService.getAllTours()
      .subscribe({
        next: (data) => {
          this.tours = data;
          this.filterAvailableTours();
          this.loading = false;
        },
        error: (err) => {
          console.error('Error loading tours:', err);
          this.error = 'Failed to load tours';
          this.loading = false;
        }
      });
  }

    loadPurchasedTours(): void {
    this.purchaseService.getPurchasedTokens()
      .subscribe({
        next: (tokens) => {
          this.purchasedTourIds = tokens.map(token => token.tour_id);
          this.filterAvailableTours();
        },
        error: (err) => {
          console.error('Error loading purchased tours:', err);
          // ako nema kupljenih tura ili greske, dalje
        }
      });
  }

  filterAvailableTours(): void {
    this.publishedTours = this.tours.filter(tour => 
      tour.status.toLowerCase() === 'published' &&
      !this.purchasedTourIds.includes(tour.id)
    );
  }

  addToCart(tour: Tour): void {
    this.cartService.addItem(tour.id, tour.name, tour.price)
      .subscribe({
        next: () => {
          this.cartService.refreshCart();
          this.cartService.showCart();
        },
        error: (err) => {
          if (err.error && err.error.includes('already in cart')) {
            alert('Tour is already in your cart!');
          } else {
            console.error('Error adding to cart:', err);
            alert('Failed to add tour to cart. Please try again.');
          }
        }
      });
  }

  getDifficultyColor(difficulty: string): string {
    switch (difficulty.toLowerCase()) {
      case 'easy': return 'bg-green-100 text-green-800';
      case 'medium': return 'bg-yellow-100 text-yellow-800';
      case 'hard': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }

  getStatusColor(status: string): string {
    switch (status.toLowerCase()) {
      case 'published': return 'bg-green-100 text-green-800';
      case 'draft': return 'bg-gray-100 text-gray-800';
      default: return 'bg-blue-100 text-blue-800';
    }
  }

  formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }

  parseTags(tags: string): string[] {
    return tags ? tags.split(',').map(tag => tag.trim()) : [];
  }
}