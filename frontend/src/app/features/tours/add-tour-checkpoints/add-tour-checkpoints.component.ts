import { CommonModule } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { Checkpoint } from '../../../core/models/checkpoint.model';
import { Tour } from '../../../core/models/tour.model';
import { TourService } from '../../../core/services/tour.service';
import { MapComponent } from '../../../shared/map/map.component';

@Component({
  selector: 'xp-add-tour-checkpoints',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule, MapComponent],
  templateUrl: './add-tour-checkpoints.component.html'
})
export class AddTourCheckpointsComponent implements OnInit {
  tour: Tour | null = null;
  tourId!: string;
  checkpoints: Checkpoint[] = [];
  isHelpModalOpen = false;
  saving = false;
  tourDistanceKm = 0;

  @ViewChild('map', { static: false }) mapComponent!: MapComponent;

  checkpointForm = new FormGroup({
    name: new FormControl('', Validators.required),
    description: new FormControl('', Validators.required),
    latitude: new FormControl('', Validators.required),
    longitude: new FormControl('', Validators.required),
    image_url: new FormControl(''),
  });

  constructor(private tourService: TourService, private router: Router) { }

  ngOnInit(): void {
    const state = history.state as { tour: Tour };
    if (state?.tour) {
      this.tour = state.tour;
      this.tourId = state.tour.id!;
      console.log('Loaded tour ID:', this.tourId);
    } else {
      console.error('No tour data found');
      this.router.navigate(['/addNewTour']);
    }
  }

  addCheckpoint(): void {
    console.log('Add checkpoint clicked');

    if (this.checkpointForm.invalid || !this.tourId) {
      console.warn('Form invalid or tourId missing');
      return;
    }

    const checkpoint: Checkpoint = {
      name: this.checkpointForm.value.name!.trim(),
      description: this.checkpointForm.value.description!.trim(),
      latitude: Number(this.checkpointForm.value.latitude),
      longitude: Number(this.checkpointForm.value.longitude),
      image_url: this.checkpointForm.value.image_url?.trim() || undefined
    };

    this.saving = true;

    // odmah šalje u backend
    this.tourService.addKeyPoint(this.tourId, checkpoint).subscribe({
      next: (savedCheckpoint) => {
        console.log('Checkpoint successfully added to DB:', savedCheckpoint);
        this.checkpoints = [...this.checkpoints, savedCheckpoint];

        this.resetForm();
        this.saving = false;
      },
      error: (err) => {
        console.error('Failed to add checkpoint:', err);
        this.saving = false;
      }
    });
  }

  onLocationSelected(location: { lat: number; lng: number }) {
    this.checkpointForm.get('latitude')?.setValue(location.lat.toString());
    this.checkpointForm.get('longitude')?.setValue(location.lng.toString());
  }

  resetForm(): void {
    this.checkpointForm.reset();
  }

  cancelTour(): void {
    this.router.navigate(['/addNewTour']);
  }

  onRouteDistanceUpdated(distanceKm: number) {
    this.tourDistanceKm = distanceKm;
    console.log('Route distance updated:', distanceKm, 'km');
  }

  handleCheckpointRemoved(index: number): void {
    if (index >= 0 && index < this.checkpoints.length) {
      this.checkpoints.splice(index, 1);
      console.log('Checkpoint removed:', index);
    }
  }

  toggleHelpModal() {
    this.isHelpModalOpen = !this.isHelpModalOpen;
  }

  finalizeTour(): void {
    if (!this.tour || this.checkpoints.length < 2) {
      console.error('Tour missing or less than 2 checkpoints');
      return;
    }

    console.log('Finalizing tour with checkpoints:', this.checkpoints);

    this.tourService.updateDistance(this.tourId, this.tourDistanceKm).subscribe({
      next: (updatedTour) => {
        console.log('Tour finalized with distance:', updatedTour);
        this.router.navigate(['/tours/author-tours']);
      },
      error: (err) => {
        console.error('Failed to finalize tour:', err);
      }
    });
  }

}
