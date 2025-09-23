import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface TourExecution {
    id: string;
    tourId: string;
    userId: string;
    status: 'active' | 'completed' | 'abandoned';
    startedAt?: string;
    endedAt?: string;
    lastActive?: string;
}
export interface KeyPoint {
  id: string;
  tour_id: string;
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  order: number;
  created_at: string;
}
@Injectable({
  providedIn: 'root'
})
export class TourExecutionService {

    private apiUrl = 'http://localhost:8000/api/tours/executions';
    private toursUrl = 'http://localhost:8000/api/tours';
    constructor(private http: HttpClient) {}

    private getNoCacheHeaders(): { headers: HttpHeaders } {
    return {
        headers: new HttpHeaders({
        'Cache-Control': 'no-cache'
        })
    };
    }

    startTour(tourId: string): Observable<TourExecution> {
    return this.http.post<TourExecution>(
        `${this.apiUrl}/start`,
        { tour_id: tourId },
        this.getNoCacheHeaders()
    );
    }

    completeTour(execId: string): Observable<void> {
    return this.http.post<void>(
        `${this.apiUrl}/complete`,
        { exec_id: execId },
        this.getNoCacheHeaders()
    );
    }

    abandonTour(execId: string): Observable<void> {
    return this.http.post<void>(
        `${this.apiUrl}/abandon`,
        { exec_id: execId },
        this.getNoCacheHeaders()
    );
    }

    checkKeyPoint(execId: string, keyPointId: string, latitude: number, longitude: number): Observable<{ reached: boolean }> {
    return this.http.post<{ reached: boolean }>(
        `${this.apiUrl}/check-keypoint`,
        { exec_id: execId, key_point_id: keyPointId, latitude, longitude },
        this.getNoCacheHeaders()
    );
    }

    getSimulatedPosition(execId: string): Observable<{ latitude: number, longitude: number }> {
    return this.http.post<{ latitude: number, longitude: number }>(
        `${this.apiUrl}/position`,
        { exec_id: execId },
        this.getNoCacheHeaders()
    );
    }
    getKeyPoints(tourId: string): Observable<KeyPoint[]> {
        return this.http.get<KeyPoint[]>(`${this.toursUrl}/${tourId}/keypoints`);
    }
}
