import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface PurchaseToken {
  id: string;
  tour_id: string;
  tourist_id: string;
  token: string;
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
}