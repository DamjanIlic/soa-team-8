import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourService, Tour } from '../../../core/services/tour.service';

@Component({
  selector: 'app-author-tours',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './author-tours.component.html',
  styleUrls: ['./author-tours.component.css']
})
export class AuthorToursComponent implements OnInit {
  authorTours: Tour[] = [];
  loading = true;
  error = '';

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    this.loadAuthorTours();
  }

  loadAuthorTours(): void {
    this.loading = true;
    console.log('Calling getAuthorTours API...');
    
    this.tourService.getAuthorTours()
      .subscribe({
        next: (tours) => {
          console.log('Received tours:', tours);
          this.authorTours = tours;
          this.loading = false;
        },
        error: (err) => {
          console.error('Error loading author tours:', err);
          this.error = 'Failed to load your tours';
          this.loading = false;
        }
      });
  }

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
}