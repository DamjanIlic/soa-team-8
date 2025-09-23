import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourExecutionService } from '../../../core/services/tourExecution.service'; 
import { ActivatedRoute } from '@angular/router';
import { Subscription, interval } from 'rxjs';
import * as L from 'leaflet';

@Component({
  selector: 'app-tour-execution',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './tour-execution.component.html',
  styleUrls: ['./tour-execution.component.css']
})
export class TourExecutionComponent implements OnInit, OnDestroy {

  tourExecutionId: string | null = null;
  status: 'active' | 'completed' | 'abandoned' = 'active';
  currentLatitude: number = 44.8176;
  currentLongitude: number = 20.4569;
  keyPoints: any[] = [];

  private intervalSub?: Subscription;
  private tourId: string = '';

  private map!: L.Map;
  private userMarker!: L.Marker;

  constructor(
    private tourExecutionService: TourExecutionService,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.tourId = this.route.snapshot.paramMap.get('id') || '';
    this.tourExecutionService.getKeyPoints(this.tourId).subscribe(kps => {
      // Dodaj reached i marker polja
      this.keyPoints = kps.map(kp => ({ ...kp, reached: false, marker: null }));
    });

    if (this.tourId) {
      this.tourExecutionService.startTour(this.tourId).subscribe(exec => {
        this.tourExecutionId = exec.id;
        this.status = exec.status;

        this.startPositionInterval();
        this.initMap();
      });
    }
  }

  private initMap(): void {
    const initialPos: L.LatLngExpression = [this.currentLatitude, this.currentLongitude];

    this.map = L.map('positionMap').setView(initialPos, 15);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    // Marker korisnika
    this.userMarker = L.marker(initialPos, { title: 'User' })
      .addTo(this.map)
      .bindPopup('User is here')
      .openPopup();

    // Keypoints kao circleMarker
    this.keyPoints.forEach(kp => {
      kp.marker = L.circleMarker([kp.latitude, kp.longitude], {
        radius: 8,
        color: kp.reached ? 'green' : 'red',
        fillColor: kp.reached ? 'green' : 'red',
        fillOpacity: 0.8
      }).addTo(this.map)
        .bindPopup(`${kp.name} (${kp.latitude.toFixed(5)}, ${kp.longitude.toFixed(5)})`);
    });

    // Klik na mapu menja trenutnu poziciju korisnika
    this.map.on('click', (e: L.LeafletMouseEvent) => this.onMapClick(e));
  }

  private onMapClick(event: L.LeafletMouseEvent) {
    const latlng = event.latlng;
    this.userMarker.setLatLng(latlng);
    this.currentLatitude = latlng.lat;
    this.currentLongitude = latlng.lng;
    this.userMarker.getPopup()?.setContent(`User is here: ${latlng.lat.toFixed(6)}, ${latlng.lng.toFixed(6)}`);
    this.userMarker.openPopup();
  }

  startPositionInterval() {
    if (!this.tourExecutionId) return;

    this.intervalSub = interval(10000).subscribe(() => {
      if (!this.tourExecutionId || !this.userMarker) return;

      // Uzmi trenutnu poziciju korisnika sa mape
      const currentPos = this.userMarker.getLatLng();
      this.currentLatitude = currentPos.lat;
      this.currentLongitude = currentPos.lng;

      // Pozovi backend samo za proveru keypoints
      this.keyPoints.forEach(kp => {
        if (!kp.reached) {
          this.tourExecutionService.checkKeyPoint(
            this.tourExecutionId!,
            kp.id,
            this.currentLatitude,
            this.currentLongitude
          ).subscribe(res => {
            if (res.reached) {
              kp.reached = true;
              kp.reachedAt = new Date();
              if (kp.marker) {
                kp.marker.setStyle({ color: 'green', fillColor: 'green' });
              }
            }
          });
        }
      });
    });
  }

  ngOnDestroy(): void {
    this.intervalSub?.unsubscribe();
    this.map?.remove();
  }

  completeTour(): void {
    if (!this.tourExecutionId) return;

    const allReached = this.keyPoints.every(kp => kp.reached);
    if (!allReached) {
      alert("You haven't reached all key points yet!");
      return;
    }

    this.tourExecutionService.completeTour(this.tourExecutionId)
      .subscribe({
        next: () => {
          this.status = 'completed';
          alert("Tour completed successfully!");
          this.intervalSub?.unsubscribe();
        },
        error: (err) => {
          console.error("Error completing tour:", err);
          alert("Failed to complete the tour.");
        }
      });
  }

  abandonTour(): void {
    if (!this.tourExecutionId) return;

    this.tourExecutionService.abandonTour(this.tourExecutionId)
      .subscribe({
        next: () => {
          this.status = 'abandoned';
          alert("Tour has been abandoned.");
          this.intervalSub?.unsubscribe();
          this.intervalSub = undefined;
        },
        error: (err) => {
          console.error("Error abandoning tour:", err);
          alert("Failed to abandon the tour.");
        }
    });
  }
}
