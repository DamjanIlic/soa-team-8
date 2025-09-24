import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { jwtDecode } from 'jwt-decode';
import { forkJoin } from 'rxjs';
import { Checkpoint, Duration } from '../../../core/models/tour.model';
import { TourService } from '../../../core/services/tour.service';

// Lokalni tipovi sa saved flagom (sada readonly)
interface TourCheckpoint extends Checkpoint {
  saved?: boolean;
}

interface TourDuration extends Duration {
  saved?: boolean;
}

interface Tour {
  id: string;
  author_id: string;
  name: string;
  description: string;
  difficulty: string;
  tags: string;
  status: 'draft' | 'published' | 'archived';
  price: number;
  distance_km?: number;
  checkpoints?: TourCheckpoint[];
  durations?: TourDuration[];
  updated_at?: string;
}

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
  statusLoading = false;

  private userId: string = '';

  constructor(private tourService: TourService) { }

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
    } catch (err) {
      console.error('Failed to decode JWT', err);
    }
  }

  loadAuthorTours(): void {
    this.loading = true;
    this.tourService.getAllTours().subscribe({
      next: (allTours) => {
        this.authorTours = allTours
          .filter(t => t.author_id === this.userId)
          .map(t => ({
            ...t,
            checkpoints: (t.checkpoints ?? []).map(kp => ({ ...kp })) as TourCheckpoint[],
            durations: (t.durations ?? []).map(d => ({ ...d })) as TourDuration[],
            price: t.price ?? 0
          }));
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading tours:', err);
        this.error = 'Failed to load your tours';
        this.loading = false;
      }
    });
  }

  editTourStatus(tour: Tour): void {
    this.selectedTour = { ...tour };
  }

  cancelEditStatus(): void {
    this.selectedTour = null;
  }

  saveStatus(newStatus?: 'published' | 'archived' | 'draft'): void {
    if (!this.selectedTour) return;

    const tour = this.selectedTour;
    const tourId = tour.id;

    // Reactivate archived to draft
    if (newStatus === 'draft' && tour.status === 'archived') {
      tour.status = 'draft';
      this.selectedTour = { ...tour };
      return;
    }

    // Validate price
    tour.price = Number(tour.price);
    if (isNaN(tour.price) || tour.price <= 0) {
      alert('Price must be greater than 0.');
      return;
    }

    // Validate required fields
    if (!tour.name || !tour.description || !tour.difficulty || !tour.tags) {
      alert('Fill in all required fields: name, description, difficulty, tags.');
      return;
    }

    if ((tour.checkpoints?.length ?? 0) < 2) {
      alert('Add at least 2 checkpoints.');
      return;
    }

    this.statusLoading = true;

    // Save price first
    this.tourService.updatePrice(tourId, tour.price).subscribe({
      next: () => {
        // Auto-calculate durations if empty
        if ((tour.durations?.length ?? 0) === 0 && tour.distance_km) {
          const speeds: Record<string, number> = { walk: 5, bike: 15, car: 60 };
          tour.durations = Object.entries(speeds).map(([transport, speed]) => ({
            transport,
            duration: Math.ceil(tour.distance_km! / speed)
          })) as TourDuration[];
        }

        // Only send durations that aren't saved yet (checkpoint-i se više ne dodaju)
        const durationRequests = (tour.durations ?? [])
          .filter(d => !d.saved)
          .map(d => this.tourService.addDuration(tourId, d));

        forkJoin([...durationRequests]).subscribe({
          next: () => {
            (tour.durations ?? []).forEach(d => d.saved = true);

            // Handle status update
            if (newStatus === 'published') {
              this.tourService.publishTour(tourId).subscribe({
                next: () => {
                  this.loadAuthorTours();
                  this.selectedTour = null;
                  this.statusLoading = false;
                },
                error: (err) => {
                  console.error('Failed to publish tour:', err);
                  alert(err.error || 'Publish failed');
                  this.statusLoading = false;
                }
              });
            } else if (newStatus === 'archived') {
              this.tourService.archiveTour(tourId).subscribe({
                next: () => {
                  this.loadAuthorTours();
                  this.selectedTour = null;
                  this.statusLoading = false;
                },
                error: (err) => {
                  console.error('Failed to archive tour:', err);
                  alert(err.error || 'Archive failed');
                  this.statusLoading = false;
                }
              });
            }
          },
          error: (err) => {
            console.error('Failed to save durations:', err);
            alert('Failed to save durations before publishing.');
            this.statusLoading = false;
          }
        });
      },
      error: (err) => {
        alert(err.error || 'Failed to set price');
        this.statusLoading = false;
      }
    });
  }

  formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }

  getStatusColor(status: string): string {
    switch (status.toLowerCase()) {
      case 'draft': return 'bg-yellow-100 text-yellow-800';
      case 'published': return 'bg-green-100 text-green-800';
      case 'archived': return 'bg-gray-100 text-gray-800';
      default: return 'bg-blue-100 text-blue-800';
    }
  }

  getDifficultyColor(diff: string): string {
    switch (diff.toLowerCase()) {
      case 'easy': return 'bg-green-100 text-green-800';
      case 'medium': return 'bg-yellow-100 text-yellow-800';
      case 'hard': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }

  parseTags(tags: string): string[] {
    return tags ? tags.split(',').map(t => t.trim()) : [];
  }

  trackByTourId(index: number, tour: Tour) {
    return tour.id;
  }
}
