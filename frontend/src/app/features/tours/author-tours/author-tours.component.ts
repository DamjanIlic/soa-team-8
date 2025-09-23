import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { jwtDecode } from 'jwt-decode';
import { Tour } from '../../../core/models/tour.model';
import { TourService } from '../../../core/services/tour.service';

@Component({
  selector: 'app-author-tours',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './author-tours.component.html',
  styleUrls: ['./author-tours.component.css']
})
export class AuthorToursComponent implements OnInit {
  authorTours: Tour[] = [];
  loading = true;
  error = '';

  selectedTour: Tour | null = null;
  statusOptions: string[] = ['draft', 'published', 'archived'];
  statusLoading = false;

  private userId: string = '';

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    this.extractUserIdFromToken();
    this.loadAuthorTours();
  }

  private extractUserIdFromToken(): void {
    const token = localStorage.getItem('access_token');
    if (!token) return;

    try {
      const decoded: any = jwtDecode(token);
      this.userId = decoded.sub || decoded.user_id || '';
      console.log('Logged-in user ID from JWT:', this.userId);
    } catch (err) {
      console.error('Failed to decode JWT', err);
    }
  }

  loadAuthorTours(): void {
    this.loading = true;
    this.tourService.getAllTours().subscribe({
      next: (allTours) => {
        this.authorTours = allTours.filter(t => t.author_id === this.userId);
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading tours:', err);
        this.error = 'Failed to load your tours';
        this.loading = false;
      }
    });
  }

  // --- Status edit functions ---
  editTourStatus(tour: Tour): void {
    this.selectedTour = { ...tour }; // kopija da ne menjaš odmah UI
  }

  cancelEditStatus(): void {
    this.selectedTour = null;
  }

  saveStatus(): void {
    if (!this.selectedTour) return;

    this.statusLoading = true;
    const tourId = this.selectedTour.id;
    const newStatus = this.selectedTour.status.toLowerCase();

    let update$;
    switch (newStatus) {
      case 'published':
        update$ = this.tourService.publishTour(tourId);
        break;
      case 'archived':
        update$ = this.tourService.archiveTour(tourId);
        break;
      case 'draft':
        update$ = this.tourService.reactivateTour(tourId);
        break;
      default:
        console.error('Unknown status:', newStatus);
        this.statusLoading = false;
        return;
    }

    update$.subscribe({
      next: (updatedTour) => {
        const index = this.authorTours.findIndex(t => t.id === updatedTour.id);
        if (index > -1) this.authorTours[index] = updatedTour;
        this.selectedTour = null;
        this.statusLoading = false;
      },
      error: (err) => {
        console.error('Failed to update status:', err);
        alert('Failed to update status. Make sure tour meets the requirements.');
        this.statusLoading = false;
      }
    });
  }

  // --- Helpers ---
  formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }

  getStatusColor(status: string): string {
    switch (status.toLowerCase()) {
      case 'published': return 'bg-green-100 text-green-800';
      case 'draft': return 'bg-yellow-100 text-yellow-800';
      case 'archived': return 'bg-gray-100 text-gray-800';
      default: return 'bg-blue-100 text-blue-800';
    }
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

  trackByTourId(index: number, tour: Tour) {
    return tour.id;
  }
}
