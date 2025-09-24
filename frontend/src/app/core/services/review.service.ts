import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface ReviewRequest {
  rating: number;
  comment: string;
  visited_at: string;
  images: string[];
}

export interface Review {
  id: string;
  tour_id: string;
  tour_name?: string;
  tourist_id: string;
  tourist_name?: string;
  rating: number;
  comment: string;
  visited_at: string;
  created_at: string;
  images: string[];
}

@Injectable({
  providedIn: 'root'
})
export class ReviewService {
  private apiUrl = 'http://localhost:8000/api';

  constructor(private http: HttpClient) {}

  createReview(tourId: string, review: ReviewRequest): Observable<Review> {
    console.log('Sending review request to:', `${this.apiUrl}/tours/${tourId}/reviews`);
    console.log('Review data:', review);
    return this.http.post<Review>(`${this.apiUrl}/tours/${tourId}/reviews`, review);
  }

  getReviewsByTour(tourId: string): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.apiUrl}/tours/${tourId}/reviews`);
  }

  // Za turiste - njihovi reviews
  getMyReviews(): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.apiUrl}/tours/reviews/my`);
  }

  // Za guide-ove - reviews njihovih tura
  getGuideReviews(): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.apiUrl}/tours/reviews/guide`);
  }

  // Za admina - svi reviews
  getAllReviews(): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.apiUrl}/tours/reviews`);
  }
}