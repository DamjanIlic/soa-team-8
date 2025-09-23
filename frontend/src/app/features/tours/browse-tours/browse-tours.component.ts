import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourService } from '../../../core/services/tour.service';
import { Tour } from '../../../core/services/tour.service';

@Component({
  selector: 'app-browse-tours',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './browse-tours.component.html',
  styleUrls: ['./browse-tours.component.css']
})
export class BrowseToursComponent implements OnInit {
  tours: Tour[] = [];
  loading = true;
  error = '';

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    this.loadTours();
  }

  loadTours(): void {
    this.loading = true;
    this.tourService.getAllTours()
      .subscribe({
        next: (data) => {
          this.tours = data;
          this.loading = false;
        },
        error: (err) => {
          console.error('Error loading tours:', err);
          this.error = 'Failed to load tours';
          this.loading = false;
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