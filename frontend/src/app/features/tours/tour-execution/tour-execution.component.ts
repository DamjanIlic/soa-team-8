import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourExecutionService } from '../../../core/services/tourExecution.service'; 
import { ActivatedRoute } from '@angular/router';
import { Subscription, interval } from 'rxjs';

@Component({
  selector: 'app-tour-execution',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './tour-execution.component.html',
  styleUrls: ['./tour-execution.component.css']
})
export class TourExecutionComponent implements OnInit, OnDestroy {

  // TourExecution ID dobijen od backa nakon starta
  tourExecutionId: string | null = null;

  // Status ture
  status: 'active' | 'completed' | 'abandoned' = 'active';

  // Trenutna pozicija
  currentLatitude: number = 44.8176;
  currentLongitude: number = 20.4569;

  // Lista keypoints - trenutno placeholder
  keyPoints: any[] = [];

  private intervalSub?: Subscription;
  private tourId: string = '';

  constructor(
    private tourExecutionService: TourExecutionService,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    // Dohvati tourId iz rute
    this.tourId = this.route.snapshot.paramMap.get('id') || '';
    this.tourExecutionService.getKeyPoints(this.tourId).subscribe(kps => {
    this.keyPoints = kps.map(kp => ({ ...kp, reached: false }));
    });
    // Pokreni turu
    if (this.tourId) {
      this.tourExecutionService.startTour(this.tourId).subscribe(exec => {
        this.tourExecutionId = exec.id;
        this.status = exec.status;

        // Startuj interval za update pozicije na svakih 10s
        this.startPositionInterval();
      });
    }
  }

  startPositionInterval() {
    if (!this.tourExecutionId) return;

    this.intervalSub = interval(10000).subscribe(() => {
      if (!this.tourExecutionId) return;

      // 1. Update trenutne pozicije
      this.tourExecutionService.getSimulatedPosition(this.tourExecutionId)
        .subscribe(pos => {
          this.currentLatitude = pos.latitude;
          this.currentLongitude = pos.longitude;

          // 2. Proveri sve keypoints
          this.keyPoints.forEach(kp => {
            if (!kp.reached) { // samo za one koji još nisu dostignuti
              this.tourExecutionService.checkKeyPoint(
                this.tourExecutionId!,
                kp.id,
                this.currentLatitude,
                this.currentLongitude
              ).subscribe(res => {
                if (res.reached) {
                  kp.reached = true;
                  kp.reachedAt = new Date();
                }
              });
            }
          });
        });
    });
  }

  ngOnDestroy(): void {
    // Očisti interval kada se komponenta uništi
    this.intervalSub?.unsubscribe();
  }
  completeTour(): void {
    if (!this.tourExecutionId) return;

    // Proveri da li su svi keypoints reached
    const allReached = this.keyPoints.every(kp => kp.reached);

    if (!allReached) {
      alert("You haven't reached all key points yet!");
      return;
    }

    // Ako jesu svi reached, pozovi backend
    this.tourExecutionService.completeTour(this.tourExecutionId)
      .subscribe({
        next: () => {
          this.status = 'completed';
          alert("Tour completed successfully!");

          // Zaustavi interval simulacije pozicije
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

          // Stopiraj simulator pozicije
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
