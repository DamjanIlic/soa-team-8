import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface PurchaseToken {
  id: string;
  tour_id: string;
  tourist_id: string;
  token: string;
  is_executed: boolean;
  is_reviewed: boolean;
  created_at: string;
}

@Injectable({
  providedIn: 'root'
})
export class PurchaseService {
  private apiUrl = 'http://localhost:8000/api';

  constructor(private http: HttpClient) {}

  getPurchasedTokens(): Observable<PurchaseToken[]> {
    return this.http.get<PurchaseToken[]>(`${this.apiUrl}/cart/tokens/purchased`);
  }

  markTokenAsExecuted(tokenId: string): Observable<any> {
    return this.http.put(`${this.apiUrl}/cart/tokens/${tokenId}/executed`, {});
  }

  markTokenAsReviewed(tokenId: string): Observable<any> {
    return this.http.put(`${this.apiUrl}/cart/tokens/${tokenId}/reviewed`, {});
  }

  getTokenByTour(tourId: string): Observable<PurchaseToken> {
    return this.http.get<PurchaseToken>(`${this.apiUrl}/cart/tokens/tour/${tourId}`);
  }
}