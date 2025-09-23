import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Checkpoint } from '../models/checkpoint.model';
import { Tour, TourRequest, TourResponse } from '../models/tour.model';

/*
export interface Tour {
  id: string;
  author_id: string;
  name: string;
  description: string;
  difficulty: string;
  tags: string;
  status: string;
  price: number;
}*/


@Injectable({
  providedIn: 'root'
})
export class TourService {

  private apiUrl = 'http://localhost:8000/api/tours';

  constructor(private http: HttpClient) {}

  private getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('access_token');
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    });
  }

  createTour(tour: TourRequest): Observable<TourResponse> {
    return this.http.post<TourResponse>(this.apiUrl, tour, { headers: this.getAuthHeaders() });
  }

  getTour(id: string): Observable<Tour> {
    return this.http.get<Tour>(`${this.apiUrl}/${id}`, { headers: this.getAuthHeaders() });
  }

  getAllTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(this.apiUrl, { headers: this.getAuthHeaders() });
  }

  getToursByAuthor(authorId: string): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.apiUrl}/authors/${authorId}`, { headers: this.getAuthHeaders() });
  }

  publishTour(id: string): Observable<TourResponse> {
    return this.http.post<TourResponse>(`${this.apiUrl}/${id}/publish`, {}, { headers: this.getAuthHeaders() });
  }

  archiveTour(id: string): Observable<TourResponse> {
    return this.http.post<TourResponse>(`${this.apiUrl}/${id}/archive`, {}, { headers: this.getAuthHeaders() });
  }

  reactivateTour(id: string): Observable<TourResponse> {
    return this.http.post<TourResponse>(`${this.apiUrl}/${id}/reactivate`, {}, { headers: this.getAuthHeaders() });
  }

  addKeyPoints(tourId: string, keyPoints: Checkpoint[]): Observable<Checkpoint[]> {
    return this.http.post<Checkpoint[]>(`${this.apiUrl}/${tourId}/keypoints`, keyPoints, { headers: this.getAuthHeaders() });
  }

  getKeyPointsByTour(tourId: string): Observable<Checkpoint[]> {
    return this.http.get<Checkpoint[]>(`${this.apiUrl}/${tourId}/keypoints`, { headers: this.getAuthHeaders() });
  }
  getAuthorTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.apiUrl}/tours/author-tours`);
  }
}

