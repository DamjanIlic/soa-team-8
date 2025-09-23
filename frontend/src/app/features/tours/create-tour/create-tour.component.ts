import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms'; // <--- OBAVEZNO
import { Router, RouterModule } from '@angular/router';
import { TourRequest } from '../../../core/models/tour.model';
import { TourService } from '../../../core/services/tour.service';

@Component({
  selector: 'app-create-tour',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule], 
  templateUrl: './create-tour.component.html'
})
export class CreateTourComponent {
  tour: TourRequest = {
    name: '',
    description: '',
    difficulty: 'EASY',
    tags: ''
  };

  message = '';

  constructor(private tourService: TourService, private router: Router) {}

  createTour() {
    if (!this.tour.name || !this.tour.description) {
      this.message = 'Name and Description are required!';
      return;
    }

    this.tourService.createTour(this.tour).subscribe({
  next: (res) => {
    this.router.navigate(['/tours/add-tour-checkpoints'], { state: { tour: res } });
  },
  error: (err) => (this.message = 'Failed to create tour. ' + (err.error || ''))
});
  }
}
