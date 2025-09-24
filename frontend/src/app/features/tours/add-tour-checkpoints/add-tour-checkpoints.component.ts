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
        // Dodaj checkpoint u listu i update referencu za mapu
        this.checkpoints = [...this.checkpoints, savedCheckpoint];
        this.updateMapMarkers();
        this.calculateDurations();
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
      this.calculateDurations();
    }
  }

  updateMapMarkers(): void {
    if (this.mapComponent) {
      // Kreira novu referencu da Angular detektuje promenu
      this.mapComponent.addedCheckpointCollection = [...this.checkpoints];
    }
  }

  calculateDurations(): void {
    if (!this.tourDistanceKm) return;

    const speedMap: Record<TransportType, number> = {
      [TransportType.Walk]: 5,
      [TransportType.Bike]: 15,
      [TransportType.Car]: 60
    };

    const durationMinutes = Math.ceil((this.tourDistanceKm / speedMap[this.selectedTransport]) * 60);

    // Ispravka: koristi 'minutes' umesto 'duration' i 'transport_type' umesto 'transport'
    const durationData = { 
      transport: this.selectedTransport, 
      minutes: durationMinutes 
    };

    const existingIndex = this.durations.findIndex(d => d.transport === this.selectedTransport);
    if (existingIndex >= 0) {
      this.durations[existingIndex].minutes = durationMinutes;
    } else {
      this.durations.push(durationData);
    }

    // Update tour durations u backend-u
    if (this.tourId) {
      const durationObservables = this.durations.map(d => 
        this.tourService.addDuration(this.tourId, {
          transport: d.transport,
          minutes: d.minutes
        })
      );
      
      // Ispravka: koristi forkJoin umesto deprecated toPromise()
      forkJoin(durationObservables).subscribe({
        next: () => console.log('Durations updated', this.durations),
        error: (err) => console.error('Failed to update durations', err)
      });
    }
  }

  onTransportChange(transport: TransportType) {
    this.selectedTransport = transport;
    this.calculateDurations();
  }

  onLocationSelected(location: { lat: number; lng: number }) {
    this.checkpointForm.get('latitude')?.setValue(location.lat.toString());
    this.checkpointForm.get('longitude')?.setValue(location.lng.toString());
  }

  resetForm(): void { this.checkpointForm.reset(); }

  cancelTour(): void { this.router.navigate(['/addNewTour']); }

  onRouteDistanceUpdated(distanceKm: number) {
    this.tourDistanceKm = distanceKm;
    this.calculateDurations();
  }

  toggleHelpModal() { this.isHelpModalOpen = !this.isHelpModalOpen; }

  finalizeTour(): void {
    if (!this.tour || this.checkpoints.length < 2) return;

    this.tourService.updateDistance(this.tourId, this.tourDistanceKm).subscribe({
      next: () => this.router.navigate(['/tours/author-tours']),
      error: (err) => console.error('Failed to finalize tour:', err)
    });
  }
}