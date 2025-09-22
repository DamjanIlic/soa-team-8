import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:8000/api';
  private authStatusSubject = new BehaviorSubject<boolean>(false);
  public authStatus$ = this.authStatusSubject.asObservable();

  constructor(private http: HttpClient) {
        this.checkInitialAuthStatus();
  }

  register(userData: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/auth/register`, userData);
  }

  login(credentials: { email: string; password: string }): Observable<any> {
    return this.http.post(`${this.apiUrl}/auth/login`, credentials);
  }

  logout(): void {
    localStorage.removeItem('access_token');
    this.authStatusSubject.next(false);
  }

  private checkInitialAuthStatus(): void {
    const token = localStorage.getItem('access_token');
    this.authStatusSubject.next(!!token);
  }

  getUserRole(): string {
    const token = localStorage.getItem('access_token');
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return payload.role || payload.user_type || 'user';
      } catch (error) {
        return '';
      }
    }
    return '';
  }

  setAuthStatus(token: string): void {
    localStorage.setItem('access_token', token);
    this.authStatusSubject.next(true);
  }
  
}