import { CommonModule } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { Checkpoint } from '../../../core/models/checkpoint.model';
import { Tour } from '../../../core/models/tour.model';
import { TourService } from '../../../core/services/tour.service';

import { MapComponent } from '../../../shared/map/map.component';

@Component({
  selector: 'xp-add-tour-checkpoints',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    MapComponent
  ],
  templateUrl: './add-tour-checkpoints.component.html'
})
export class AddTourCheckpointsComponent implements OnInit {
  tour: Tour | null = null;
  tourId!: string;
  checkpoints: Checkpoint[] = [];
  checkpointCollection: Checkpoint[] = [];
  isHelpModalOpen = false;

  @ViewChild('map', { static: false }) mapComponent!: MapComponent;

  checkpointForm = new FormGroup({
    name: new FormControl('', Validators.required),
    description: new FormControl('', Validators.required),
    latitude: new FormControl('', Validators.required),
    longitude: new FormControl('', Validators.required),
    image_url: new FormControl(''), // URL slike
  });

  constructor(
    private tourService: TourService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    const state = history.state as { tour: Tour };
    if (state?.tour) {
      this.tour = state.tour;
      this.tourId = state.tour.id!;
    } else {
      console.error('No tour data found');
      this.router.navigate(['/addNewTour']);
    }
  }

  addCheckpoint(): void {
    // Kreiraj payload koji backend očekuje
    const checkpoint: Checkpoint = {
      name: this.checkpointForm.value.name!,
      description: this.checkpointForm.value.description!,
      latitude: Number(this.checkpointForm.value.latitude),
      longitude: Number(this.checkpointForm.value.longitude),
      image_url: this.checkpointForm.value.image_url || undefined
    };

    this.checkpoints.push(checkpoint);
    this.checkpointCollection = [...this.checkpoints];
    this.resetForm();
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

  finalizeTour(): void {
    if (!this.tour || this.checkpoints.length < 2) {
      console.error('Tour missing or less than 2 checkpoints');
      return;
    }

    // Pošalji payload sa URL-om slike
    this.tourService.addKeyPoints(this.tour.id!, this.checkpoints).subscribe({
      next: (res) => {
        console.log('Checkpoints added:', res);
        this.router.navigate(['/mytours']);
      },
      error: (err) => console.error('Error adding checkpoints:', err)
    });
  }

  toggleHelpModal() {
    this.isHelpModalOpen = !this.isHelpModalOpen;
  }

  handleCheckpointRemoved(index: number): void {
    if (index >= 0 && index < this.checkpoints.length) {
      this.checkpoints.splice(index, 1);
      this.checkpointCollection = [...this.checkpoints];
    }
  }
}
