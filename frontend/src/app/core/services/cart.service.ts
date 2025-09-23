import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';

export interface CartItem {
  id: string;
  tour_id: string;
  tour_name: string;
  price: number;
  quantity?: number;
}

export interface Cart {
  id: string;
  user_id: string;
  items: CartItem[];
  total: number;
}

@Injectable({
  providedIn: 'root'
})
export class CartService {
  private apiUrl = 'http://localhost:8000/api';
  private cartVisibleSubject = new BehaviorSubject<boolean>(false);
  public cartVisible$ = this.cartVisibleSubject.asObservable();
  
  private cartSubject = new BehaviorSubject<Cart | null>(null);
  public cart$ = this.cartSubject.asObservable();

  constructor(private http: HttpClient) {
    this.loadCart();
  }

  // Cart visibility
  toggleCart(): void {
    this.cartVisibleSubject.next(!this.cartVisibleSubject.value);
  }

  hideCart(): void {
    this.cartVisibleSubject.next(false);
  }

  showCart(): void {
    this.cartVisibleSubject.next(true);
  }

  // Cart API calls
  createCart(): Observable<Cart> {
    return this.http.post<Cart>(`${this.apiUrl}/cart`, {});
  }

    getCart(): Observable<Cart> {
    console.log('Making cart request to:', `${this.apiUrl}/cart`);
    return this.http.get<Cart>(`${this.apiUrl}/cart`);
    }

    addItem(tourId: string, tourName: string, price: number): Observable<any> {
    const item = { tour_id: tourId, name: tourName, price: price };
    console.log('Adding item:', item);
    console.log('To URL:', `${this.apiUrl}/cart/items`);
    return this.http.post(`${this.apiUrl}/cart/items`, item);
    }

  removeItem(itemId: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/cart/items/${itemId}`);
  }

  getTotal(): Observable<{ total: number }> {
    return this.http.get<{ total: number }>(`${this.apiUrl}/cart/total`);
  }

  checkout(): Observable<any> {
    return this.http.post(`${this.apiUrl}/cart/checkout`, {});
  }

  private getUserRole(): string {
    const token = localStorage.getItem('access_token');
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return payload.role || payload.user_type || '';
      } catch (error) {
        return '';
      }
    }
    return '';
  }

  // Load cart data
  loadCart(): void {
    const userRole = this.getUserRole();
    if (userRole !== 'tourist') {
      return; // ne ucitava cart ako nije turista
    }
    this.getCart().subscribe({
      next: (cart) => {
        this.cartSubject.next(cart);
      },
      error: (err) => {
        console.error('Error loading cart:', err);
        // Create cart if it doesn't exist
        this.createCart().subscribe({
          next: (newCart) => {
            this.cartSubject.next(newCart);
          },
          error: (createErr) => {
            console.error('Error creating cart:', createErr);
          }
        });
      }
    });
  }

  // Refresh cart after operations
  refreshCart(): void {
    this.loadCart();
  }
}