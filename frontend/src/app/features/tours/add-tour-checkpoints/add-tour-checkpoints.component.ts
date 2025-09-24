import { CommonModule } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { forkJoin } from 'rxjs';
import { Checkpoint } from '../../../core/models/checkpoint.model';
import { Duration, Tour } from '../../../core/models/tour.model';
import { TourService } from '../../../core/services/tour.service';
import { MapComponent } from '../../../shared/map/map.component';

export enum TransportType {
  Walk = 'walk',
  Bike = 'bike',
  Car = 'car'
}

@Component({
  selector: 'xp-add-tour-checkpoints',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, RouterModule, MapComponent],
  templateUrl: './add-tour-checkpoints.component.html'
})
export class AddTourCheckpointsComponent implements OnInit {
  selectedTransport: TransportType = TransportType.Walk;
  TransportType = TransportType;
  minutes = 0;

  tour: Tour | null = null;
  tourId!: string;
  checkpoints: Checkpoint[] = [];
  durations: Duration[] = [];
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
    } else {
      console.error('No tour data found');
      this.router.navigate(['/addNewTour']);
    }
  }

  addCheckpoint(): void {
    if (this.checkpointForm.invalid || !this.tourId) return;

    const checkpoint: Checkpoint = {
      name: this.checkpointForm.value.name!.trim(),
      description: this.checkpointForm.value.description!.trim(),
      latitude: Number(this.checkpointForm.value.latitude),
      longitude: Number(this.checkpointForm.value.longitude),
      image_url: this.checkpointForm.value.image_url?.trim() || undefined
    };

    this.saving = true;
    this.tourService.addKeyPoint(this.tourId, checkpoint).subscribe({
      next: (savedCheckpoint) => {
        this.checkpoints = [...this.checkpoints, savedCheckpoint];
        this.updateMapMarkers();
        this.resetForm();
        this.saving = false;
      },
      error: (err) => {
        console.error('Failed to add checkpoint:', err);
        this.saving = false;
      }
    });
  }

  handleCheckpointRemoved(index: number): void {
    if (index >= 0 && index < this.checkpoints.length) {
      this.checkpoints.splice(index, 1);
      this.updateMapMarkers();
    }
  }

  updateMapMarkers(): void {
    if (this.mapComponent) {
      this.mapComponent.addedCheckpointCollection = [...this.checkpoints];
    }
  }

  onTransportChange(transport: TransportType) {
    this.selectedTransport = transport;
    this.updateMinutes();
  }

  updateMinutes() {
    const speedMap: Record<TransportType, number> = {
      [TransportType.Walk]: 5,
      [TransportType.Bike]: 15,
      [TransportType.Car]: 60
    };
    const speed = speedMap[this.selectedTransport] ?? 1;
    this.minutes = Math.ceil((this.tourDistanceKm / speed) * 60);
  }

  addDuration() {
    if (this.durations.some(d => d.transport === this.selectedTransport)) return;
    this.durations.push({ transport: this.selectedTransport, minutes: this.minutes });
    this.syncDurationsToBackend();
  }

  syncDurationsToBackend() {
    if (!this.tourId) return;
    const observables = this.durations.map(d => this.tourService.addDuration(this.tourId, d));
    forkJoin(observables).subscribe({
      next: () => console.log('Durations updated', this.durations),
      error: (err) => console.error('Failed to update durations', err)
    });
  }

  onLocationSelected(location: { lat: number; lng: number }) {
    this.checkpointForm.get('latitude')?.setValue(location.lat.toString());
    this.checkpointForm.get('longitude')?.setValue(location.lng.toString());
  }

  resetForm(): void { this.checkpointForm.reset(); }

  cancelTour(): void { this.router.navigate(['/addNewTour']); }

  onRouteDistanceUpdated(distanceKm: number) {
    this.tourDistanceKm = distanceKm;
    this.updateMinutes();
  }

  toggleHelpModal() { this.isHelpModalOpen = !this.isHelpModalOpen; }

  finalizeTour(): void {
    if (!this.tour || this.checkpoints.length < 2) {
      alert('Tour must have at least 2 checkpoints before finalizing.');
      return;
    }

    // Update distance first
    this.tourService.updateDistance(this.tourId, this.tourDistanceKm).subscribe({
      next: () => {
        // Optionally navigate or show success message
        this.router.navigate(['/tours/author-tours']);
      },
      error: (err) => console.error('Failed to finalize tour:', err)
    });
  }
}
