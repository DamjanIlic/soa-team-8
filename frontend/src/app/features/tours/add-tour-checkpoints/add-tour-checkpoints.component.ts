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
  checkpointCollection: any[] = [];
  selectedImage: File | null = null;
  imagePreview: string | ArrayBuffer | null = null;
  clearMarkersFlag = false;
  isHelpModalOpen = false;

  @ViewChild('map', { static: false }) mapComponent!: MapComponent;

  checkpointForm = new FormGroup({
    name: new FormControl('', Validators.required),
    description: new FormControl('', Validators.required),
    secret: new FormControl('', Validators.required),
    latitude: new FormControl('', Validators.required),
    longitude: new FormControl('', Validators.required),
    image: new FormControl(''),
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

  onFileSelect(event: any): void {
    const file = event.target.files[0];
    if (file) {
      this.selectedImage = file;
      const reader = new FileReader();
      reader.onload = () => this.imagePreview = reader.result;
      reader.readAsDataURL(file);
    }
  }

  addCheckpoint(): void {
    const checkpoint: Checkpoint = {
      name: this.checkpointForm.value.name!,
      description: this.checkpointForm.value.description!,
      secret: this.checkpointForm.value.secret!,
      latitude: Number(this.checkpointForm.value.latitude),
      longitude: Number(this.checkpointForm.value.longitude),
      image: this.selectedImage ? {
        data: (this.imagePreview as string).split(',')[1],
        mimeType: this.selectedImage.type,
        uploadedAt: new Date().toISOString()
      } : undefined
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
    this.selectedImage = null;
    this.imagePreview = null;
  }

  cancelTour(): void {
    this.router.navigate(['/addNewTour']);
  }

  finalizeTour(): void {
    if (!this.tour || this.checkpoints.length < 2) {
      console.error('Tour missing or less than 2 checkpoints');
      return;
    }

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
